package intelligentIP

import (
	"fmt"
	"regexp"
	"sync"
)

// Pattern defines the intelligentIP Pattern
// It replaces all instances of a found IP with a unique new IP to preserve the logged information while still scrubbing sensitive data
// FileIPs stores a per file (first key) ip<->replacement mapping
type Pattern struct {
	Regex   *regexp.Regexp
	FileIPs map[string]map[string]string
	m       sync.RWMutex

	Suffix string
}

// New .
func New() *Pattern {
	regex, err := regexp.Compile(`\b(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}\b`)
	if err != nil {
		panic(err)
	}
	fip := make(map[string]map[string]string)
	return &Pattern{
		Regex:   regex,
		FileIPs: fip,
		Suffix:  ".ip.fscrub.org",
	}
}

// Find returns how often the regexp was found in string
func (p *Pattern) Find(s string, file string) (int, error) {
	p.checkFile(file)
	find := p.Regex.FindAllString(s, -1)
	if find == nil {
		return 0, nil
	}

	return len(find), nil
}

// Handle returns the regexp handled with target
func (p *Pattern) Handle(s string, file string) (string, error) {
	p.checkFile(file)
	s = p.Regex.ReplaceAllStringFunc(s, func(ip string) string {
		return p.checkIP(file, ip)
	})
	return s, nil
}

// String gives a representation of the pattern for logging
func (p *Pattern) String() string {
	return fmt.Sprintf("intelligentIP")
}

func (p *Pattern) checkFile(file string) {
	p.m.Lock()
	defer p.m.Unlock()
	_, ok := p.FileIPs[file]
	if !ok {
		p.FileIPs[file] = make(map[string]string)
	}
}

func (p *Pattern) checkIP(file, ip string) string {
	p.m.Lock()
	defer p.m.Unlock()
	repl, ok := p.FileIPs[file][ip]
	if !ok {
		repl = fmt.Sprintf("client%d%s", len(p.FileIPs[file]), p.Suffix)
		p.FileIPs[file][ip] = repl
	}
	return repl
}

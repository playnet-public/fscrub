package fscrub

import (
	"errors"
	"strings"
)

type patternType string

func (p patternType) String() string {
	return string(p)
}

// types lists kinds of patterns
var pTypes = struct {
	String patternType
	Regex  patternType
}{
	String: "string",
	Regex:  "regex",
}

// PatternConfig defines the json containing patterns
type PatternConfig struct {
	Patterns []Pattern `json:"patterns"`
}

// Pattern defines a search and replace pattern
type Pattern struct {
	Type   patternType `json:"type"`
	Source string      `json:"source"`
	Target string      `json:"target"`
}

// Find returns how often the source was found in string
func (p *Pattern) Find(s string) int {
	if p.Type == pTypes.String {
		return strings.Count(s, p.Source)
	}

	if p.Type == pTypes.Regex {
		return 1
	}
	return -1
}

// Handle returns the string handled based pattern
func (p *Pattern) Handle(s string) (string, error) {
	if p.Type == pTypes.String {
		return strings.Replace(s, p.Source, p.Target, -1), nil
	}

	if p.Type == pTypes.Regex {
		return s, errors.New("not implemented")
	}
	return s, errors.New("unknown type")
}

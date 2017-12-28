package fscrub

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// Patterns .
type Patterns []Pattern

// Pattern .
type Pattern interface {
	Find(s string) (int, error)
	Handle(s string) (string, error)
	String() string
}

// PatternConfig defines the json containing patterns
type PatternConfig struct {
	Patterns Patterns `json:"patterns"`
}

// UnmarshalJSON stored in PatternConfig
func (c *PatternConfig) UnmarshalJSON(b []byte) error {
	var objMap map[string]*json.RawMessage
	err := json.Unmarshal(b, &objMap)
	if err != nil {
		return err
	}

	var rawMessagesForPatterns []*json.RawMessage
	err = json.Unmarshal(*objMap["patterns"], &rawMessagesForPatterns)
	if err != nil {
		return err
	}

	c.Patterns = make([]Pattern, len(rawMessagesForPatterns))

	var m map[string]string
	for index, rawMessage := range rawMessagesForPatterns {
		err = json.Unmarshal(*rawMessage, &m)
		if err != nil {
			return err
		}

		if m["type"] == "string" {
			var p StringPattern
			err := json.Unmarshal(*rawMessage, &p)
			if err != nil {
				return err
			}
			c.Patterns[index] = &p
		} else if m["type"] == "regex" {
			var p RegexPattern
			err := json.Unmarshal(*rawMessage, &p)
			if err != nil {
				return err
			}
			c.Patterns[index] = &p
		} else {
			return errors.New("unsupported type found")
		}
	}

	// That's it!  We made it the whole way with no errors, so we can return `nil`
	return nil
}

// StringPattern defines a search and replace pattern
type StringPattern struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

// NewStringPattern returns new pattern
func NewStringPattern(src, target string) *StringPattern {
	return &StringPattern{
		Source: src,
		Target: target,
	}
}

// Find returns how often the source was found in string
func (p *StringPattern) Find(s string) (int, error) {
	return strings.Count(s, p.Source), nil
}

// Handle returns the string handled based pattern
func (p *StringPattern) Handle(s string) (string, error) {
	return strings.Replace(s, p.Source, p.Target, -1), nil
}

// String gives a representation of the pattern for logging
func (p *StringPattern) String() string {
	return fmt.Sprintf("Source: %s - Target: %s", p.Source, p.Target)
}

// RegexPattern defines a regex search with static replace
type RegexPattern struct {
	RegexString string `json:"exp"`
	Regex       *regexp.Regexp
	Target      string `json:"target"`
}

// NewRegexPattern compiles the regex and returns pattern
func NewRegexPattern(exp, target string) *RegexPattern {
	return &RegexPattern{
		RegexString: exp,
		Regex:       nil,
		Target:      target,
	}
}

// Find returns how often the regexp was found in string
func (p *RegexPattern) Find(s string) (int, error) {
	if p.Regex == nil {
		regex, err := regexp.Compile(p.RegexString)
		if err != nil {
			return -1, err
		}
		p.Regex = regex
	}
	return 1, nil
}

// Handle returns the regexp handled with target
func (p *RegexPattern) Handle(s string) (string, error) {
	if p.Regex == nil {
		regex, err := regexp.Compile(p.RegexString)
		if err != nil {
			return s, err
		}
		p.Regex = regex
	}
	return s, errors.New("not implemented")
}

// String gives a representation of the pattern for logging
func (p *RegexPattern) String() string {
	return fmt.Sprintf("Regex: %s - Target: %s", p.RegexString, p.Target)
}

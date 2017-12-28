package fscrub

import (
	"testing"
)

// TODO: Find tests for reaching all errors
func TestPatternConfig_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		c       *PatternConfig
		b       []byte
		wantErr bool
	}{
		{
			"basic",
			&PatternConfig{},
			[]byte(`{
				"patterns": [
					{"type": "string", "source": "foo", "target": "bar"},
					{"type": "regex", "exp": "t\\s\\*\\w+", "target": "f *foo"}
				] 
			}`),
			false,
		},
		{
			"unmErr",
			&PatternConfig{},
			[]byte(`.{
				"patterns": [
					{"type": "string", "source": "foo", "target": "bar"}
				] 
			}`),
			true,
		},
		{
			"unknownType",
			&PatternConfig{},
			[]byte(`{
				"patterns": [
					{"type": "foo", "source": "foo", "target": "bar"}
				] 
			}`),
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.UnmarshalJSON(tt.b); (err != nil) != tt.wantErr {
				t.Errorf("PatternConfig.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPattern_Find(t *testing.T) {
	tests := []struct {
		name    string
		p       Pattern
		s       string
		want    int
		wantErr bool
	}{
		{
			"basicString",
			NewStringPattern("foo", "bar"),
			"abc",
			0,
			false,
		},
		{
			"basicRegex",
			NewRegexPattern("t\\s\\*\\w+", "f *foo"),
			"abc",
			1,
			false,
		},
		{
			"oneResult",
			NewStringPattern("foo", "bar"),
			"foo bar",
			1,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.p.Find(tt.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("Pattern.Find() = error %v, want %v", got, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("Pattern.Find() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPattern_Handle(t *testing.T) {
	tests := []struct {
		name    string
		p       Pattern
		s       string
		want    string
		wantErr bool
	}{
		{
			"basicString",
			NewStringPattern("foo", "bar"),
			"abc",
			"abc",
			false,
		},
		{
			"basicRegex",
			NewRegexPattern("t\\s\\*\\w+", "f *foo"),
			"abc",
			"abc",
			true,
		},
		{
			"oneResult",
			NewStringPattern("foo", "bar"),
			"foo bar",
			"bar bar",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.p.Handle(tt.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("Pattern.Handle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Pattern.Handle() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewRegexPattern(t *testing.T) {
	type args struct {
		exp    string
		target string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"basic", args{"t\\s(*\\w+", "bar"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pattern := NewRegexPattern(tt.args.exp, tt.args.target)
			if _, err := pattern.Find(""); (err != nil) != tt.want {
				t.Errorf("NewRegexPattern() = %v, want %v", err, tt.want)
			}
			pattern = NewRegexPattern(tt.args.exp, tt.args.target)
			if _, err := pattern.Handle(""); (err != nil) != tt.want {
				t.Errorf("NewRegexPattern() = %v, want %v", err, tt.want)
			}
		})
	}
}

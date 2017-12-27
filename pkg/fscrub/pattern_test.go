package fscrub

import (
	"testing"
)

func TestPatternType_String(t *testing.T) {
	tests := []struct {
		name string
		p    patternType
		want string
	}{
		{
			"basic",
			"test",
			"test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.String(); got != tt.want {
				t.Errorf("Pattern.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPattern_Find(t *testing.T) {
	tests := []struct {
		name string
		p    *Pattern
		s    string
		want int
	}{
		{
			"basicString",
			&Pattern{Type: pTypes.String, Source: "foo", Target: "bar"},
			"abc",
			0,
		},
		{
			"basicRegex",
			&Pattern{Type: pTypes.Regex, Source: "foo", Target: "bar"},
			"abc",
			1,
		},
		{
			"unknownType",
			&Pattern{Type: "unknown", Source: "foo", Target: "bar"},
			"abc",
			-1,
		},
		{
			"oneResult",
			&Pattern{Type: pTypes.String, Source: "foo", Target: "bar"},
			"foo bar",
			1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.Find(tt.s); got != tt.want {
				t.Errorf("Pattern.Find() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPattern_Handle(t *testing.T) {
	tests := []struct {
		name    string
		p       *Pattern
		s       string
		want    string
		wantErr bool
	}{
		{
			"basicString",
			&Pattern{Type: pTypes.String, Source: "foo", Target: "bar"},
			"abc",
			"abc",
			false,
		},
		{
			"basicRegex",
			&Pattern{Type: pTypes.Regex, Source: "foo", Target: "bar"},
			"abc",
			"abc",
			true,
		},
		{
			"unknownType",
			&Pattern{Type: "unknown", Source: "foo", Target: "bar"},
			"abc",
			"abc",
			true,
		},
		{
			"oneResult",
			&Pattern{Type: pTypes.String, Source: "foo", Target: "bar"},
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

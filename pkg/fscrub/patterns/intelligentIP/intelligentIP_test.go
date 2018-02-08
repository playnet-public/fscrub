package intelligentIP

import (
	"testing"
)

func TestPattern_Find(t *testing.T) {
	p := New()
	tests := []struct {
		name    string
		file    string
		s       string
		want    int
		wantErr bool
	}{
		{
			"file1first",
			"file1",
			"127.0.0.1",
			1,
			false,
		},
		{
			"file2first",
			"file2",
			"127.0.0.1",
			1,
			false,
		},
		{
			"file1second",
			"file1",
			"127.0.0.2",
			1,
			false,
		},
		{
			"file1third",
			"file1",
			"127.0.0.1",
			1,
			false,
		},
		{
			"noFind",
			"file1",
			"abcccccccc.a..2.1.c.1.2.a.2.2.a.x",
			0,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := p.Find(tt.s, tt.file)
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
	p := New()
	tests := []struct {
		name    string
		file    string
		s       string
		want    string
		wantErr bool
	}{
		{
			"file1first",
			"file1",
			"127.0.0.1",
			"client0.ip.fscrub.org",
			false,
		},
		{
			"file2first",
			"file2",
			"127.0.0.1",
			"client0.ip.fscrub.org",
			false,
		},
		{
			"file1second",
			"file1",
			"127.0.0.2",
			"client1.ip.fscrub.org",
			false,
		},
		{
			"file1third",
			"file1",
			"127.0.0.1",
			"client0.ip.fscrub.org",
			false,
		},
		{
			"noFind",
			"file1",
			"abcccccccc.a..2.1.c.1.2.a.2.2.a.x",
			"abcccccccc.a..2.1.c.1.2.a.2.2.a.x",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := p.Handle(tt.s, tt.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("Pattern.Find() = error %v, want %v", got, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("Pattern.Find() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestString(t *testing.T) {
	p := New()
	if p.String() != "intelligentIP" {
		t.Error("invalid String()")
	}
}

// func TestPattern_Handle(t *testing.T) {
// 	tests := []struct {
// 		name    string
// 		p       Pattern
// 		s       string
// 		want    string
// 		wantErr bool
// 	}{
// 		{
// 			"basicString",
// 			New(),
// 			"abc",
// 			"abc",
// 			false,
// 		},
// 		{
// 			"basicRegex",
// 			New(),
// 			"abc",
// 			"abc",
// 			false,
// 		},
// 		{
// 			"findRegex",
// 			New(),
// 			"t *testing.T",
// 			"f *foo.T",
// 			false,
// 		},
// 		{
// 			"findString",
// 			New(),
// 			"foo bar",
// 			"bar bar",
// 			false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got, err := tt.p.Handle(tt.s)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("Pattern.Handle() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if got != tt.want {
// 				t.Errorf("Pattern.Handle() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

func TestPattern_String(t *testing.T) {
	tests := []struct {
		name string
		p    *Pattern
		want string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.String(); got != tt.want {
				t.Errorf("Pattern.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

package model

import (
	"os"
	"testing"
)

func TestDirectory_String(t *testing.T) {
	tests := []struct {
		name string
		d    Directory
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
			if got := tt.d.String(); got != tt.want {
				t.Errorf("Directory.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNoOpAction(t *testing.T) {
	type args struct {
		path string
		file os.FileInfo
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"basic",
			args{"test", nil},
			false,
		},
		{
			"invalid path",
			args{"", nil},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := NoOpAction(tt.args.path, tt.args.file); (err != nil) != tt.wantErr {
				t.Errorf("NoOpAction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

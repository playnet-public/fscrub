package model

import (
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"
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

func TestDirFlags(t *testing.T) {
	dirs := &Directories{}
	want := &Directories{
		"test",
	}
	dirs.Set("test")
	if !reflect.DeepEqual(dirs, want) {
		t.Errorf("Directory.Set() = %v, want %v", dirs, want)
	}
	wantStr := "test"
	if str := dirs.String(); str != wantStr {
		t.Errorf("Directory.String() = %v, want %v", str, wantStr)
	}
	want = &Directories{
		"test",
		"test2",
	}
	dirs.Set("test2")
	if !reflect.DeepEqual(dirs, want) {
		t.Errorf("Directory.Set() = %v, want %v", dirs, want)
	}
	wantStr = "test,test2"
	if str := dirs.String(); str != wantStr {
		t.Errorf("Directory.String() = %v, want %v", str, wantStr)
	}
}

func TestNoOpHandler(t *testing.T) {
	handler := &NoOpHandler{make(chan bool)}
	erc := make(chan error)
	count := 0

	go handler.Run("", erc)

	for {
		select {
		case err := <-erc:
			if err != nil {
				t.Errorf("NoOpHandler.Run() = %v, want %v", err, nil)
			}
			count = count + 1
		case <-time.After(time.Millisecond * 15):
			fmt.Println("timeout 1s")
			return
		}
		if count > 1 {
			t.Errorf("NoOpHandler.Stop() did not work")
		}
		handler.Stop()

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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := NoOpAction(tt.args.path, tt.args.file); (err != nil) != tt.wantErr {
				t.Errorf("NoOpAction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

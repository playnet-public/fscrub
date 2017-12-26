package fscrub

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"go.uber.org/zap"
)

func TestNewFscrub(t *testing.T) {
	log := zap.NewNop()
	tests := []struct {
		name    string
		fscrub  *Fscrub
		wantErr bool
	}{
		{
			"basic",
			NewFscrub(log),
			false,
		},
		{
			"invalidFscrub",
			&Fscrub{log, nil},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fscrub.Validate(); (got != nil) != tt.wantErr {
				t.Errorf("NewFscrub() = %v, want %v", got, tt.wantErr)
			}
		})
	}
}

// TODO: Fix coverage for file tests
func TestFscrub_Handle(t *testing.T) {
	log := zap.NewNop()
	type args struct {
		path     string
		fileInfo os.FileInfo
	}
	tests := []struct {
		name    string
		f       *Fscrub
		args    args
		wantErr bool
	}{
		{
			"basic",
			&Fscrub{log, mockOpenFile("testdata.txt", "ABC\nDEF\nGHI\n")},
			args{"testdata.txt", newMockFileInfo(false)},
			false,
		},
		{
			"dir",
			&Fscrub{log, mockOpenFile("testdata", "")},
			args{"testdata", newMockFileInfo(true)},
			false,
		},
		{
			"fileNotExist",
			&Fscrub{log, mockOpenFile("notexist.txt", "")},
			args{"notexist.txt", newMockFileInfo(false)},
			true,
		},
		{
			"fileNoPerm",
			&Fscrub{log, mockOpenFile("noperm.txt", "")},
			args{"noperm.txt", newMockFileInfo(false)},
			true,
		},
		{
			"fileTimeout",
			&Fscrub{log, mockOpenFile("timeout.txt", "")},
			args{"timeout.txt", newMockFileInfo(false)},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.f.Handle(tt.args.path, tt.args.fileInfo); (err != nil) != tt.wantErr {
				t.Errorf("Fscrub.Handle() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFscrub_HandleLine(t *testing.T) {
	log := zap.NewNop()
	type args struct {
		path   string
		lineNo int
		line   string
	}
	tests := []struct {
		name    string
		f       *Fscrub
		args    args
		want    string
		wantErr bool
	}{
		{
			"basic",
			NewFscrub(log),
			args{"testfile.txt", 0, "ABC"},
			"ABC",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.f.HandleLine(tt.args.path, tt.args.lineNo, tt.args.line)
			if (err != nil) != tt.wantErr {
				t.Errorf("Fscrub.HandleLine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Fscrub.HandleLine() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TODO: Solve problem where this won't test
/*func TestFscrub_openFile(t *testing.T) {
	log := zap.NewNop()
	tests := []struct {
		name    string
		f       *Fscrub
		path    string
		wantErr bool
	}{
		{
			"basic",
			NewFscrub(log),
			"testdata.txt",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, path, err := createTempFile(tt.path)
			if err != nil {
				t.Errorf("Fscrub.createTempFile() error = %v", err)
			}
			file.Close()
			_, err = tt.f.openFile(path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Fscrub.openFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}*/

type mockFileInfo struct {
	dir bool
}

func (f *mockFileInfo) Name() string {
	return "test"
}
func (f *mockFileInfo) Size() int64 {
	return 1
}
func (f *mockFileInfo) Mode() os.FileMode {
	return os.ModeTemporary
}
func (f *mockFileInfo) ModTime() time.Time {
	return time.Now()
}
func (f *mockFileInfo) IsDir() bool {
	return f.dir
}
func (f *mockFileInfo) Sys() interface{} {
	return mockFileInfo{}
}

func newMockFileInfo(isDir bool) os.FileInfo {
	info := &mockFileInfo{isDir}
	return info
}

func mockOpenFile(path, content string) func(path string) (*os.File, error) {
	byteSlice := []byte(content)
	return func(path string) (*os.File, error) {
		if path == "notexist.txt" {
			return nil, os.ErrNotExist
		}
		if path == "noperm.txt" {
			return nil, os.ErrPermission
		}
		file, _, err := createTempFile(path)

		_, err = file.Write(byteSlice)
		if err != nil {
			return nil, err
		}
		err = file.Sync()
		if err != nil {
			fmt.Print(err)
			return nil, err
		}
		if path == "timeout.txt" {
			file.Close()
		}
		return file, nil
	}
}

func createTempFile(path string) (*os.File, string, error) {
	tmpDir, err := ioutil.TempDir("", "fscrubTests")
	if err != nil {
		return nil, "", err
	}

	fileName := filepath.Base(path)

	file, err := ioutil.TempFile(tmpDir, fileName)
	if err != nil {
		return nil, "", err
	}
	return file, filepath.Join(tmpDir, fileName), nil
}

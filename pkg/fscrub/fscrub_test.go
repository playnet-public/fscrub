package fscrub

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/playnet-public/fscrub/pkg/primitives"
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
			NewFscrub(log, true),
			false,
		},
		{
			"invalidFscrub",
			&Fscrub{log: log, fileOpener: nil},
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
	patterns := []Pattern{
		{pTypes.String, "foo", "bar"},
	}
	errPatterns := []Pattern{
		{pTypes.Regex, "foo", "bar"},
	}
	type args struct {
		path     string
		fileInfo os.FileInfo
	}
	tests := []struct {
		name    string
		f       *Fscrub
		content string
		args    args
		wantErr bool
	}{
		{
			"basic",
			&Fscrub{log: log,
				fileOpener:  primitives.OpenFile(log),
				fileWriter:  mockWriteFile("testdata.txt"),
				fileUpdater: FileUpdater(NewFscrub(log, false)),
			},
			"ABC\nDEF\nGHI\n",
			args{"testdata.txt", newMockFileInfo(false)},
			false,
		},
		{
			"basicReplace",
			&Fscrub{log: log,
				fileOpener:  primitives.OpenFile(log),
				fileWriter:  mockWriteFile("testdata.txt"),
				fileUpdater: mockUpdateFile,
				patterns:    patterns,
			},
			"foo\nbar\nfoo\n",
			args{"testdata.txt", newMockFileInfo(false)},
			false,
		},
		{
			"basicSkip",
			&Fscrub{log: log,
				fileOpener:  primitives.OpenFile(log),
				fileWriter:  mockWriteFile("testdata.txt"),
				fileUpdater: mockUpdateFile,
				patterns:    patterns,
			},
			"//-ignore: github.com/playnet-public/fscrub",
			args{"testdata.txt", newMockFileInfo(false)},
			false,
		},
		{
			"handleErr",
			&Fscrub{log: log,
				fileOpener:  primitives.OpenFile(log),
				fileWriter:  mockWriteFile("testdata.txt"),
				fileUpdater: mockUpdateFile,
				patterns:    errPatterns,
			},
			"ABC\nDEF\nGHI\n",
			args{"testdata.txt", newMockFileInfo(false)},
			true,
		},
		{
			"dir",
			&Fscrub{log: log,
				fileOpener:  mockOpenFile("testdata", ""),
				fileWriter:  mockWriteFile("testdata"),
				fileUpdater: mockUpdateFile,
			},
			"",
			args{"testdata", newMockFileInfo(true)},
			false,
		},
		{
			"fileNotExist",
			&Fscrub{log: log,
				fileOpener:  mockOpenFile("notexist.txt", ""),
				fileWriter:  mockWriteFile("notexist.txt"),
				fileUpdater: mockUpdateFile,
			},
			"",
			args{"notexist.txt", newMockFileInfo(false)},
			true,
		},
		{
			"fileNoPerm",
			&Fscrub{log: log,
				fileOpener:  mockOpenFile("noperm.txt", ""),
				fileWriter:  mockWriteFile("noperm.txt"),
				fileUpdater: mockUpdateFile,
			},
			"",
			args{"noperm.txt", newMockFileInfo(false)},
			true,
		},
		{
			"fileTimeout",
			&Fscrub{log: log,
				fileOpener:  mockOpenFile("timeout.txt", ""),
				fileWriter:  mockWriteFile("timeout.txt"),
				fileUpdater: mockUpdateFile,
			},
			"",
			args{"timeout.txt", newMockFileInfo(false)},
			true,
		},
		{
			"fileUndefErr",
			&Fscrub{log: log,
				fileOpener:  mockOpenFile("undefErr.txt", ""),
				fileWriter:  mockWriteFile("undefErr.txt"),
				fileUpdater: mockUpdateFile,
			},
			"",
			args{"undefErr.txt", newMockFileInfo(false)},
			true,
		},
		{
			"failUpdate",
			&Fscrub{log: log,
				fileOpener:  primitives.OpenFile(log),
				fileWriter:  mockWriteFile("failupdate.txt"),
				fileUpdater: mockUpdateFile,
				patterns:    patterns,
			},
			"foo\nbar\nfoo\n",
			args{"failupdate.txt", newMockFileInfo(false)},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, path, _ := createTempFile(tt.args.path, tt.content)
			file.Close()
			file, err := os.OpenFile(path, os.O_RDONLY, 0666)
			if err != nil {
				t.Errorf("Fscrub.UpdateFile() error = %v when opening file", err)
			}
			data, err := ioutil.ReadAll(file)
			if err != nil {
				t.Errorf("Fscrub.UpdateFile() error = %v when reading file", err)
			}
			if string(data) != tt.content {
				t.Errorf("Fscrub.UpdateFile() newContent = %s, want %v", data, tt.content)
			}

			if err := tt.f.Handle(path, tt.args.fileInfo); (err != nil) != tt.wantErr {
				t.Errorf("Fscrub.Handle() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFscrub_HandleLine(t *testing.T) {
	log := zap.NewNop()
	patterns := []Pattern{
		{pTypes.String, "foo", "bar"},
	}
	errPatterns := []Pattern{
		{pTypes.Regex, "foo", "bar"},
	}
	tests := []struct {
		name    string
		f       *Fscrub
		line    Line
		want    string
		wantErr bool
	}{
		{
			"basic",
			NewFscrub(log, true, patterns...),
			Line{"testfile.txt", 0, "ABC", false},
			"ABC",
			false,
		},
		{
			"findFoo",
			NewFscrub(log, true, patterns...),
			Line{"testfile.txt", 0, "foo", false},
			"foo",
			false,
		},
		{
			"handleFoo",
			NewFscrub(log, false, patterns...),
			Line{"testfile.txt", 0, "foo", false},
			"bar",
			false,
		},
		{
			"handleFooErr",
			NewFscrub(log, false, errPatterns...),
			Line{"testfile.txt", 0, "foo", false},
			"foo",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.f.HandleLine(tt.line)
			if (err != nil) != tt.wantErr {
				t.Errorf("Fscrub.HandleLine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.Text != tt.want {
				t.Errorf("Fscrub.HandleLine() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFscrub_FileUpdater(t *testing.T) {
	log := zap.NewNop()
	tests := []struct {
		name    string
		f       *Fscrub
		path    string
		content string
		wantErr bool
	}{
		{
			"basic",
			&Fscrub{log: log,
				fileOpener: mockOpenFile("testdata.txt", "ABC\nDEF\nGHI\n"),
				fileWriter: primitives.WriteFile(log),
			},
			"testdata.txt",
			"foo\nbar\n",
			false,
		},
		{
			"fileNotExist",
			&Fscrub{log: log,
				fileOpener: mockOpenFile("notexist.txt", "foo"),
				fileWriter: mockWriteFile("notexist.txt"),
			},
			"notexist.txt",
			"sometempfile",
			true,
		},
		{
			"fileNoPerm",
			&Fscrub{log: log,
				fileOpener: mockOpenFile("noperm.txt", "foo"),
				fileWriter: mockWriteFile("noperm.txt"),
			},
			"noperm.txt",
			"sometempfile",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			file, path, _ := createTempFile(tt.path)
			file.Close()

			if err := FileUpdater(tt.f)(path, tt.content); (err != nil) != tt.wantErr {
				t.Errorf("Fscrub.FileUpdater() error = %v, wantErr %v", err, tt.wantErr)
			}
			file, err := os.OpenFile(path, os.O_RDONLY, 0666)
			if err != nil {
				t.Errorf("Fscrub.FileUpdater() error = %v when opening file", err)
			}
			data, err := ioutil.ReadAll(file)
			if err != nil {
				t.Errorf("Fscrub.FileUpdater() error = %v when reading file", err)
			}
			if string(data) != tt.content {
				t.Errorf("Fscrub.FileUpdater() newContent = %s, want %v", data, tt.content)
			}
		})
	}
}

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
		if strings.Contains(path, "notexist.txt") {
			return nil, os.ErrNotExist
		}
		if strings.Contains(path, "noperm.txt") {
			return nil, os.ErrPermission
		}
		if strings.Contains(path, "undefErr.txt") {
			return nil, errors.New("undefErr")
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
		if strings.Contains(path, "timeout.txt") {
			file.Close()
		}
		return file, nil
	}
}

func mockWriteFile(path string) func(path string, data []byte) error {
	return func(path string, data []byte) error {
		if strings.Contains(path, "notexist.txt") {
			return os.ErrNotExist
		}
		if strings.Contains(path, "noperm.txt") {
			return os.ErrPermission
		}
		file, tmpPath, err := createTempFile(path)
		if err != nil {
			log.Fatal(err)
			return err
		}
		file.Close()
		err = ioutil.WriteFile(tmpPath, data, 0666)
		if err != nil {
			log.Fatal(err)
			return err
		}
		return nil
	}
}

func mockUpdateFile(path, content string) error {
	if strings.Contains(path, "failupdate.txt") {
		return errors.New("update error")
	}
	return nil
}

func createTempFile(path string, content ...string) (*os.File, string, error) {

	tmpDir, err := ioutil.TempDir("", "fscrubTests")
	if err != nil {
		return nil, "", err
	}

	fileName := filepath.Base(path)

	file, err := ioutil.TempFile(tmpDir, fileName)
	if err != nil {
		return nil, "", err
	}
	if len(content) > 0 {
		_, err = file.Write([]byte(content[0]))

	} else {
		_, err = file.Write([]byte("sometempfile"))
	}
	if err != nil {
		return nil, "", err
	}
	err = file.Sync()
	if err != nil {
		fmt.Print(err)
		return nil, "", err
	}

	return file, file.Name(), nil
}

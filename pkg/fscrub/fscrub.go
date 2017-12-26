package fscrub

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/zap"
)

// Fscrub defines an action for scrubbing text files
type Fscrub struct {
	log        *zap.Logger
	fileOpener func(path string) (*os.File, error)
}

// NewFscrub with logger
func NewFscrub(log *zap.Logger) *Fscrub {
	f := &Fscrub{
		log: log,
	}
	f.fileOpener = f.openFile
	return f
}

// Validate that fscrub has fileOpener
func (f *Fscrub) Validate() error {
	if f.fileOpener == nil {
		return errors.New("fscrub missing fileOpener")
	}
	return nil
}

// Handle path and take actions if fileInfo matches required criteria
func (f *Fscrub) Handle(path string, fileInfo os.FileInfo) error {

	if fileInfo.IsDir() {
		return nil
	}

	f.log.Info(
		"running action",
		zap.String("action", "fscrub"),
		zap.String("path", path),
		zap.String("file", fileInfo.Name()),
	)

	file, err := f.fileOpener(path)
	defer file.Close()
	if err != nil {
		if os.IsNotExist(err) {
			f.log.Error("file does not exist", zap.String("file", path), zap.Error(err))
			return err
		}
		if os.IsPermission(err) {
			f.log.Error("file permission denied", zap.String("file", path), zap.Error(err))
			return err
		}
	}

	f.log.Info("file scan started", zap.String("file", path))
	scanner := bufio.NewScanner(file)
	lineNo := 0
	fmt.Println(scanner.Scan())
	fmt.Println(scanner.Text())
	for scanner.Scan() {
		line := scanner.Text()
		newLine, err := f.HandleLine(path, lineNo, line)
		if err != nil {
			f.log.Error("failed handling line", zap.String("file", path), zap.Int("line", lineNo), zap.String("text", line))
			return err
		}
		line = newLine
		lineNo = lineNo + 1
	}
	err = scanner.Err()
	if err == nil {
		f.log.Info("file scan finished", zap.String("file", path))
	} else {
		f.log.Error("file scan failed", zap.String("file", path), zap.Error(err))
		return err
	}

	return nil
}

// HandleLine and return new line or error
func (f *Fscrub) HandleLine(path string, lineNo int, line string) (string, error) {
	f.log.Debug("handling line", zap.String("file", path), zap.Int("line", lineNo), zap.String("text", line))
	return line, nil
}

func (f *Fscrub) openFile(path string) (*os.File, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		f.log.Error("could not get abs path", zap.String("file", path), zap.Error(err))
		return nil, err
	}
	return os.OpenFile(path, os.O_RDWR, 0666)
}

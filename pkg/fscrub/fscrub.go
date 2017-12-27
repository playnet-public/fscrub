package fscrub

import (
	"bufio"
	"errors"
	"os"

	"github.com/playnet-public/fscrub/pkg/primitives"
	"go.uber.org/zap"
)

// Fscrub defines an action for scrubbing text files
type Fscrub struct {
	patterns []Pattern
	dry      bool

	log        *zap.Logger
	fileOpener func(path string) (*os.File, error)
}

// NewFscrub with logger
func NewFscrub(log *zap.Logger, dryrun bool, patterns ...Pattern) *Fscrub {
	f := &Fscrub{
		patterns: patterns,
		log:      log,
		dry:      dryrun,
	}
	f.fileOpener = primitives.OpenFile(log)
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
			f.log.Error("file does not exist",
				zap.String("file", path),
				zap.Error(err))
			return err
		}
		if os.IsPermission(err) {
			f.log.Error("file permission denied",
				zap.String("file", path),
				zap.Error(err))
			return err
		}
	}

	var newFile []string

	f.log.Info("file scan started", zap.String("file", path))
	scanner := bufio.NewScanner(file)
	lineNo := 0
	for scanner.Scan() {
		line := Line{
			Path:    path,
			No:      lineNo,
			Text:    scanner.Text(),
			Changed: false,
		}
		new, err := f.HandleLine(line)
		if err != nil {
			f.log.Error("failed handling line",
				zap.String("file", path),
				zap.Int("line", lineNo), zap.String("text", line.Text))
			return err
		}
		line = new
		lineNo = lineNo + 1
		newFile = append(newFile, line.Text)
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

// Line represents a line handled by Fscrub
type Line struct {
	Path    string
	No      int
	Text    string
	Changed bool
}

// HandleLine and return new line or error
func (f *Fscrub) HandleLine(line Line) (Line, error) {
	//f.log.Debug("handling line",
	//	zap.String("file", line.Path),
	//	zap.Int("line", line.No),
	//	zap.String("text", line.Text))

	for _, p := range f.patterns {
		count := p.Find(line.Text)
		if count > 0 {
			f.log.Info("found pattern",
				zap.String("file", line.Path),
				zap.Int("line", line.No),
				zap.String("text", line.Text),
				zap.String("type", p.Type.String()),
				zap.String("pattern", p.Source),
			)
			f.log.Info("handling pattern",
				zap.String("file", line.Path),
				zap.Int("line", line.No),
				zap.String("text", line.Text),
				zap.String("type", p.Type.String()),
				zap.String("pattern", p.Source),
				zap.String("with", p.Target),
			)
			if !f.dry {
				new, err := p.Handle(line.Text)
				if err != nil {
					f.log.Error("handling pattern failed",
						zap.String("file", line.Path),
						zap.Int("line", line.No),
						zap.String("text", line.Text),
						zap.String("type", p.Type.String()),
						zap.String("pattern", p.Source),
						zap.Error(err),
					)
					return line, err
				}
				line.Text = new
				line.Changed = true
			}
		}
	}
	return line, nil
}

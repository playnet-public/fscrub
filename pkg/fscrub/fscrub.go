package fscrub

import (
	"bufio"
	"errors"
	"github.com/playnet-public/libs/log"
	"os"
	"strings"

	"github.com/playnet-public/fscrub/pkg/primitives"
	"go.uber.org/zap"
)

// Fscrub defines an action for scrubbing text files
type Fscrub struct {
	patterns Patterns
	dry      bool

	log         *log.Logger
	fileOpener  func(path string) (*os.File, error)
	fileWriter  func(path string, data []byte) error
	fileUpdater func(path, content string) error
}

// NewFscrub with logger
func NewFscrub(log *log.Logger, dryrun bool, patterns ...Pattern) *Fscrub {
	f := &Fscrub{
		patterns: patterns,
		log:      log,
		dry:      dryrun,
	}
	f.fileOpener = primitives.OpenFile(log)
	f.fileWriter = primitives.WriteFile(log)
	f.fileUpdater = FileUpdater(f)
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

	changed := false
	var newFile []string
	{
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
			f.log.Error("undefined file error",
				zap.String("file", path),
				zap.Error(err))
			return err
		}

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
			if line.Text == primitives.BuildIgnoreHeader() {
				f.log.Info(
					"skipping file",
					zap.String("action", "fscrub"),
					zap.String("path", path),
					zap.String("reason", "skip-header"),
				)
				return nil
			}
			skipLine := false
			for _, hl := range primitives.BuildHeader() {
				if line.Text == hl {
					skipLine = true
				}
			}
			if skipLine {
				continue
			}
			new, err := f.HandleLine(line)
			if err != nil {
				f.log.Error("failed handling line",
					zap.String("file", path),
					zap.Int("line", lineNo),
					zap.String("text", line.Text),
					zap.Error(err))
				return err
			}
			if new.Changed {
				changed = true
				line = new
			}
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
	}

	if changed {
		newFile = append(primitives.BuildHeader(), newFile...)
		err := f.fileUpdater(path, strings.Join(newFile, "\n"))
		if err != nil {
			f.log.Error("updating file failed",
				zap.String("file", path),
				zap.Error(err))
			return err
		}
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
		count, err := p.Find(line.Text)
		if err != nil {
			f.log.Error("finding pattern failed",
				zap.String("file", line.Path),
				zap.Int("line", line.No),
				zap.String("text", line.Text),
				zap.String("pattern", p.String()),
				zap.Error(err),
			)
			return line, err
		}
		if count > 0 {
			f.log.Info("found pattern",
				zap.String("file", line.Path),
				zap.Int("line", line.No),
				zap.String("text", line.Text),
				zap.String("pattern", p.String()),
			)
			f.log.Info("handling pattern",
				zap.String("file", line.Path),
				zap.Int("line", line.No),
				zap.String("text", line.Text),
				zap.String("pattern", p.String()),
			)
			if !f.dry {
				new, err := p.Handle(line.Text)
				if err != nil {
					f.log.Error("handling pattern failed",
						zap.String("file", line.Path),
						zap.Int("line", line.No),
						zap.String("text", line.Text),
						zap.String("pattern", p.String()),
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

// FileUpdater returns update function for files
func FileUpdater(f *Fscrub) func(path, content string) error {
	return func(path, content string) error {
		err := f.fileWriter(path, []byte(content))
		if err == nil {
			f.log.Info("updating file finished", zap.String("file", path))
		} else {
			f.log.Error("updating file failed",
				zap.String("file", path),
				zap.Error(err))
			return err
		}
		return nil
	}
}

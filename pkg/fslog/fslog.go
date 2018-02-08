package fslog

import (
	"os"

	"github.com/playnet-public/libs/log"

	"go.uber.org/zap"
)

// FsLogger provides a logging action for the fshandler
type FsLogger struct {
	log *log.Logger
}

// NewFsLogger with logger
func NewFsLogger(log *log.Logger) *FsLogger {
	return &FsLogger{
		log: log,
	}
}

// Log the provided path and file
func (f *FsLogger) Log(path string, file os.FileInfo) error {
	f.log.Info(
		"running action",
		zap.String("action", "fslog"),
		zap.String("path", path),
		zap.String("file", file.Name()),
	)
	return nil
}

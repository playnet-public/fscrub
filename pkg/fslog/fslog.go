package fslog

import (
	"os"

	"go.uber.org/zap"
)

// FsLogger provides a logging action for the fshandler
type FsLogger struct {
	log *zap.Logger
}

// NewFsLogger with logger
func NewFsLogger(log *zap.Logger) *FsLogger {
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

package primitives

import (
	"os"
	"path/filepath"

	"go.uber.org/zap"
)

// OpenFile returns a primitive function for opening files
func OpenFile(log *zap.Logger) func(path string) (*os.File, error) {
	return func(path string) (*os.File, error) {
		path, err := filepath.Abs(path)
		if err != nil {
			log.Error("could not get abs path", zap.String("file", path), zap.Error(err))
			return nil, err
		}
		return os.OpenFile(path, os.O_RDWR, 0666)
	}
}

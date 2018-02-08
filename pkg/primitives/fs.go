package primitives

import (
	"io/ioutil"
	"github.com/playnet-public/libs/log"
	"os"
	"path/filepath"

	"go.uber.org/zap"
)

// OpenFile returns a primitive function for opening and reading files
func OpenFile(log *log.Logger) func(path string) (*os.File, error) {
	return func(path string) (*os.File, error) {
		path, err := filepath.Abs(path)
		if err != nil {
			log.Error("could not get abs path", zap.String("file", path), zap.Error(err))
			return nil, err
		}
		return os.OpenFile(path, os.O_RDONLY, 0666)
	}
}

// WriteFile returns a primitive function for writing files replacing their content
func WriteFile(log *log.Logger) func(path string, data []byte) error {
	return func(path string, data []byte) error {
		path, err := filepath.Abs(path)
		if err != nil {
			log.Error("could not get abs path", zap.String("file", path), zap.Error(err))
			return err
		}
		return ioutil.WriteFile(path, data, 0666)
	}
}

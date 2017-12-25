package model

import (
	"os"
)

// Directory used over multiple packages
type Directory string

func (d Directory) String() string {
	return string(d)
}

// Handler defines the functions for handling directories
type Handler func(dir Directory, erc chan error)

// NoOpHandler does nothing
func NoOpHandler(dir Directory, erc chan error) {
	erc <- nil
}

// Action defines the functions for processing files
type Action func(path string, file os.FileInfo) error

// NoOpAction does nothing
func NoOpAction(path string, file os.FileInfo) error {
	return nil
}

package model

import (
	"os"
	"strings"
	"time"
)

// Directory used over multiple packages
type Directory string

func (d Directory) String() string {
	return string(d)
}

// Directories array of Directory
type Directories []Directory

func (d *Directories) String() string {
	var str []string
	for _, dir := range *d {
		str = append(str, dir.String())
	}
	return strings.Join(str, ",")
}

// Set new Directory to Directories
func (d *Directories) Set(value string) error {
	*d = append(*d, Directory(value))
	return nil
}

// Handler defines the functions for handling directories
type Handler interface {
	Run(dir Directory, erc chan error)
	Stop()
}

// NoOpHandler does nothing
type NoOpHandler struct {
	interrupt chan bool
}

// Run the NoOpHandler while sleeping for 1 Sec every run
func (h *NoOpHandler) Run(dir Directory, erc chan error) {
	for {
		select {
		case erc <- nil:
			time.Sleep(time.Millisecond * 5)
		case <-h.interrupt:
			return
		case <-time.After(time.Millisecond * 5):
		}
	}
}

// Stop NoOp
func (h *NoOpHandler) Stop() {
	h.interrupt <- true
}

// Action defines the functions for processing files
type Action func(path string, file os.FileInfo) error

// NoOpAction does nothing
func NoOpAction(path string, file os.FileInfo) error {
	return nil
}

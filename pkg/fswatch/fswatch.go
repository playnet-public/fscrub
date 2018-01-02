package fswatch

import (
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/playnet-public/fscrub/pkg/model"
	"go.uber.org/zap"
)

// Watcher defines the dir watching handler
type Watcher struct {
	log       *zap.Logger
	interrupt chan bool
	actions   []model.Action
	watcher   *fsnotify.Watcher
}

// NewWatcher with logger
func NewWatcher(log *zap.Logger, actions ...model.Action) *Watcher {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Error("error initializing fswatch", zap.Error(err))
		return nil
	}

	interrupt := make(chan bool)

	w := &Watcher{
		log:       log,
		interrupt: interrupt,
		actions:   actions,
		watcher:   watcher,
	}

	go w.watch()

	return w
}

func (w *Watcher) watch() {
	for {
		select {
		case event := <-w.watcher.Events:
			w.log.Debug("file event captured", zap.String("event", event.String()))
			if event.Op&fsnotify.Write == fsnotify.Write {
				w.log.Info("handling file event", zap.String("type", "modified"), zap.String("file", event.Name))
				err := w.handle(event.Name)
				if err != nil {
					w.log.Error("failed handling file event",
						zap.String("type", ""),
						zap.String("file", event.Name),
					)
				}
			}
			if event.Op&fsnotify.Create == fsnotify.Create {
				w.log.Info("handling file event", zap.String("type", "created"), zap.String("file", event.Name))
				err := w.handle(event.Name)
				if err != nil {
					w.log.Error("failed handling file event",
						zap.String("type", ""),
						zap.String("file", event.Name),
					)
				}
			}
		case err := <-w.watcher.Errors:
			w.log.Error("error in fswatch", zap.Error(err))
		case <-w.interrupt:
			return
		}
	}
}

// Run the watcher for dir
func (w *Watcher) Run(dir model.Directory, erc chan error) {
	err := w.watcher.Add(dir.String())
	if err != nil {
		erc <- err
	}
}

func (w *Watcher) handle(path string) error {
	file, err := os.Lstat(path)
	if err != nil {
		w.log.Error("unable to handle path", zap.String("path", path), zap.Error(err))
		return err
	}
	w.log.Info("handling path", zap.String("path", path))
	for _, a := range w.actions {
		err := a(path, file)
		if err != nil {
			return err
		}
	}
	return nil
}

// Stop watching all dirs
func (w *Watcher) Stop() {
	w.log.Info("stopping watcher")
	w.interrupt <- true
	w.watcher.Close()
	w.log.Info("watcher stopped")
}

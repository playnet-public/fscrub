package fscrub

import (
	"github.com/playnet-public/fscrub/pkg/model"
	"go.uber.org/zap"

	"github.com/pkg/errors"
)

// Fscrub holds the main config
type Fscrub struct {
	Dirs    []model.Directory
	Watcher model.Handler
	Crawler model.Handler
	Actions []model.Action

	log *zap.Logger
}

// NewFscrub creates a fscrub instance with defaults if necessary
func NewFscrub(
	dirs []model.Directory,
	watcher model.Handler,
	crawler model.Handler,
	actions []model.Action,
	log *zap.Logger,
) *Fscrub {
	if watcher == nil {
		watcher = model.NoOpHandler
	}
	if crawler == nil {
		crawler = model.NoOpHandler
	}
	f := &Fscrub{
		Dirs:    dirs,
		Watcher: watcher,
		Crawler: crawler,
		Actions: actions,
		log:     log,
	}
	return f
}

// Validate Fscrub integrity
func (f *Fscrub) Validate() error {
	if len(f.Dirs) < 1 {
		return errors.New("at least one dir required")
	}
	if len(f.Actions) < 1 {
		return errors.New("at least one action required")
	}
	for _, act := range f.Actions {
		if act == nil {
			return errors.New("action must not be nil")
		}
	}
	return nil
}

// Run watcher and crawler for all dirs
func (f *Fscrub) Run() error {
	f.log.Info("running fscrub")
	erc := make(chan error)
	for _, d := range f.Dirs {
		f.log.Info("starting watcher", zap.String("dir", d.String()))
		go f.Watcher(d, erc)
		f.log.Info("starting crawler", zap.String("dir", d.String()))
		go f.Crawler(d, erc)
	}
	err := <-erc
	if err != nil {
		f.log.Error("fscrub handler error", zap.Error(err))
		return errors.Wrap(err, "error encountered while running fscrub")
	}
	f.log.Info("fscrub finished")
	return nil
}

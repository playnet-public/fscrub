package fshandle

import (
	"reflect"

	"github.com/playnet-public/fscrub/pkg/model"
	"go.uber.org/zap"

	"github.com/pkg/errors"
)

// FsHandler holds the main config
type FsHandler struct {
	Dirs     model.Directories
	Handlers []model.Handler

	log *zap.Logger
}

// NewFsHandler creates a fscrub instance with defaults if necessary
func NewFsHandler(
	dirs model.Directories,
	handlers []model.Handler,
	log *zap.Logger,
) *FsHandler {
	if len(dirs) < 1 {
		log.Error("invalid handler", zap.Error(errors.New("at least one dir required")))
		return nil
	}
	f := &FsHandler{
		Dirs:     dirs,
		Handlers: handlers,
		log:      log,
	}
	return f
}

// Run watcher and crawler for all dirs
func (f *FsHandler) Run() error {
	f.log.Info("running fscrub")
	erc := make(chan error)
	for _, d := range f.Dirs {
		for _, h := range f.Handlers {
			f.log.Info("starting handler", zap.String("type", reflect.TypeOf(h).String()), zap.String("dir", d.String()))
			go h.Run(d, erc)
			defer h.Stop()
		}
	}
	err := <-erc
	if err != nil {
		f.log.Error("fscrub handler error", zap.Error(err))
		return errors.Wrap(err, "error encountered while running fscrub")
	}
	f.log.Info("fscrub finished")
	return nil
}

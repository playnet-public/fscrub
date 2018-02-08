package fscrawl

import (
	"errors"
	"github.com/playnet-public/libs/log"
	"os"
	"path/filepath"

	"github.com/playnet-public/fscrub/pkg/model"
	"go.uber.org/zap"
)

// Crawler defines the dir crawling handler
type Crawler struct {
	log       *log.Logger
	interrupt chan bool
	actions   []model.Action
}

// NewCrawler with logger
func NewCrawler(log *log.Logger, actions ...model.Action) *Crawler {
	if actions == nil {
		actions = []model.Action{model.NoOpAction}
	}
	return &Crawler{
		log:       log,
		interrupt: make(chan bool),
		actions:   actions,
	}
}

// Validate crawler integrity
func (c *Crawler) Validate() error {
	if c.log == nil {
		return errors.New("log must not be nil")
	}
	return nil
}

// Run the crawler for dir
func (c *Crawler) Run(dir model.Directory, erc chan error) {
	c.log.Info("start handling", zap.String("dir", dir.String()), zap.String("handler", "crawler"))
	defer c.log.Info("stop handling", zap.String("dir", dir.String()), zap.String("handler", "crawler"))
	if len(dir) < 1 {
		erc <- errors.New("invalid dir")
		return
	}
	err := filepath.Walk(dir.String(), c.handle)
	erc <- err
}

func (c *Crawler) handle(path string, file os.FileInfo, err error) error {
	if err != nil {
		c.log.Error("unable to handle path", zap.String("path", path), zap.Error(err))
		return err
	}
	c.log.Info("handling path", zap.String("path", path))
	for _, a := range c.actions {
		err := a(path, file)
		if err != nil {
			return err
		}
	}
	return nil
}

// Stop crawling
func (c *Crawler) Stop() {}

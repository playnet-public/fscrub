package fscrawl

import (
	"errors"

	"github.com/playnet-public/fscrub/pkg/model"
	"go.uber.org/zap"
)

// Crawler .
type Crawler struct {
	log *zap.Logger
}

// NewCrawler with logger
func NewCrawler(log *zap.Logger) *Crawler {
	return &Crawler{
		log: log,
	}
}

// Validate crawler integrity
func (c *Crawler) Validate() error {
	if c.log == nil {
		return errors.New("log must not be nil")
	}
	return nil
}

// Run the crawler for path
func (c *Crawler) Run(dir model.Directory, erc chan error) {
	c.log.Info("running crawler", zap.String("dir", dir.String()))
	if len(dir) < 1 {
		erc <- errors.New("invalid dir")
		return
	}
	erc <- nil
}

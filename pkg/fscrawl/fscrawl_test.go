package fscrawl

import (
	"testing"

	"github.com/playnet-public/fscrub/pkg/model"
	"go.uber.org/zap"
)

func TestCrawler_Validate(t *testing.T) {
	log := zap.NewNop()
	type args struct {
		log *zap.Logger
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"basic",
			args{log},
			false,
		},
		{
			"nilLog",
			args{nil},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			crawl := NewCrawler(tt.args.log)
			if err := crawl.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Crawler.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCrawler_Run(t *testing.T) {

	log := zap.NewNop()

	tests := []struct {
		name    string
		c       *Crawler
		path    model.Directory
		wantErr bool
	}{
		{
			"basic",
			&Crawler{log},
			"test",
			false,
		},
		{
			"invalid path",
			&Crawler{log},
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			erc := make(chan error)
			go tt.c.Run(tt.path, erc)
			err := <-erc
			if (err != nil) != tt.wantErr {
				t.Errorf("Crawler.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

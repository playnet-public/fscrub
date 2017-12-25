package fscrub

import (
	"os"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/playnet-public/fscrub/pkg/fscrawl"
	"github.com/playnet-public/fscrub/pkg/model"
	"go.uber.org/zap"
)

func TestFscrub_Validate(t *testing.T) {
	crawler := fscrawl.NewCrawler(zap.NewNop())
	log := zap.NewNop()
	type args struct {
		dirs    []model.Directory
		watcher model.Handler
		crawler model.Handler
		actions []model.Action
		log     *zap.Logger
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"basic",
			args{[]model.Directory{"test"}, nil, nil, []model.Action{model.NoOpAction}, log},
			false,
		},
		{
			"withCrawler",
			args{[]model.Directory{"test"}, nil, crawler.Run, []model.Action{model.NoOpAction}, log},
			false,
		},
		{
			"dirEmptyError",
			args{[]model.Directory{}, nil, nil, []model.Action{model.NoOpAction}, log},
			true,
		},
		{
			"dirNilError",
			args{nil, nil, nil, []model.Action{model.NoOpAction}, log},
			true,
		},
		{
			"actionEmptyError",
			args{[]model.Directory{"test"}, nil, nil, []model.Action{}, log},
			true,
		},
		{
			"actionNilError",
			args{[]model.Directory{"test"}, nil, nil, nil, log},
			true,
		},
		{
			"actionNilEntryError",
			args{[]model.Directory{"test"}, nil, nil, []model.Action{nil}, log},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fscrub := NewFscrub(tt.args.dirs, tt.args.watcher, tt.args.crawler, tt.args.actions, tt.args.log)
			if err := fscrub.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Fscrub.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFscrub_Run(t *testing.T) {
	log := zap.NewNop()
	tests := []struct {
		name    string
		f       *Fscrub
		wantErr bool
	}{
		{
			"basic",
			&Fscrub{
				[]model.Directory{"test"},
				mockHandler,
				mockHandler,
				[]model.Action{mockAction},
				log,
			},
			false,
		},
		{
			"error",
			&Fscrub{
				[]model.Directory{""},
				mockHandler,
				mockHandler,
				[]model.Action{mockAction},
				log,
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.f.Run(); (err != nil) != tt.wantErr {
				t.Errorf("Fscrub.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func mockHandler(dir model.Directory, erc chan error) {
	if len(dir) < 1 {
		erc <- errors.New("invalid dir")
	}
	if dir == "fail" {
		erc <- errors.New("handling dir failed")
	}
	time.Sleep(time.Second * 3)
	erc <- nil
}

func mockAction(path string, file os.FileInfo) error {
	if len(path) < 1 {
		return errors.New("invalid path")
	}
	if path == "fail" {
		return errors.New("action failed")
	}
	return nil
}

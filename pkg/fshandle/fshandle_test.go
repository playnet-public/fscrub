package fshandle

import (
	"os"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/playnet-public/fscrub/pkg/fscrawl"
	"github.com/playnet-public/fscrub/pkg/model"
	"go.uber.org/zap"
)

func TestFsHandler_New(t *testing.T) {
	crawler := fscrawl.NewCrawler(zap.NewNop())
	log := zap.NewNop()
	type args struct {
		dirs     model.Directories
		handlers []model.Handler
		log      *zap.Logger
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"basic",
			args{model.Directories{"test"}, []model.Handler{nil, nil}, log},
			false,
		},
		{
			"withCrawler",
			args{model.Directories{"test"}, []model.Handler{nil, crawler}, log},
			false,
		},
		{
			"dirEmptyError",
			args{model.Directories{}, []model.Handler{nil, nil}, log},
			true,
		},
		{
			"dirNilError",
			args{nil, []model.Handler{nil, nil}, log},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fshandle := NewFsHandler(tt.args.dirs, tt.args.handlers, tt.args.log)
			if (fshandle == nil) != tt.wantErr {
				t.Errorf("Fscrub.NewFsHandler() fshandle = %v, want %v", fshandle, tt.wantErr)
			}
		})
	}
}

func TestFsHandler_Run(t *testing.T) {
	log := zap.NewNop()
	tests := []struct {
		name    string
		f       *FsHandler
		wantErr bool
	}{
		{
			"basic",
			NewFsHandler(
				model.Directories{"test"},
				[]model.Handler{
					&mockHandler{},
					&mockHandler{},
				},
				log,
			),
			false,
		},
		{
			"error",
			NewFsHandler(
				model.Directories{""},
				[]model.Handler{
					&mockHandler{},
					&mockHandler{},
				},
				log,
			),
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

type mockHandler struct{}

func (h *mockHandler) Run(dir model.Directory, erc chan error) {
	if len(dir) < 1 {
		erc <- errors.New("invalid dir")
	}
	if dir == "fail" {
		erc <- errors.New("handling dir failed")
	}
	time.Sleep(time.Second * 3)
	erc <- nil
}

func (h *mockHandler) Stop() {}

func mockAction(path string, file os.FileInfo) error {
	if len(path) < 1 {
		return errors.New("invalid path")
	}
	if path == "fail" {
		return errors.New("action failed")
	}
	return nil
}

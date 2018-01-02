package fswatch

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/playnet-public/fscrub/pkg/model"
	"go.uber.org/zap"
)

// TODO: Add tests for file system events
func TestNewWatcher(t *testing.T) {
	log := zap.NewNop()
	tests := []struct {
		name    string
		actions []model.Action
		dirs    model.Directories
		wantNil bool
		wantErr bool
	}{
		{
			"basic",
			[]model.Action{},
			model.Directories{},
			false,
			false,
		},
		{
			"noDir",
			[]model.Action{model.NoOpAction},
			model.Directories{},
			false,
			false,
		},
		{
			"dir",
			[]model.Action{model.NoOpAction},
			model.Directories{""},
			false,
			false,
		},
		{
			"invalidDir",
			[]model.Action{model.NoOpAction},
			model.Directories{"invalidDir"},
			false,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := NewWatcher(log, tt.actions...)
			if (w != nil) == tt.wantNil {
				t.Errorf("NewWatcher() = %v, want nil: %v", w, tt.wantNil)
			}
			erc := make(chan error)
			for _, d := range tt.dirs {
				go w.Run(d, erc)
				select {
				case <-time.After(time.Millisecond * 5):
				case err := <-erc:
					if (err != nil) != tt.wantErr {
						t.Errorf("Watcher.Run() = %v, wantErr %v", err, tt.wantErr)
					}
				}
			}
			w.Stop()
		})
	}
}

func TestWatcher_handle(t *testing.T) {
	tests := []struct {
		name    string
		w       *Watcher
		path    string
		wantErr bool
	}{
		{
			"basic",
			&Watcher{
				log:     zap.NewNop(),
				actions: []model.Action{model.NoOpAction},
			},
			"fswatch_test.go",
			false,
		},
		{
			"fileErr",
			&Watcher{
				log:     zap.NewNop(),
				actions: []model.Action{model.NoOpAction},
			},
			"",
			true,
		},
		{
			"actionErr",
			&Watcher{
				log:     zap.NewNop(),
				actions: []model.Action{errorAction},
			},
			"fswatch_test.go",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.w.handle(tt.path); (err != nil) != tt.wantErr {
				t.Errorf("Watcher.handle() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func errorAction(path string, file os.FileInfo) error {
	return errors.New("testError")
}

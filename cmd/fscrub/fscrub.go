package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/playnet-public/fscrub/pkg/fscrub"

	"github.com/playnet-public/fscrub/pkg/fscrawl"
	"github.com/playnet-public/fscrub/pkg/fshandle"
	"github.com/playnet-public/fscrub/pkg/fslog"

	raven "github.com/getsentry/raven-go"
	"github.com/golang/glog"
	"github.com/kolide/kit/version"
	"github.com/pkg/errors"
	"github.com/playnet-public/fscrub/pkg/model"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	app    = "fscrub"
	appKey = "fscrub"
)

var (
	maxprocsPtr = flag.Int("maxprocs", runtime.NumCPU(), "max go procs")
	sentryDsn   = flag.String("sentrydsn", "", "sentry dsn key")
	dbgPtr      = flag.Bool("debug", false, "debug printing")
	versionPtr  = flag.Bool("version", true, "show or hide version info")
	patternPtr  = flag.String("patterns", "", "path where additional patterns are stored")

	dirs   model.Directories
	sentry *raven.Client
)

func main() {
	flag.Var(&dirs, "dir", "directories to scrub")
	flag.Parse()

	if *versionPtr {
		fmt.Printf("-- PlayNet %s --\n", app)
		version.PrintFull()
	}
	runtime.GOMAXPROCS(*maxprocsPtr)

	// prepare glog
	defer glog.Flush()
	glog.CopyStandardLogTo("info")

	var zapFields []zapcore.Field
	// hide app and version information when debugging
	if !*dbgPtr {
		zapFields = []zapcore.Field{
			zap.String("app", appKey),
			zap.String("version", version.Version().Version),
		}
	}

	// prepare zap logging
	log := newLogger(*dbgPtr).With(zapFields...)
	defer log.Sync()
	log.Info("preparing")

	var err error

	// prepare sentry error logging
	sentry, err = raven.New(*sentryDsn)
	if err != nil {
		panic(err)
	}
	err = raven.SetDSN(*sentryDsn)
	if err != nil {
		panic(err)
	}
	errs := make(chan error)

	// catch system interrupts
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	// catch errors and throw fatal //TODO: is this good?
	go func() {
		ret := <-errs
		if ret != nil {
			log.Fatal(ret.Error())
		}
	}()

	// run main code
	log.Info("starting")
	raven.CapturePanicAndWait(func() {
		if err := do(log); err != nil {
			log.Fatal("fatal error encountered", zap.Error(err))
			raven.CaptureErrorAndWait(err, map[string]string{"isFinal": "true"})
			errs <- err
		}
	}, nil)
	log.Info("finished")
}

func do(log *zap.Logger) error {
	logAction := fslog.NewFsLogger(log)
	patterns, err := parsePatterns(*patternPtr)
	if err != nil {
		return err
	}
	fscrubAction := fscrub.NewFscrub(log, false, patterns...)

	actions := []model.Action{
		logAction.Log,
		fscrubAction.Handle,
	}

	crawler := fscrawl.NewCrawler(log, actions...)

	fshandler := fshandle.NewFsHandler(
		dirs,
		[]model.Handler{
			crawler,
		},
		log,
	)

	err = fshandler.Run()
	if err != nil {
		return errors.Wrap(err, "running fscrub failed")
	}

	return nil
}

func parsePatterns(path string) ([]fscrub.Pattern, error) {
	if path == "" {
		return []fscrub.Pattern{}, nil
	}
	config := &fscrub.PatternConfig{}
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return []fscrub.Pattern{}, err
	}
	if err := json.Unmarshal(content, config); err != nil {
		return []fscrub.Pattern{}, err
	}
	return config.Patterns, nil
}

//TODO: Move this to playnet common libs
func newLogger(dbg bool) *zap.Logger {
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel
	})

	consoleDebugging := zapcore.Lock(os.Stdout)
	consoleErrors := zapcore.Lock(os.Stderr)
	consoleConfig := zap.NewDevelopmentEncoderConfig()
	consoleEncoder := zapcore.NewConsoleEncoder(consoleConfig)
	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, consoleErrors, highPriority),
		zapcore.NewCore(consoleEncoder, consoleDebugging, lowPriority),
	)
	logger := zap.New(core)
	if dbg {
		logger = logger.WithOptions(
			zap.AddCaller(),
			zap.AddStacktrace(zap.ErrorLevel),
		)
	} else {
		logger = logger.WithOptions(
			zap.AddStacktrace(zap.FatalLevel),
		)
	}
	return logger
}

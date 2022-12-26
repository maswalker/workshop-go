package main

import (
	"fmt"
	"os"

	"github.com/maswalker/workshop-go/pkg/log"
	"github.com/samber/lo"
)

func InitLog(debug bool) {
	var logger *log.Logger
	if debug {
		logger = log.New(os.Stderr, log.InfoLevel, log.WithCaller(true), log.AddCallerSkip(1))
	} else {
		tops := lo.Map(lo.RangeWithSteps(log.WarnLevel, log.FatalLevel, 1), func(lvl log.Level, _ int) log.TeeOption {
			return log.TeeOption{
				Filename: fmt.Sprintf("%s.log", lvl.String()),
				Ropt: log.RotateOptions{
					MaxSize:    1,
					MaxAge:     1,
					MaxBackups: 3,
					Compress:   true,
				},
				Lef: func(lvl log.Level) bool {
					return lvl == log.InfoLevel
				},
			}
		})
		logger = log.NewTeeWithRotate(tops, log.WithCaller(true), log.AddCallerSkip(1))
	}
	log.ResetDefault(logger)
}

func main() {
	InitLog(true)
	defer log.Sync()
	log.Info("demo3:", log.String("app", "start ok"),
		log.Int("major version", 3))
	log.Error("demo3:", log.String("app", "crash"),
		log.Int("reason", -1))
}

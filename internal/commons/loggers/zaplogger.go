package loggers

import (
	"log"

	"github.com/swallowarc/tictactoe_battle_backend/internal/commons/mode"
	"go.uber.org/zap"
)

func NewZapLogger(m mode.Mode) *zap.Logger {
	switch m {
	case mode.Test:
		fallthrough
	case mode.Release:
		return getInstance(zap.NewProductionConfig())
	}
	c := zap.NewDevelopmentConfig()
	c.DisableStacktrace = true
	return getInstance(c)
}

// getInstance get zap logger instance.
func getInstance(config zap.Config) *zap.Logger {
	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stderr"}
	zapLogger, err := config.Build()
	if err != nil {
		log.Println("create logger failed.")
		log.Fatal(err)
	}
	return zapLogger
}

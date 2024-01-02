package logger

import (
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// mode is Zap logger mode
type mode int

const (
	ModeProd mode = iota // production mode
	ModeDev              // development mode
)

// SLogger is customized zap.SugaredLogger
type SLogger struct {
	*zap.SugaredLogger
}

var (
	log    *SLogger
	config zap.Config
	once   sync.Once
)

// Log is zap.SugaredLogger singleton
func Log() *SLogger {
	if log == nil {
		SetMode(ModeDev)
	}
	return log
}

// SetMode allows to set up logging mode only once
func SetMode(mode mode) {
	once.Do(func() {
		switch mode {
		case ModeProd:
			config = zap.NewProductionConfig()
		case ModeDev:
			config = zap.NewDevelopmentConfig()
			config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.DateTime + ".000")
			config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
			config.DisableStacktrace = true
		}
		l, _ := config.Build()
		log = &SLogger{SugaredLogger: l.Sugar()}
	})
}

// SetLevel sets global logger level
func SetLevel(level int) {
	zLevel := zapcore.Level(level)
	if zLevel < zapcore.DebugLevel || zLevel > zapcore.FatalLevel {
		zLevel = zapcore.InfoLevel
	}
	config.Level.SetLevel(zLevel)
	l, _ := config.Build()
	log.SugaredLogger = l.Sugar()
}

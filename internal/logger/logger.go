package logger

import (
	"fmt"
	"github.com/acya-skulskaya/go-musthave-diploma/internal/config/logging"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger = zap.NewNop()

func Init(config logging.Config) error {
	lvl, err := zap.ParseAtomicLevel(config.Level)
	if err != nil {
		return fmt.Errorf("could not parse logging level: %w", err)
	}

	cfg := zap.NewProductionConfig()
	//cfg.ErrorOutputPaths = []string{"stderr", "./logs/error.log"}
	//cfg.OutputPaths = []string{"stderr", "./logs/output.log"}
	cfg.Level = lvl
	cfg.EncoderConfig = zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		TimeKey:        "ts",
		FunctionKey:    zapcore.OmitKey,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02T15:04:05.000Z"), // Формат ISO 8601
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
	}

	zl, err := cfg.Build()
	if err != nil {
		return fmt.Errorf("could not build logger: %w", err)
	}
	// устанавливаем синглтон
	Log = zl
	return nil
}

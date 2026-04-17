package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLogger(cfg *Config) (*zap.Logger, error) {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	var cores []zapcore.Core

	if cfg.Log.EnableConsole {
		consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
		consoleWriter := zapcore.AddSync(os.Stdout)
		cores = append(cores, zapcore.NewCore(consoleEncoder, consoleWriter, zapcore.DebugLevel))
	}

	if cfg.Log.EnableFile {
		dir := filepath.Dir(cfg.Log.FilePath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create log directory: %w", err)
		}

		fileEncoder := zapcore.NewJSONEncoder(encoderConfig)

		lumberjackLogger := &lumberjack.Logger{
			Filename:   cfg.Log.FilePath,
			MaxSize:    cfg.Log.MaxSize,
			MaxBackups: cfg.Log.MaxBackups,
			MaxAge:     cfg.Log.MaxAge,
			Compress:   cfg.Log.Compress,
		}

		fileWriter := zapcore.AddSync(lumberjackLogger)
		cores = append(cores, zapcore.NewCore(fileEncoder, fileWriter, zapcore.InfoLevel))
	}

	if len(cores) == 0 {
		return zap.NewNop(), nil
	}

	core := zapcore.NewTee(cores...)
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	return logger, nil
}

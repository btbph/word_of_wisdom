package logger

import (
	"errors"
	"log/slog"
	"os"
)

const (
	DebugLogLevel = "debug"
	InfoLogLevel  = "info"
	WarnLogLevel  = "warn"
	ErrorLogLevel = "error"
)

func ParseLogLevel(loglevel string) (slog.Level, error) {
	var level slog.Level

	switch loglevel {
	case DebugLogLevel:
		level = slog.LevelDebug
	case InfoLogLevel:
		level = slog.LevelInfo
	case WarnLogLevel:
		level = slog.LevelWarn
	case ErrorLogLevel:
		level = slog.LevelError
	default:
		return 0, errors.New("unkown log level")
	}

	return level, nil
}

func Init(logLevel string) (*slog.Logger, error) {
	level, err := ParseLogLevel(logLevel)
	if err != nil {
		return nil, err
	}

	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})), nil
}

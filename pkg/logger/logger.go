package logger

import (
	"log/slog"
	"os"
	"strings"
)

var Log *slog.Logger = slog.Default()

func InitLogger(env string) {
	var handler slog.Handler

	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	if strings.ToLower(env) == "development" {
		opts.Level = slog.LevelDebug
		handler = slog.NewTextHandler(os.Stdout, opts)
	} else {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}

	Log = slog.New(handler)
	slog.SetDefault(Log)
}

func Info(msg string, args ...any) {
	if Log != nil {
		Log.Info(msg, args...)
	}
}

func Error(msg string, args ...any) {
	if Log != nil {
		Log.Error(msg, args...)
	}
}

func Debug(msg string, args ...any) {
	if Log != nil {
		Log.Debug(msg, args...)
	}
}

func Warn(msg string, args ...any) {
	if Log != nil {
		Log.Warn(msg, args...)
	}
}

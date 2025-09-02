package log

import (
	"os"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Options : semua output JSON. Console & file bisa diaktifkan bersamaan.
type Options struct {
	// "debug" | "info" | "warn" | "error" | "dpanic" | "panic" | "fatal"
	Level string

	// Console stdout (JSON)
	ConsoleEnabled bool

	// File (JSON) + rotation; kosong => nonaktif
	FilePath   string
	MaxSizeMB  int  // default 50
	MaxBackups int  // default 5
	MaxAgeDays int  // default 30
	Compress   bool // default true

	// Stacktrace mulai level apa (default "error")
	StacktraceAt string

	// Caller skip (default 1)
	CallerSkip int
}

func New(opts Options) (*zap.Logger, error) {
	// defaults
	if strings.TrimSpace(opts.Level) == "" {
		opts.Level = "info"
	}
	if opts.CallerSkip <= 0 {
		opts.CallerSkip = 1
	}
	if opts.MaxSizeMB <= 0 {
		opts.MaxSizeMB = 50
	}
	if opts.MaxBackups < 0 {
		opts.MaxBackups = 5
	}
	if opts.MaxAgeDays <= 0 {
		opts.MaxAgeDays = 30
	}
	if strings.TrimSpace(opts.StacktraceAt) == "" {
		opts.StacktraceAt = "error"
	}
	if !opts.ConsoleEnabled && opts.FilePath == "" {
		// default minimal: console on
		opts.ConsoleEnabled = true
	}

	lvl := parseLevel(opts.Level)
	stackLvl := parseLevel(opts.StacktraceAt)

	// JSON encoder config dengan timestamp human-readable (ms)
	encCfg := zapcore.EncoderConfig{
		TimeKey:       "ts",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "caller",
		MessageKey:    "msg",
		StacktraceKey: "stack",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			// contoh: 2025-09-02 13:42:01.123
			enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
		},
		EncodeLevel: func(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(strings.ToUpper(l.String()))
		},
		EncodeCaller: zapcore.ShortCallerEncoder,
		EncodeName:   zapcore.FullNameEncoder,
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(d.String())
		},
	}

	var cores []zapcore.Core

	// Console JSON
	if opts.ConsoleEnabled {
		consoleEnc := zapcore.NewJSONEncoder(encCfg)
		consoleWS := zapcore.AddSync(os.Stdout)
		cores = append(cores, zapcore.NewCore(consoleEnc, consoleWS, lvl))
	}

	// File JSON + rotation
	if opts.FilePath != "" {
		lj := &lumberjack.Logger{
			Filename:   opts.FilePath,
			MaxSize:    opts.MaxSizeMB,
			MaxBackups: opts.MaxBackups,
			MaxAge:     opts.MaxAgeDays,
			Compress:   opts.Compress,
		}
		fileWS := zapcore.AddSync(lj)
		fileEnc := zapcore.NewJSONEncoder(encCfg)
		cores = append(cores, zapcore.NewCore(fileEnc, fileWS, lvl))
	}

	core := zapcore.NewTee(cores...)

	logger := zap.New(
		core,
		zap.AddCaller(),
		zap.AddCallerSkip(opts.CallerSkip),
		zap.AddStacktrace(stackLvl),
	)

	return logger, nil
}

func parseLevel(s string) zapcore.LevelEnabler {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn", "warning":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "dpanic":
		return zapcore.DPanicLevel
	case "panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

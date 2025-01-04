package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func newZapLogger(level, output, encoding string, sampleRate int) (*zap.SugaredLogger, error) {
	zapLevel := zap.InfoLevel
	switch level {
	case "debug":
		zapLevel = zap.DebugLevel
	case "info":
		zapLevel = zap.InfoLevel
	case "warn":
		zapLevel = zap.WarnLevel
	case "error":
		zapLevel = zap.ErrorLevel
	}
	development := zapLevel == zap.DebugLevel
	cfg := &zap.Config{
		Level:       zap.NewAtomicLevelAt(zapLevel),
		Development: development,
		Sampling: &zap.SamplingConfig{
			Initial:    sampleRate,
			Thereafter: sampleRate,
		},
		Encoding: encoding,
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    "func",
			MessageKey:     "msg",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{output},
		ErrorOutputPaths: []string{output},
	}
	// AddCaller() will annotate the log with filename and line
	// AddCallerSkip(1) ensures we are not just annotating with the wrapper filename and line
	// https://github.com/uber-go/zap/blob/master/options.go
	l, err := cfg.Build(zap.AddCaller(), zap.AddCallerSkip(1))
	if err != nil {
		return nil, err
	}
	return l.Sugar(), nil
}

func NewZapLogger(opts *Opts) (*zap.SugaredLogger, error) {
	return newZapLogger(opts.Level, opts.Output, opts.Encoding, opts.SampleRate)
}

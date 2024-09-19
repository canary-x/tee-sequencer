package logger

import (
	"io"
	"os"

	"github.com/mdlayher/vsock"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/canary-x/tee-sequencer/internal/config"
)

var logger *ZapLogger

type Logger interface {
	Info(msg string, v ...any)
	Debug(msg string, v ...any)
	Warn(msg string, v ...any)
	Error(msg string, v ...any)
}

type ZapLogger struct {
	zap *zap.SugaredLogger
}

func Init(cfg config.Config) Logger {
	logger = newZapVSockLogger(cfg)
	return logger
}

func Instance() Logger {
	if logger == nil {
		panic("called logger.Instance() but logger was never initialized")
	}
	return logger
}

func (l *ZapLogger) Info(msg string, v ...any) {
	if len(v) > 0 {
		l.zap.Infof(msg, v)
	} else {
		l.zap.Info(msg)
	}
}

func (l *ZapLogger) Debug(msg string, v ...any) {
	if len(v) > 0 {
		l.zap.Debugf(msg, v)
	} else {
		l.zap.Debug(msg)
	}
}

func (l *ZapLogger) Warn(msg string, v ...any) {
	if len(v) > 0 {
		l.zap.Warnf(msg, v)
	} else {
		l.zap.Warn(msg)
	}
}

func (l *ZapLogger) Error(msg string, v ...any) {
	if len(v) > 0 {
		l.zap.Errorf(msg, v)
	} else {
		l.zap.Error(msg)
	}
}

func newZapVSockLogger(cfg config.Config) *ZapLogger {
	var z *zap.Logger
	vsockConn, vsockErr := vsock.Dial(vsock.Host, cfg.VSockPort, &vsock.Config{})
	// no need to vsockConn.Close() as this won't be closed until the application exits
	if vsockErr == nil {
		z = initZapLoggerWithVsock(vsockConn)
	} else {
		z = initZapLogger()
	}
	if vsockErr != nil {
		z.Sugar().Warnf("Failed to establish vsock connection for logs: %v", vsockErr)
		z.Warn("Logs will be streamed to console only")
	}
	return &ZapLogger{zap: z.Sugar()}
}

func initZapConsole() (zapcore.Encoder, zapcore.Core) {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Use the console encoder prints human-readable logs instead of json
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
	core := zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel)
	return consoleEncoder, core
}

func initZapLogger() *zap.Logger {
	_, consoleCore := initZapConsole()
	return zap.New(consoleCore, zap.AddCaller(), zap.AddStacktrace(zapcore.DPanicLevel))
}

// initZapLoggerWithVsock initializes a Zap logger that tees logs to both os.Stdout and the vsock connection.
func initZapLoggerWithVsock(vsockConn io.Writer) *zap.Logger {
	consoleEncoder, consoleCore := initZapConsole()
	fileCore := zapcore.NewCore(consoleEncoder, zapcore.AddSync(vsockConn), zapcore.DebugLevel)
	// duplicate logs across the console and the file socket
	combinedCore := zapcore.NewTee(fileCore, consoleCore)
	return zap.New(combinedCore, zap.AddCaller(), zap.AddStacktrace(zapcore.DPanicLevel))
}

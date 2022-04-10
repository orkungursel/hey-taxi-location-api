package logger

import (
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/orkungursel/hey-taxi-location-api/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	AppName  string
	LogLevel string
	DevMode  bool
	Encoder  string
}

type ILogger interface {
	Debug(args ...interface{})
	Debugf(template string, args ...interface{})
	Info(args ...interface{})
	Infof(template string, args ...interface{})
	Warn(args ...interface{})
	Warnf(template string, args ...interface{})
	WarnMsg(msg string, err error)
	Error(args ...interface{})
	Errorf(template string, args ...interface{})
	Err(msg string, err error)
	DPanic(args ...interface{})
	DPanicf(template string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(template string, args ...interface{})
	Printf(template string, args ...interface{})
	EchoCtx(e echo.Context, start time.Time)
	Sync()
}

type appLogger struct {
	appname     string
	level       string
	devMode     bool
	encoding    string
	sugarLogger *zap.SugaredLogger
	logger      *zap.Logger
}

// New App Logger constructor
func New(c *config.Config) *appLogger {
	loggerConfig := &Config{
		AppName:  c.App.Name,
		LogLevel: "debug",
		Encoder:  "console",
		DevMode:  true,
	}
	if c.IsProduction() {
		loggerConfig.LogLevel = "info"
		loggerConfig.Encoder = "json"
		loggerConfig.DevMode = false
	}

	return (&appLogger{
		appname:  loggerConfig.AppName,
		level:    loggerConfig.LogLevel,
		devMode:  loggerConfig.DevMode,
		encoding: loggerConfig.Encoder,
	}).init()
}

var levels = map[string]zapcore.Level{
	"debug":  zapcore.DebugLevel,
	"info":   zapcore.InfoLevel,
	"warn":   zapcore.WarnLevel,
	"error":  zapcore.ErrorLevel,
	"dpanic": zapcore.DPanicLevel,
	"panic":  zapcore.PanicLevel,
	"fatal":  zapcore.FatalLevel,
}

func (l *appLogger) getLevel() zapcore.Level {
	level, exist := levels[l.level]
	if !exist {
		return zapcore.DebugLevel
	}

	return level
}

// init Init logger
func (l *appLogger) init() *appLogger {
	logLevel := l.getLevel()
	logWriter := zapcore.AddSync(os.Stdout)

	var encoderCfg zapcore.EncoderConfig
	if l.devMode {
		encoderCfg = zap.NewDevelopmentEncoderConfig()
	} else {
		encoderCfg = zap.NewProductionEncoderConfig()
	}

	var encoder zapcore.Encoder
	encoderCfg.NameKey = "[SERVICE]"
	encoderCfg.TimeKey = "[TIME]"
	encoderCfg.LevelKey = "[LEVEL]"
	encoderCfg.CallerKey = "[LINE]"
	encoderCfg.MessageKey = "[MESSAGE]"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderCfg.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderCfg.EncodeCaller = zapcore.ShortCallerEncoder
	encoderCfg.EncodeDuration = zapcore.StringDurationEncoder

	if l.encoding == "console" {
		encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoderCfg.EncodeCaller = zapcore.ShortCallerEncoder
		encoderCfg.ConsoleSeparator = " | "
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
	} else {
		encoderCfg.FunctionKey = "[CALLER]"
		encoderCfg.EncodeName = zapcore.FullNameEncoder
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	}

	core := zapcore.NewCore(encoder, logWriter, zap.NewAtomicLevelAt(logLevel))
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	l.logger = logger
	l.sugarLogger = logger.Sugar()

	if l.appname != "" {
		l.logger = l.logger.Named(l.appname)
		l.sugarLogger = l.sugarLogger.Named(l.appname)
	}

	return l
}

// Sync flushes any buffered log entries
func (l *appLogger) Sync() {
	go func() {
		_ = l.logger.Sync()
	}()
	_ = l.sugarLogger.Sync()
}

// Debug uses fmt.Sprint to construct and log a message.
func (l *appLogger) Debug(args ...interface{}) {
	l.sugarLogger.Debug(args...)
}

// Debugf uses fmt.Sprintf to log a templated message
func (l *appLogger) Debugf(template string, args ...interface{}) {
	l.sugarLogger.Debugf(template, args...)
}

// Info uses fmt.Sprint to construct and log a message
func (l *appLogger) Info(args ...interface{}) {
	l.sugarLogger.Info(args...)
}

// Infof uses fmt.Sprintf to log a templated message.
func (l *appLogger) Infof(template string, args ...interface{}) {
	l.sugarLogger.Infof(template, args...)
}

// Printf uses fmt.Sprintf to log a templated message
func (l *appLogger) Printf(template string, args ...interface{}) {
	l.sugarLogger.Infof(template, args...)
}

// Warn uses fmt.Sprint to construct and log a message.
func (l *appLogger) Warn(args ...interface{}) {
	l.sugarLogger.Warn(args...)
}

// WarnMsg log error message with warn level.
func (l *appLogger) WarnMsg(msg string, err error) {
	l.logger.Warn(msg, zap.String("error", err.Error()))
}

// Warnf uses fmt.Sprintf to log a templated message.
func (l *appLogger) Warnf(template string, args ...interface{}) {
	l.sugarLogger.Warnf(template, args...)
}

// Error uses fmt.Sprint to construct and log a message.
func (l *appLogger) Error(args ...interface{}) {
	l.sugarLogger.Error(args...)
}

// Errorf uses fmt.Sprintf to log a templated message.
func (l *appLogger) Errorf(template string, args ...interface{}) {
	l.sugarLogger.Errorf(template, args...)
}

// Err uses error to log a message.
func (l *appLogger) Err(msg string, err error) {
	l.logger.Error(msg, zap.Error(err))
}

// DPanic uses fmt.Sprint to construct and log a message. In development, the logger then panics. (See DPanicLevel for details.)
func (l *appLogger) DPanic(args ...interface{}) {
	l.sugarLogger.DPanic(args...)
}

// DPanicf uses fmt.Sprintf to log a templated message. In development, the logger then panics. (See DPanicLevel for details.)
func (l *appLogger) DPanicf(template string, args ...interface{}) {
	l.sugarLogger.DPanicf(template, args...)
}

// Panic uses fmt.Sprint to construct and log a message, then panics.
func (l *appLogger) Panic(args ...interface{}) {
	l.sugarLogger.Panic(args...)
}

// Panicf uses fmt.Sprintf to log a templated message, then panics
func (l *appLogger) Panicf(template string, args ...interface{}) {
	l.sugarLogger.Panicf(template, args...)
}

// Fatal uses fmt.Sprint to construct and log a message, then calls os.Exit.
func (l *appLogger) Fatal(args ...interface{}) {
	l.sugarLogger.Fatal(args...)
}

// Fatalf uses fmt.Sprintf to log a templated message, then calls os.Exit.
func (l *appLogger) Fatalf(template string, args ...interface{}) {
	l.sugarLogger.Fatalf(template, args...)
}

// EchoCtx is a middleware that logs the request as it goes in and the response as it goes out.
func (l *appLogger) EchoCtx(c echo.Context, start time.Time) {
	req := c.Request()
	res := c.Response()

	id := req.Header.Get(echo.HeaderXRequestID)
	if id == "" {
		id = res.Header().Get(echo.HeaderXRequestID)
	}

	fields := []zapcore.Field{
		zap.Int("status", res.Status),
		zap.String("latency", time.Since(start).String()),
		zap.String("id", id),
		zap.String("method", req.Method),
		zap.String("uri", req.RequestURI),
		zap.String("host", req.Host),
		zap.String("remote_ip", c.RealIP()),
	}

	n := res.Status
	switch {
	case n >= 500:
		l.logger.Error("Error", fields...)
	case n >= 400:
		l.logger.Warn("Warning", fields...)
	case n >= 300:
		l.logger.Info("Redirecting", fields...)
	default:
		//l.logger.Info("Success", fields...)
	}
}

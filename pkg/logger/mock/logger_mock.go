package mock

import (
	"time"

	"github.com/labstack/echo/v4"
)

type MockLogger struct{}

func NewLoggerMock() *MockLogger { return &MockLogger{} }

func (l *MockLogger) Debug(args ...interface{})                      {}
func (l *MockLogger) Debugf(template string, args ...interface{})    {}
func (l *MockLogger) Info(args ...interface{})                       {}
func (l *MockLogger) Infof(template string, args ...interface{})     {}
func (l *MockLogger) Infow(msg string, keysAndValues ...interface{}) {}
func (l *MockLogger) Printf(template string, args ...interface{})    {}
func (l *MockLogger) Warn(args ...interface{})                       {}
func (l *MockLogger) WarnMsg(msg string, err error)                  {}
func (l *MockLogger) Warnf(template string, args ...interface{})     {}
func (l *MockLogger) Error(args ...interface{})                      {}
func (l *MockLogger) Errorf(template string, args ...interface{})    {}
func (l *MockLogger) Err(msg string, err error)                      {}
func (l *MockLogger) DPanic(args ...interface{})                     {}
func (l *MockLogger) DPanicf(template string, args ...interface{})   {}
func (l *MockLogger) Panic(args ...interface{})                      {}
func (l *MockLogger) Panicf(template string, args ...interface{})    {}
func (l *MockLogger) Fatal(args ...interface{})                      {}
func (l *MockLogger) Fatalf(template string, args ...interface{})    {}
func (l *MockLogger) EchoCtx(c echo.Context, start time.Time)        {}
func (l *MockLogger) Sync()                                          {}

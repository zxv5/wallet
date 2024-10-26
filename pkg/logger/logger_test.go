package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mocking the Logger for testing
type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Debug(args ...interface{}) {
	m.Called(args)
}

func (m *MockLogger) Info(args ...interface{}) {
	m.Called(args)
}

func (m *MockLogger) Warn(args ...interface{}) {
	m.Called(args)
}

func (m *MockLogger) Error(args ...interface{}) {
	m.Called(args)
}

func (m *MockLogger) DPanic(args ...interface{}) {
	m.Called(args)
}

func (m *MockLogger) Panic(args ...interface{}) {
	m.Called(args)
}

func (m *MockLogger) Fatal(args ...interface{}) {
	m.Called(args)
}

// Testing the Init function and logging methods
func TestInit(t *testing.T) {
	conf := &Config{
		EnableConsole: true,
		Levels:        DefaultLevelConf(),
	}

	Init(conf)

	// Assuming Logger is a global variable
	assert.NotNil(t, Logger)
}

func TestSetConfig(t *testing.T) {
	var conf Config
	setConfig(&conf)

	assert.NotNil(t, conf.Levels)
	assert.Equal(t, 3, len(conf.Levels)) // Expecting 3 default levels
}

func TestLoggingMethods(t *testing.T) {
	mockLogger := new(MockLogger)
	// Logger = mockLogger

	mockLogger.On("Debug", "debug message").Return()
	mockLogger.On("Debugf", "debugf message").Return()
	mockLogger.On("Info", "info message").Return()
	mockLogger.On("Infof", "infof message").Return()
	mockLogger.On("Warn", "warn message").Return()
	mockLogger.On("Warnf", "warnf message").Return()
	mockLogger.On("Error", "error message").Return()
	mockLogger.On("Errorf", "errorf message").Return()

	Debug("debug message")
	Debugf("debugf message")
	Info("info message")
	Infof("infof message")
	Warn("warn message")
	Warnf("warnf message")
	Error("error message")
	Errorf("errorf message")
}

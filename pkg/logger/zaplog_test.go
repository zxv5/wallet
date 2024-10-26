package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Mock for lumberjack.Logger
type MockLumberjack struct {
	mock.Mock
}

func (m *MockLumberjack) Write(p []byte) (n int, err error) {
	args := m.Called(p)
	return args.Int(0), args.Error(1)
}

// Mock for ZapBuilder
func TestNewZapBuilder(t *testing.T) {
	conf := []LevelConf{
		{
			Level:      "info",
			Path:       "logs/info.log",
			MaxSize:    10,
			MaxBackups: 5,
			MaxAge:     30,
			Compress:   true,
			LocalTime:  true,
		},
	}

	zapBuilder := NewZapBuilder(conf)

	assert.NotNil(t, zapBuilder)
	assert.Equal(t, true, zapBuilder.consoleEnable)
	assert.Equal(t, zapcore.DebugLevel, zapBuilder.consoleLevel) // Default console level
	assert.Len(t, zapBuilder.hooks, 1)                           // Expecting one hook based on LevelConf
}

func TestZapBuilder_Build(t *testing.T) {
	conf := []LevelConf{
		{
			Level:      "info",
			Path:       "logs/info.log",
			MaxSize:    10,
			MaxBackups: 5,
			MaxAge:     30,
			Compress:   true,
			LocalTime:  true,
		},
	}

	zapBuilder := NewZapBuilder(conf)
	logger := zapBuilder.Build()

	assert.NotNil(t, logger)
	assert.IsType(t, &zap.SugaredLogger{}, logger.Sugar())
}

func TestZapBuilder_EnableConsole(t *testing.T) {
	zapBuilder := NewZapBuilder(nil)
	zapBuilder.DisableFile().EnableConsole()

	logger := zapBuilder.Build()

	assert.NotNil(t, logger)
}

func TestZapBuilder_DisableFile(t *testing.T) {
	zapBuilder := NewZapBuilder(nil)
	zapBuilder.DisableFile()

	logger := zapBuilder.Build()

	assert.NotNil(t, logger)
}

func TestSetLogger(t *testing.T) {
	logger := zap.NewNop()
	SetLogger(logger)

	assert.NotNil(t, Logger)
	assert.Equal(t, logger, Desugar)
}

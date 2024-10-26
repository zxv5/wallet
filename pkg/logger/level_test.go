package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

func TestDefaultLevelConf(t *testing.T) {
	expected := []LevelConf{
		{
			Level:      "debug",
			Path:       "logs/debug.log",
			MaxSize:    50,
			MaxBackups: 30,
			MaxAge:     30,
			Compress:   true,
			LocalTime:  false,
		},
		{
			Level:      "info",
			Path:       "logs/info.log",
			MaxSize:    50,
			MaxBackups: 30,
			MaxAge:     30,
			Compress:   true,
			LocalTime:  false,
		},
		{
			Level:      "warn",
			Path:       "logs/error.log",
			MaxSize:    50,
			MaxBackups: 30,
			MaxAge:     360,
			Compress:   true,
			LocalTime:  false,
		},
	}

	actual := DefaultLevelConf()
	assert.Equal(t, expected, actual)
}

func TestZapLevel(t *testing.T) {
	tests := []struct {
		level string
		want  zapcore.LevelEnabler
	}{
		{"debug", zapcore.DebugLevel},
		{"info", zapcore.InfoLevel},
		{"warn", zapcore.WarnLevel},
		{"error", zapcore.ErrorLevel},
		{"dpanic", zapcore.DPanicLevel},
		{"panic", zapcore.PanicLevel},
		{"fatal", zapcore.FatalLevel},
		{"unknown", zapcore.InfoLevel}, // default case
	}

	for _, tt := range tests {
		t.Run(tt.level, func(t *testing.T) {
			got := zapLevel(tt.level)
			assert.Equal(t, tt.want, got)
		})
	}
}

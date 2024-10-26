package logger

import "go.uber.org/zap/zapcore"

type LevelConf struct {
	Level      string // level
	Path       string // The path of the log file
	MaxSize    int    // The maximum size of each log file saved unit：M
	MaxBackups int    // How many backups can be saved for log files at most
	MaxAge     int    // How many days can the file be saved
	Compress   bool   // Whether to compress
	LocalTime  bool
}

func DefaultLevelConf() []LevelConf {
	return []LevelConf{
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
}

func zapLevel(l string) zapcore.LevelEnabler {
	switch l {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
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

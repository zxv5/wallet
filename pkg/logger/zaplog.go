package logger

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// hooksConfig .
type hooksConfig struct {
	hook  lumberjack.Logger
	level zapcore.LevelEnabler
}

// ZapBuilder .
type ZapBuilder struct {
	hooks         []hooksConfig
	fileEnable    bool
	consoleEnable bool
	consoleLevel  zapcore.LevelEnabler
}

// NewZapBuilder .
func NewZapBuilder(confs []LevelConf) *ZapBuilder {
	if confs == nil {
		confs = DefaultLevelConf()
	}
	var hooks []hooksConfig
	for _, conf := range confs {
		hooks = append(hooks, hooksConfig{
			hook: lumberjack.Logger{
				Filename:   conf.Path,
				MaxSize:    conf.MaxSize,
				MaxBackups: conf.MaxBackups,
				MaxAge:     conf.MaxAge,
				Compress:   conf.Compress,
				LocalTime:  conf.LocalTime,
			},
			level: zapLevel(conf.Level),
		})
	}

	return &ZapBuilder{
		hooks:         hooks,
		fileEnable:    false,
		consoleEnable: true,
		consoleLevel:  zapcore.DebugLevel,
	}
}

// DisableFile .
func (opts *ZapBuilder) DisableFile() *ZapBuilder {
	opts.fileEnable = false
	return opts
}

// EnableConsole .
func (opts *ZapBuilder) EnableConsole() *ZapBuilder {
	opts.consoleEnable = true
	return opts
}

// Build .
func (zb *ZapBuilder) Build() *zap.Logger {
	var (
		cores []zapcore.Core
		opts  []zap.Option
	)

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "name",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		FunctionKey:    "function",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder, // Uppercase encoder
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}

	if zb.fileEnable {
		for i, l := 0, len(zb.hooks); i < l; i++ {
			hook := &zb.hooks[i].hook
			fileW := zapcore.AddSync(hook)
			core := zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), fileW, zb.hooks[i].level)
			cores = append(cores, core)
		}
	}

	if zb.consoleEnable {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
		}
		consoleW := zapcore.Lock(os.Stdout)
		core := zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), consoleW, zb.consoleLevel)
		cores = append(cores, core)
	}
	// else if zb.serverName != "" {
	// 	// Add custom fields serverName
	// 	opts = append(opts, zap.Fields(zap.String("ServerName", zb.serverName)))
	// }

	// Open the file and line number
	opts = append(opts, zap.AddCaller(), zap.AddCallerSkip(1))

	core := zapcore.NewTee(cores...)
	return zap.New(core, opts...)
}

var (
	// Desugar .
	Desugar *zap.Logger
	// Logger .
	Logger *zap.SugaredLogger
)

// SetLogger .
func SetLogger(logger *zap.Logger) {
	Desugar = logger
	Logger = logger.Sugar()
}

func init() {
	Desugar = NewZapBuilder(nil).Build()
	Logger = Desugar.Sugar()
}

package logger

type Config struct {
	EnableConsole bool
	Levels        []LevelConf
}

func setConfig(conf *Config) {
	if conf.Levels == nil {
		conf.Levels = DefaultLevelConf()
	}
}

func Init(conf *Config) {
	if conf == nil {
		conf = &Config{}
	}
	setConfig(conf)

	zapBuilder := NewZapBuilder(conf.Levels)
	if conf.EnableConsole {
		zapBuilder.EnableConsole()
	}

	logger := zapBuilder.Build()
	SetLogger(logger)
}

func Debug(args ...interface{}) {
	Logger.Debug(args...)
}
func Debugf(template string, args ...interface{}) {
	Logger.Debugf(template, args...)
}

func Info(args ...interface{}) {
	Logger.Info(args...)
}
func Infof(template string, args ...interface{}) {
	Logger.Infof(template, args...)
}

func Warn(args ...interface{}) {
	Logger.Warn(args...)
}
func Warnf(template string, args ...interface{}) {
	Logger.Warnf(template, args...)
}

func Error(args ...interface{}) {
	Logger.Error(args...)
}
func Errorf(template string, args ...interface{}) {
	Logger.Errorf(template, args...)
}

func DPanic(args ...interface{}) {
	Logger.DPanic(args...)
}
func DPanicf(template string, args ...interface{}) {
	Logger.DPanicf(template, args...)
}

func Panic(args ...interface{}) {
	Logger.Panic(args...)
}
func Panicf(template string, args ...interface{}) {
	Logger.Panicf(template, args...)
}

// Fatal logs a message, then calls os.Exit(1).
func Fatal(args ...interface{}) {
	Logger.Fatal(args...)
}
func Fatalf(template string, args ...interface{}) {
	Logger.Fatalf(template, args...)
}

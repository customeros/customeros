package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
)

type Config struct {
	LogLevel string `env:"LOGGER_LEVEL" envDefault:"info"`
	DevMode  bool   `env:"LOGGER_DEV_MODE" envDefault:"false"`
	Encoder  string `env:"LOGGER_ENCODER" envDefault:"console"`
}

// Logger methods interface
type Logger interface {
	InitLogger()
	Sync() error
	Logger() *zap.Logger
	SugarLogger() *zap.SugaredLogger
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
	WithName(name string)
}

// Application logger
type AppLogger struct {
	Level       string
	DevMode     bool
	Encoding    string
	sugarLogger *zap.SugaredLogger
	logger      *zap.Logger
}

// NewAppLogger App Logger constructor
func NewAppLogger(cfg *Config) *AppLogger {
	return &AppLogger{Level: strings.ToLower(cfg.LogLevel), DevMode: cfg.DevMode, Encoding: cfg.Encoder}
}

// For mapping config logger to email_service logger levels
var loggerLevelMap = map[string]zapcore.Level{
	"debug":  zapcore.DebugLevel,
	"info":   zapcore.InfoLevel,
	"warn":   zapcore.WarnLevel,
	"error":  zapcore.ErrorLevel,
	"dpanic": zapcore.DPanicLevel,
	"panic":  zapcore.PanicLevel,
	"fatal":  zapcore.FatalLevel,
}

func (l *AppLogger) getLoggerLevel() zapcore.Level {
	level, exist := loggerLevelMap[l.Level]
	if !exist {
		return zapcore.DebugLevel
	}

	return level
}

// InitLogger Init logger
func (l *AppLogger) InitLogger() {
	logLevel := l.getLoggerLevel()

	logWriter := zapcore.AddSync(os.Stdout)

	var encoderCfg zapcore.EncoderConfig
	if l.DevMode {
		encoderCfg = zap.NewDevelopmentEncoderConfig()
	} else {
		encoderCfg = zap.NewProductionEncoderConfig()
	}

	var encoder zapcore.Encoder
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderCfg.EncodeDuration = zapcore.StringDurationEncoder
	encoderCfg.EncodeLevel = zapcore.CapitalLevelEncoder

	if l.Encoding == "console" {
		if l.DevMode {
			encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
		}
		encoderCfg.ConsoleSeparator = " | "
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
	} else {
		encoderCfg.NameKey = "[SERVICE]"
		encoderCfg.TimeKey = "[TIME]"
		encoderCfg.LevelKey = "[LEVEL]"
		encoderCfg.CallerKey = "[LINE]"
		encoderCfg.MessageKey = "[MESSAGE]"
		encoderCfg.FunctionKey = "[CALLER]"
		encoderCfg.StacktraceKey = "[STACKTRACE]"
		encoderCfg.EncodeName = zapcore.FullNameEncoder
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	}

	core := zapcore.NewCore(encoder, logWriter, zap.NewAtomicLevelAt(logLevel))
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	l.logger = logger
	l.sugarLogger = logger.Sugar()

	zap.ReplaceGlobals(logger)
}

// Logger methods

// WithName add logger microservice name
func (l *AppLogger) WithName(name string) {
	l.logger = l.logger.Named(name)
	l.sugarLogger = l.logger.Sugar()
	zap.ReplaceGlobals(l.logger)
}

// Debug uses fmt.Sprint to construct and log a message.
func (l *AppLogger) Debug(args ...interface{}) {
	l.sugarLogger.Debug(args...)
}

// Debugf uses fmt.Sprintf to log a templated message
func (l *AppLogger) Debugf(template string, args ...interface{}) {
	l.sugarLogger.Debugf(template, args...)
}

// Info uses fmt.Sprint to construct and log a message
func (l *AppLogger) Info(args ...interface{}) {
	l.sugarLogger.Info(args...)
}

// Infof uses fmt.Sprintf to log a templated message.
func (l *AppLogger) Infof(template string, args ...interface{}) {
	l.sugarLogger.Infof(template, args...)
}

// Printf uses fmt.Sprintf to log a templated message
func (l *AppLogger) Printf(template string, args ...interface{}) {
	l.sugarLogger.Infof(template, args...)
}

// Warn uses fmt.Sprint to construct and log a message.
func (l *AppLogger) Warn(args ...interface{}) {
	l.sugarLogger.Warn(args...)
}

// WarnMsg log error message with warn Level.
func (l *AppLogger) WarnMsg(msg string, err error) {
	l.logger.Warn(msg, zap.String("error", err.Error()))
}

// Warnf uses fmt.Sprintf to log a templated message.
func (l *AppLogger) Warnf(template string, args ...interface{}) {
	l.sugarLogger.Warnf(template, args...)
}

// Error uses fmt.Sprint to construct and log a message.
func (l *AppLogger) Error(args ...interface{}) {
	l.sugarLogger.Error(args...)
}

// Errorf uses fmt.Sprintf to log a templated message.
func (l *AppLogger) Errorf(template string, args ...interface{}) {
	l.sugarLogger.Errorf(template, args...)
}

// Err uses error to log a message.
func (l *AppLogger) Err(msg string, err error) {
	l.logger.Error(msg, zap.Error(err))
}

// DPanic uses fmt.Sprint to construct and log a message. In development, the logger then panics. (See DPanicLevel for details.)
func (l *AppLogger) DPanic(args ...interface{}) {
	l.sugarLogger.DPanic(args...)
}

// DPanicf uses fmt.Sprintf to log a templated message. In development, the logger then panics. (See DPanicLevel for details.)
func (l *AppLogger) DPanicf(template string, args ...interface{}) {
	l.sugarLogger.DPanicf(template, args...)
}

// Panic uses fmt.Sprint to construct and log a message, then panics.
func (l *AppLogger) Panic(args ...interface{}) {
	l.sugarLogger.Panic(args...)
}

// Panicf uses fmt.Sprintf to log a templated message, then panics
func (l *AppLogger) Panicf(template string, args ...interface{}) {
	l.sugarLogger.Panicf(template, args...)
}

// Fatal uses fmt.Sprint to construct and log a message, then calls os.Exit.
func (l *AppLogger) Fatal(args ...interface{}) {
	l.sugarLogger.Fatal(args...)
}

// Fatalf uses fmt.Sprintf to log a templated message, then calls os.Exit.
func (l *AppLogger) Fatalf(template string, args ...interface{}) {
	l.sugarLogger.Fatalf(template, args...)
}

// Sync flushes any buffered log entries
func (l *AppLogger) Sync() error {
	go l.logger.Sync() // nolint: errcheck
	return l.sugarLogger.Sync()
}

func (l *AppLogger) SugarLogger() *zap.SugaredLogger {
	return l.sugarLogger
}

func (l *AppLogger) Logger() *zap.Logger {
	return l.logger
}

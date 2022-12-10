package logger

var (
	Logger LoggerInterface
)

// Colors
const (
	Reset   = "\033[0m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Magenta = "\033[35m"
)

// LogLevel log level
type LogLevel int

const (
	// Silent silent log level
	Silent LogLevel = iota + 1
	// Error error log level
	Error
	// Warn warn log level
	Warn
	// Info info log level
	Info
)

// Writer log writer interface
type Writer interface {
	Printf(string, ...interface{})
}

// Config logger config
type Config struct {
	Colorful                  bool
	IgnoreRecordNotFoundError bool
	LogLevel                  LogLevel
}

type LoggerInterface interface {
	LogMode(LogLevel) LoggerInterface
	Info(string, ...interface{})
	Warn(string, ...interface{})
	Error(string, ...interface{})
}

func New(writer Writer, config Config) LoggerInterface {
	var (
		infoStr = "INFO "
		warnStr = "WARN "
		errStr  = "ERROR "
	)

	if config.Colorful {
		infoStr = Green + "INFO " + Reset
		warnStr = Magenta + "WARN " + Reset
		errStr = Red + "ERROR " + Reset
	}

	return &logger{
		Writer:  writer,
		Config:  config,
		infoStr: infoStr,
		warnStr: warnStr,
		errStr:  errStr,
	}
}

type logger struct {
	Writer
	Config
	infoStr, warnStr, errStr string
}

func (l *logger) LogMode(level LogLevel) LoggerInterface {
	newlogger := *l
	newlogger.LogLevel = level
	return &newlogger
}

func (l logger) Info(msg string, data ...interface{}) {
	if l.LogLevel >= Info {
		l.Printf(l.infoStr+msg, data...)
	}
}

func (l logger) Warn(msg string, data ...interface{}) {
	if l.LogLevel >= Warn {
		l.Printf(l.warnStr+msg, data...)
	}
}

func (l logger) Error(msg string, data ...interface{}) {
	if l.LogLevel >= Error {
		l.Printf(l.errStr+msg, data...)
	}
}

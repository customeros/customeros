package config

import "github.com/sirupsen/logrus"

func InitLogger(cfg *Config) {
	switch cfg.LogLevel {
	case "PANIC":
		logrus.SetLevel(logrus.PanicLevel)
	case "FATAL":
		logrus.SetLevel(logrus.FatalLevel)
	case "ERROR":
		logrus.SetLevel(logrus.ErrorLevel)
	case "WARN":
		logrus.SetLevel(logrus.WarnLevel)
	case "DEBUG":
		logrus.SetLevel(logrus.DebugLevel)
	case "TRACE":
		logrus.SetLevel(logrus.TraceLevel)
	case "INFO":
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}
}

package utils

import (
	"fmt"
)

const timeFormat = "2006-01-02 15:04:05.000"

type BoltLogger interface {
	LogClientMessage(context string, msg string, args ...any)
	LogServerMessage(context string, msg string, args ...any)
}

type ConsoleBoltNoLogger struct {
}

func (cbl *ConsoleBoltNoLogger) LogClientMessage(id, msg string, args ...any) {
	//cbl.logBoltMessage("C", id, msg, args)
}

func (cbl *ConsoleBoltNoLogger) LogServerMessage(id, msg string, args ...any) {
	//cbl.logBoltMessage("S", id, msg, args)
}

func (cbl *ConsoleBoltNoLogger) logBoltMessage(src, id string, msg string, args []any) {
	//_, _ = fmt.Fprintf(os.Stdout, "%s   BOLT  %s%s: %s\n", time.Now().Format(timeFormat), formatId(id), src, fmt.Sprintf(msg, args...))
}

func formatId(id string) string {
	if id == "" {
		return ""
	}
	return fmt.Sprintf("[%s] ", id)
}

func ConsoleBoltNoLoggerrr() *ConsoleBoltNoLogger {
	return &ConsoleBoltNoLogger{}
}

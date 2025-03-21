package utils

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

func ConsoleBoltNoLoggerrr() *ConsoleBoltNoLogger {
	return &ConsoleBoltNoLogger{}
}

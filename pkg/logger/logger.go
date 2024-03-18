package logger

type Logger interface {
	Debugw(string, ...interface{})
	Infow(string, ...interface{})
	Errorw(string, ...interface{})
}

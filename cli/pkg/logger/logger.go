package logger

type Logger interface {
	Debugf(template string, args ...interface{})
	Infof(template string, args ...interface{})
	Warnf(template string, args ...interface{})
	Errorf(template string, args ...interface{})
	Fatalf(template string, args ...interface{})
	Debugw(msg string, kvs ...interface{})
	Infow(msg string, kvs ...interface{})
	Warnw(msg string, kvs ...interface{})
	Errorw(msg string, kvs ...interface{})
	Fatalw(msg string, kvs ...interface{})
}

var logger Logger

func init() {
	logger = newStdLogger()
}

func Set(l Logger) {
	logger = l
}

func Debugf(template string, args ...interface{}) {
	logger.Debugf(template, args...)
}

func Infof(template string, args ...interface{}) {
	logger.Infof(template, args...)
}

func Warnf(template string, args ...interface{}) {
	logger.Warnf(template, args...)
}

func Errorf(template string, args ...interface{}) {
	logger.Errorf(template, args...)
}

func Fatalf(template string, args ...interface{}) {
	logger.Fatalf(template, args...)
}

func Debugw(msg string, kvs ...interface{}) {
	logger.Debugw(msg, kvs...)
}

func Infow(msg string, kvs ...interface{}) {
	logger.Infow(msg, kvs...)
}

func Warnw(msg string, kvs ...interface{}) {
	logger.Warnw(msg, kvs...)
}

func Errorw(msg string, kvs ...interface{}) {
	logger.Errorw(msg, kvs...)
}

func Fatalw(msg string, kvs ...interface{}) {
	logger.Fatalw(msg, kvs...)
}

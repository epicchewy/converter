package logger

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type stdLogger struct{}

func newStdLogger() Logger {
	return &stdLogger{}
}

func (l *stdLogger) Debugf(template string, args ...interface{}) {
	log.Printf(template, args...)
}

func (l *stdLogger) Infof(template string, args ...interface{}) {
	log.Printf(template, args...)
}

func (l *stdLogger) Warnf(template string, args ...interface{}) {
	log.Printf(template, args...)
}

func (l *stdLogger) Errorf(template string, args ...interface{}) {
	log.Printf(template, args...)
}

func (l *stdLogger) Fatalf(template string, args ...interface{}) {
	log.Printf(template, args...)
	os.Exit(1)
}

func (l *stdLogger) Debugw(msg string, kvs ...interface{}) {
	fields := createFieldsFromKvs(kvs)
	log.Printf("%s %s", msg, fields)
}

func (l *stdLogger) Infow(msg string, kvs ...interface{}) {
	fields := createFieldsFromKvs(kvs)
	log.Printf("%s %s", msg, fields)
}

func (l *stdLogger) Warnw(msg string, kvs ...interface{}) {
	fields := createFieldsFromKvs(kvs)
	log.Printf("%s %s", msg, fields)
}

func (l *stdLogger) Errorw(msg string, kvs ...interface{}) {
	fields := createFieldsFromKvs(kvs)
	log.Printf("%s %s", msg, fields)
}

func (l *stdLogger) Fatalw(msg string, kvs ...interface{}) {
	fields := createFieldsFromKvs(kvs)
	log.Printf("%s %s", msg, fields)
	os.Exit(1)
}

// Create a formatted string from the list of key value pairs passed into the
// log messages. These pairs will be in a format of key=value
func createFieldsFromKvs(kvs []interface{}) string {
	pairs := make([]string, len(kvs)/2+1)
	for i := 0; i < len(kvs); i += 2 {
		var value interface{}
		if len(kvs) > i+1 {
			value = kvs[i+1]
		}
		pairs[i/2] = fmt.Sprintf("%v=%v", kvs[i], value)
	}
	return strings.Join(pairs, " ")
}

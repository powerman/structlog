package structlog

import (
	"io"
	"log"
)

func SetLogFormat(format logFormat) *Logger {
	DefaultLogger = DefaultLogger.SetLogFormat(format)
	return DefaultLogger
}

func SetLogLevel(level logLevel) *Logger {
	DefaultLogger = DefaultLogger.SetLogLevel(level)
	return DefaultLogger
}

func SetKeyValFormat(format string) *Logger {
	DefaultLogger = DefaultLogger.SetKeyValFormat(format)
	return DefaultLogger
}

func SetTimeFormat(format string) *Logger {
	DefaultLogger = DefaultLogger.SetTimeFormat(format)
	return DefaultLogger
}

func SetTimeValFormat(format string) *Logger {
	DefaultLogger = DefaultLogger.SetTimeValFormat(format)
	return DefaultLogger
}

func AddCallDepth(depth int) *Logger {
	DefaultLogger = DefaultLogger.AddCallDepth(depth)
	return DefaultLogger
}

func SetDefaultKeyvals(keyvals ...interface{}) *Logger {
	DefaultLogger = DefaultLogger.SetDefaultKeyvals(keyvals...)
	return DefaultLogger
}

func SetPrefixKeys(keys ...string) *Logger {
	DefaultLogger = DefaultLogger.SetPrefixKeys(keys...)
	return DefaultLogger
}

func AppendPrefixKeys(keys ...string) *Logger {
	DefaultLogger = DefaultLogger.AppendPrefixKeys(keys...)
	return DefaultLogger
}

func PrependSuffixKeys(keys ...string) *Logger {
	DefaultLogger = DefaultLogger.PrependSuffixKeys(keys...)
	return DefaultLogger
}

func SetSuffixKeys(keys ...string) *Logger {
	DefaultLogger = DefaultLogger.SetSuffixKeys(keys...)
	return DefaultLogger
}

func SetKeysFormat(keysFormat map[string]string) *Logger {
	DefaultLogger = DefaultLogger.SetKeysFormat(keysFormat)
	return DefaultLogger
}

func IsInfo() bool {
	return DefaultLogger.IsInfo()
}

func IsDebug() bool {
	return DefaultLogger.IsDebug()
}

func Recover(err *error, keyvals ...interface{}) {
	DefaultLogger.Recover(err, keyvals...)
}

func Err(msg interface{}, keyvals ...interface{}) error {
	return DefaultLogger.Err(msg, keyvals...)
}

func Warn(msg interface{}, keyvals ...interface{}) {
	DefaultLogger.Warn(msg, keyvals...)
}

func Info(msg string, keyvals ...interface{}) {
	DefaultLogger.Info(msg, keyvals...)
}

func Debug(msg interface{}, keyvals ...interface{}) {
	DefaultLogger.Debug(msg, keyvals...)
}

func Print(v ...interface{}) {
	DefaultLogger.Print(v...)
}

func Printf(format string, v ...interface{}) {
	DefaultLogger.Printf(format, v...)
}

func Println(v ...interface{}) {
	DefaultLogger.Println(v...)
}

func Fatal(v ...interface{}) {
	DefaultLogger.Fatal(v...)
}

func Fatalf(format string, v ...interface{}) {
	DefaultLogger.Fatalf(format, v...)
}

func Fatalln(v ...interface{}) {
	DefaultLogger.Fatalln(v...)
}

func Panic(v ...interface{}) {
	DefaultLogger.Panic(v...)
}

func Panicf(format string, v ...interface{}) {
	DefaultLogger.Panicf(format, v...)
}

func Panicln(v ...interface{}) {
	DefaultLogger.Panicln(v...)
}

func Error(msg string, keyvals ...interface{}) {
	DefaultLogger.Error(msg, keyvals...)
}

func Errorf(format string, v ...interface{}) {
	DefaultLogger.Errorf(format, v...)
}

func Infof(format string, v ...interface{}) {
	DefaultLogger.Infof(format, v...)
}

func Debugf(format string, v ...interface{}) {
	DefaultLogger.Debugf(format, v...)
}

func WithField(key string, value interface{}) *Logger {
	return DefaultLogger.WithField(key, value)
}

func WithFields(fields map[string]interface{}) *Logger {
	return DefaultLogger.WithFields(fields)
}

func SetOutput(w io.Writer) {
	log.SetOutput(w)
}

func SetFormatter(format logFormat) *Logger {
	DefaultLogger = DefaultLogger.SetFormatter(format)
	return DefaultLogger
}

func SetLevel(level logLevel) *Logger {
	DefaultLogger = DefaultLogger.SetLogLevel(level)
	return DefaultLogger
}

func StandardLogger() *Logger {
	return New()
}

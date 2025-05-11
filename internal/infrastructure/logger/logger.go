package logger

import (
	"io"

	"github.com/sirupsen/logrus"
)

// Logger define a interface que sua aplicação usa.
type Logger interface {
    Debug(args ...interface{})
    Info(args ...interface{})
    Warn(args ...interface{})
    Error(args ...interface{})
    Fatal(args ...interface{})
    WithField(key string, value interface{}) Logger
}

// logrusLogger implementa Logger usando logrus.Logger.
type logrusLogger struct {
    *logrus.Entry
}

// New cria e configura um Logger,
// level: "debug", "info", "warn", "error", etc.
// out: onde escrever (stdout, arquivo, etc.).
func New(level, format string, out io.Writer) Logger {
    l := logrus.New()
    l.SetOutput(out)
    // nivel de logging
    lvl, err := logrus.ParseLevel(level)
    if err != nil {
        lvl = logrus.InfoLevel
    }
    l.SetLevel(lvl)

    // formato: "json" ou "text"
    if format == "json" {
        l.SetFormatter(&logrus.JSONFormatter{})
    } else {
        l.SetFormatter(&logrus.TextFormatter{
            FullTimestamp: true,
        })
    }

    return &logrusLogger{Entry: logrus.NewEntry(l)}
}

func (l *logrusLogger) WithField(key string, value interface{}) Logger {
    return &logrusLogger{Entry: l.Entry.WithField(key, value)}
}

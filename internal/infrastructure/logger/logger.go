package logger

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

// Logger é a interface que sua aplicação usa.
type Logger interface {
    Debug(args ...interface{})
    Debugf(format string, args ...interface{})
    Info(args ...interface{})
    Infof(format string, args ...interface{})
    Warn(args ...interface{})
    Warnf(format string, args ...interface{})
    Error(args ...interface{})
    Errorf(format string, args ...interface{})
    Fatal(args ...interface{})
    Fatalf(format string, args ...interface{})
    WithField(key string, value interface{}) Logger
}

// logrusLogger implementa Logger usando logrus.Entry.
type logrusLogger struct {
    *logrus.Entry
}

func (l *logrusLogger) Debug(args ...interface{}) {
    l.Entry.Debug(args...)
}
func (l *logrusLogger) Debugf(format string, args ...interface{}) {
    l.Entry.Debugf(format, args...)
}

func (l *logrusLogger) Info(args ...interface{}) {
    l.Entry.Info(args...)
}
func (l *logrusLogger) Infof(format string, args ...interface{}) {
    l.Entry.Infof(format, args...)
}

func (l *logrusLogger) Warn(args ...interface{}) {
    l.Entry.Warn(args...)
}
func (l *logrusLogger) Warnf(format string, args ...interface{}) {
    l.Entry.Warnf(format, args...)
}

func (l *logrusLogger) Error(args ...interface{}) {
    l.Entry.Error(args...)
}
func (l *logrusLogger) Errorf(format string, args ...interface{}) {
    l.Entry.Errorf(format, args...)
}

func (l *logrusLogger) Fatal(args ...interface{}) {
    l.Entry.Fatal(args...)
}
func (l *logrusLogger) Fatalf(format string, args ...interface{}) {
    l.Entry.Fatalf(format, args...)
}

func (l *logrusLogger) WithField(key string, value interface{}) Logger {
    return &logrusLogger{Entry: l.Entry.WithField(key, value)}
}

// New cria e configura um Logger:
//  - level: "debug", "info", "warn", "error"
//  - format: "json" ou "text"
//  - out: io.Writer (os.Stdout, arquivo, etc.)
func New(level, format string, out io.Writer) Logger {
    base := logrus.New()
    base.SetOutput(out)

    lvl, err := logrus.ParseLevel(level)
    if err != nil {
        lvl = logrus.InfoLevel
    }
    base.SetLevel(lvl)

    if format == "json" {
        base.SetFormatter(&logrus.JSONFormatter{})
    } else {
        base.SetFormatter(&logrus.TextFormatter{
            FullTimestamp:    true,
            TimestampFormat:  "2006-01-02 15:04:05",
            ForceColors:      true,
            PadLevelText:     true,
            QuoteEmptyFields: true,
        })
    }

    return &logrusLogger{Entry: logrus.NewEntry(base)}
}

// NewDefault é um helper para usar stdout, nível info e texto colorido
func NewDefault() Logger {
    return New("info", "text", os.Stdout)
}

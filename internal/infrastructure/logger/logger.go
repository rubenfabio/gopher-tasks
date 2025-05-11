package logger

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

// Logger é a interface que sua aplicação usa.
type Logger interface {
    Debug(args ...interface{})
    Info(args ...interface{})
    Warn(args ...interface{})
    Error(args ...interface{})
    Fatal(args ...interface{})
    WithField(key string, value interface{}) Logger
}

// logrusLogger implementa Logger usando logrus.Entry.
type logrusLogger struct {
    *logrus.Entry
}

// Debug chama Entry.Debug
func (l *logrusLogger) Debug(args ...interface{}) {
    l.Entry.Debug(args...)
}

// Info chama Entry.Info
func (l *logrusLogger) Info(args ...interface{}) {
    l.Entry.Info(args...)
}

// Warn chama Entry.Warn
func (l *logrusLogger) Warn(args ...interface{}) {
    l.Entry.Warn(args...)
}

// Error chama Entry.Error
func (l *logrusLogger) Error(args ...interface{}) {
    l.Entry.Error(args...)
}

// Fatal chama Entry.Fatal
func (l *logrusLogger) Fatal(args ...interface{}) {
    l.Entry.Fatal(args...)
}

// WithField adiciona um campo e retorna um novo Logger
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

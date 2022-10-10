//go:build !windows && !nacl && !plan9
// +build !windows,!nacl,!plan9

package syslog

import (
	"fmt"
	"log/syslog"
	"os"

	"github.com/sirupsen/logrus"
)

// Sink to send logs via syslog.
type Sink struct {
	Writer    *syslog.Writer
	Formatter logrus.Formatter
}

// NewSink Creates a hook to be added to an instance of logger. This is called with
// `sink, err := NewSink("udp", "localhost:514", syslog.LOG_DEBUG, "")`
func NewSink(f logrus.Formatter, network, raddr string, priority syslog.Priority, tag string) (*Sink, error) {
	w, err := syslog.Dial(network, raddr, priority, tag)
	return &Sink{w, f}, err
}

// NewSinkWriter Creates a hook to be added to an instance of logger. This is called with
func NewSinkWriter(f logrus.Formatter, w *syslog.Writer) *Sink {
	return &Sink{w, f}
}

func (sink *Sink) Emit(entry *logrus.Entry) error {
	line, err := sink.Formatter.Format(entry)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Unable to read entry, %v", err)
		return err
	}

	switch entry.Level {
	case logrus.PanicLevel:
		return sink.Writer.Crit(string(line))
	case logrus.FatalLevel:
		return sink.Writer.Crit(string(line))
	case logrus.ErrorLevel:
		return sink.Writer.Err(string(line))
	case logrus.WarnLevel:
		return sink.Writer.Warning(string(line))
	case logrus.InfoLevel:
		return sink.Writer.Info(string(line))
	case logrus.DebugLevel, logrus.TraceLevel:
		return sink.Writer.Debug(string(line))
	default:
		return nil
	}
}

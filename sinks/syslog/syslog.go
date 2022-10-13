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
	writer *syslog.Writer
	logrus.SinkBase
}

// NewSink Creates a hook to be added to an instance of logger. This is called with
// `sink, err := NewSink("udp", "localhost:514", syslog.LOG_DEBUG, "")`
func NewSink(f logrus.Formatter, lvl logrus.Level, network, raddr string, priority syslog.Priority, tag string) (*Sink, error) {
	w, err := syslog.Dial(network, raddr, priority, tag)
	return NewSinkWriter(f, w, lvl), err
}

// NewSinkWriter Creates a hook to be added to an instance of logger. This is called with
func NewSinkWriter(f logrus.Formatter, w *syslog.Writer, lvl logrus.Level) *Sink {
	return &Sink{writer: w, SinkBase: logrus.SinkBase{Formatter: f, Level: lvl}}
}

func (s *Sink) Emit(entry *logrus.Entry) error {
	line, err := s.Formatter.Format(entry)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Unable to read entry, %v", err)
		return err
	}

	switch entry.Level {
	case logrus.PanicLevel:
		return s.writer.Crit(string(line))
	case logrus.FatalLevel:
		return s.writer.Crit(string(line))
	case logrus.ErrorLevel:
		return s.writer.Err(string(line))
	case logrus.WarnLevel:
		return s.writer.Warning(string(line))
	case logrus.InfoLevel:
		return s.writer.Info(string(line))
	case logrus.DebugLevel, logrus.TraceLevel:
		return s.writer.Debug(string(line))
	default:
		return nil
	}
}

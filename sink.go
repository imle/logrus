package logrus

import (
	"fmt"
	"os"
	"sync"
	"sync/atomic"
)

var defaultSink = SinkWriter{
	mu:  sync.Mutex{},
	out: os.Stderr,
	SinkBase: SinkBase{
		Level:     InfoLevel,
		Formatter: &TextFormatter{},
	},
}

type Sink interface {
	Emit(*Entry) error
	EnabledAtLevel(Level) bool
	SetLevel(Level)
	GetLevel() Level
	SetFormatter(Formatter)
}

type SinkBase struct {
	Level     Level
	Formatter Formatter
}

func (s *SinkBase) EnabledAtLevel(lvl Level) bool {
	return s.Level >= lvl
}

func (s *SinkBase) SetLevel(lvl Level) {
	atomic.StoreUint32((*uint32)(&s.Level), uint32(lvl))
}

func (s *SinkBase) GetLevel() Level {
	return s.Level
}

func (s *SinkBase) SetFormatter(formatter Formatter) {
	s.Formatter = formatter
}

type LevelSinks []Sink

func (sinks LevelSinks) EnabledAtLevel(lvl Level) bool {
	for _, sink := range sinks {
		if sink.EnabledAtLevel(lvl) {
			return true
		}
	}

	return false
}

func (sinks LevelSinks) emit(entry *Entry) {
	if sinks.Empty() && defaultSink.EnabledAtLevel(entry.Level) {
		if err := emitWithBuffer(&defaultSink, entry); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "failed to write to log, %v\n", err)
		}
	}

	for _, sink := range sinks {
		if !sink.EnabledAtLevel(entry.Level) {
			continue
		}
		if err := emitWithBuffer(sink, entry); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "failed to write to log, %v\n", err)
		}
	}
}

func emitWithBuffer(sink Sink, entry *Entry) error {
	bufPool := entry.getBufferPool()
	buffer := bufPool.Get()
	defer func() {
		entry.Buffer = nil
		buffer.Reset()
		bufPool.Put(buffer)
	}()
	buffer.Reset()
	entry.Buffer = buffer

	return sink.Emit(entry)
}

func (sinks LevelSinks) Empty() bool {
	return len(sinks) == 0
}

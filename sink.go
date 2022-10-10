package logrus

import (
	"fmt"
	"os"
	"sync"
)

var defaultSink = SinkWriter{
	mu:        sync.Mutex{},
	Out:       os.Stderr,
	Formatter: &TextFormatter{},
}

const defaultLevel = InfoLevel

type LevelSinks map[Level][]Sink

type Sink interface {
	Emit(*Entry) error
}

// Add a hook to an instance of logger with the level specified and all above it.
func (sinks LevelSinks) Add(sink Sink, level Level) {
	for _, lvl := range AllLevels {
		if lvl > level {
			break
		}
		sinks[lvl] = append(sinks[lvl], sink)
	}
}

// AddLevels a hook to an instance of logger with the specific set of levels passed in.
func (sinks LevelSinks) AddLevels(sink Sink, levels []Level) {
	for _, level := range levels {
		sinks[level] = append(sinks[level], sink)
	}
}

func (sinks LevelSinks) emit(entry *Entry) {
	if sinks.Empty() && entry.Level <= InfoLevel {
		if err := emitWithBuffer(&defaultSink, entry); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "failed to write to log, %v\n", err)
		}
	}

	for _, sink := range sinks[entry.Level] {
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
	if len(sinks) == 0 {
		return true
	}

	for _, level := range AllLevels {
		if len(sinks[level]) != 0 {
			return false
		}
	}

	return true
}

func (sinks LevelSinks) LevelEmpty(level Level) bool {
	return len(sinks) == 0 || len(sinks[level]) == 0
}

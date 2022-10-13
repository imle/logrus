package logrus

import (
	"fmt"
	"io"
	"sync"
)

type SinkWriter struct {
	out io.Writer
	SinkBase

	mu               sync.Mutex
	terminalInitOnce sync.Once
	outIsTerminal    bool
}

func NewSinkWriter(out io.Writer, formatter Formatter, lvl Level) *SinkWriter {
	return &SinkWriter{out: out, SinkBase: SinkBase{Formatter: formatter, Level: lvl}}
}

func (s *SinkWriter) Emit(entry *Entry) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.terminalInitOnce.Do(func() {
		s.outIsTerminal = checkIfTerminal(s.out)
	})

	var err error
	var serialized []byte
	if s.Formatter == nil {
		serialized, err = (&TextFormatter{}).Format(entry)
	} else if tf, ok := s.Formatter.(TerminalFormatter); s.outIsTerminal && ok {
		serialized, err = tf.FormatForTerminal(entry)
	} else {
		serialized, err = s.Formatter.Format(entry)
	}

	if err != nil {
		return fmt.Errorf("failed to format: %w", err)
	}

	if _, err := s.out.Write(serialized); err != nil {
		return fmt.Errorf("failed to write to out: %w", err)
	}

	return nil
}

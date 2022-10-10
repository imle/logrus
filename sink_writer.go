package logrus

import (
	"fmt"
	"io"
	"sync"
)

type SinkWriter struct {
	Out       io.Writer
	Formatter Formatter

	mu               sync.Mutex
	terminalInitOnce sync.Once
	outIsTerminal    bool
}

func NewSinkWriter(out io.Writer) *SinkWriter {
	return &SinkWriter{Out: out, Formatter: &TextFormatter{}}
}

func (s *SinkWriter) Emit(entry *Entry) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.terminalInitOnce.Do(func() {
		s.outIsTerminal = checkIfTerminal(s.Out)
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

	if _, err := s.Out.Write(serialized); err != nil {
		return fmt.Errorf("failed to write to Out: %w", err)
	}

	return nil
}

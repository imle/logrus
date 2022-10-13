package logrus

import (
	"fmt"
	"io"
	"sync"
)

type SinkWriter struct {
	out io.Writer
	SinkBase

	mu            sync.Mutex
	outIsTerminal bool
}

func NewSinkWriter(out io.Writer, formatter Formatter, lvl Level) *SinkWriter {
	s := &SinkWriter{out: out, SinkBase: SinkBase{Formatter: formatter, Level: lvl}}
	s.outIsTerminal = checkIfTerminal(s.out)
	return s
}

func (s *SinkWriter) Emit(entry *Entry) error {
	s.mu.Lock()
	defer s.mu.Unlock()

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

func (s *SinkWriter) SetOutput(out io.Writer) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.out = out
	s.outIsTerminal = checkIfTerminal(s.out)
}

package logrus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFieldValueError(t *testing.T) {
	buf := &bytes.Buffer{}
	l := &Logger{
		Hooks: make(LevelHooks),
	}
	l.RegisterSink(&SinkWriter{Out: buf, Formatter: &JSONFormatter{}}, DebugLevel)
	l.WithField("func", func() {}).Info("test")
	fmt.Println(buf.String())
	var data map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &data); err != nil {
		t.Error("unexpected error", err)
	}
	_, ok := data[FieldKeyLogrusError]
	require.True(t, ok, `cannot found expected "logrus_error" field: %v`, data)
}

func TestNoFieldValueError(t *testing.T) {
	buf := &bytes.Buffer{}
	l := &Logger{
		Hooks: make(LevelHooks),
	}
	l.RegisterSink(&SinkWriter{Out: buf, Formatter: &JSONFormatter{}}, DebugLevel)
	l.WithField("str", "str").Info("test")
	fmt.Println(buf.String())
	var data map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &data); err != nil {
		t.Error("unexpected error", err)
	}
	_, ok := data[FieldKeyLogrusError]
	require.False(t, ok)
}

func TestWarninglnNotEqualToWarning(t *testing.T) {
	formatter := new(TextFormatter)
	formatter.DisableTimestamp = true
	formatter.DisableLevelTruncation = true

	buf := &bytes.Buffer{}
	log := &Logger{Hooks: make(LevelHooks)}
	log.RegisterSink(&SinkWriter{Out: buf, Formatter: formatter}, DebugLevel)
	log.Warning("hello,", "world")

	bufLn := &bytes.Buffer{}
	logLn := &Logger{Hooks: make(LevelHooks)}
	logLn.RegisterSink(&SinkWriter{Out: bufLn, Formatter: formatter}, DebugLevel)
	logLn.Warningln("hello,", "world")

	assert.NotEqual(t, buf.String(), bufLn.String(), "Warning() and Wantingln() should not be equal")
}

type testBufferPool struct {
	buffers []*bytes.Buffer
	get     int
}

func (p *testBufferPool) Get() *bytes.Buffer {
	p.get++
	return new(bytes.Buffer)
}

func (p *testBufferPool) Put(buf *bytes.Buffer) {
	p.buffers = append(p.buffers, buf)
}

func TestLogger_SetBufferPool(t *testing.T) {
	out := &bytes.Buffer{}
	l := New()
	l.RegisterSink(&SinkWriter{Out: out}, InfoLevel)

	pool := new(testBufferPool)
	l.SetBufferPool(pool)

	l.Info("test")

	assert.Equal(t, pool.get, 1, "Logger.SetBufferPool(): The BufferPool.Get() must be called")
	assert.Len(t, pool.buffers, 1, "Logger.SetBufferPool(): The BufferPool.Put() must be called")
}

package logrus_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	log "github.com/sirupsen/logrus"
)

func TestLogger_LogFn(t *testing.T) {
	log.RegisterSink(&log.SinkWriter{Out: os.Stderr, Formatter: &log.JSONFormatter{}}, log.WarnLevel)

	notCalled := 0
	log.InfoFn(func() []interface{} {
		notCalled++
		return []interface{}{
			"Hello",
		}
	})
	assert.Equal(t, 0, notCalled)

	called := 0
	log.ErrorFn(func() []interface{} {
		called++
		return []interface{}{
			"Oopsi",
		}
	})
	assert.Equal(t, 1, called)
}

package testutils

import (
	"bytes"
	"encoding/json"
	"strconv"
	"strings"
	"testing"

	. "github.com/sirupsen/logrus"

	"github.com/stretchr/testify/require"
)

func LogAndAssertJSON(t *testing.T, log func(*Logger, *SinkWriter), assertions func(fields Fields)) {
	var buffer bytes.Buffer
	var fields Fields

	logger := New()
	sink := &SinkWriter{Out: &buffer, Formatter: &JSONFormatter{}}
	logger.RegisterSink(sink, InfoLevel)

	log(logger, sink)

	err := json.Unmarshal(buffer.Bytes(), &fields)
	require.Nil(t, err)

	assertions(fields)
}

func LogAndAssertText(t *testing.T, log func(*Logger, *SinkWriter), assertions func(fields map[string]string)) {
	var buffer bytes.Buffer

	logger := New()
	sink := &SinkWriter{Out: &buffer, Formatter: &TextFormatter{DisableColors: true}}
	logger.RegisterSink(sink, InfoLevel)

	log(logger, sink)

	fields := make(map[string]string)
	for _, kv := range strings.Split(strings.TrimRight(buffer.String(), "\n"), " ") {
		if !strings.Contains(kv, "=") {
			continue
		}
		kvArr := strings.Split(kv, "=")
		key := strings.TrimSpace(kvArr[0])
		val := kvArr[1]
		if kvArr[1][0] == '"' {
			var err error
			val, err = strconv.Unquote(val)
			require.NoError(t, err)
		}
		fields[key] = val
	}
	assertions(fields)
}

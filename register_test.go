package logrus

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLocalhostAddAndPrint(t *testing.T) {
	log := New()
	formatter := &TextFormatter{DisableTimestamp: true, DisableColors: true}

	require.Len(t, log.sinks[PanicLevel], 0, "sink was not added properly")
	require.Len(t, log.sinks[FatalLevel], 0, "sink was not added properly")
	require.Len(t, log.sinks[ErrorLevel], 0, "sink was not added properly")
	require.Len(t, log.sinks[WarnLevel], 0, "sink was not added properly")
	require.Len(t, log.sinks[InfoLevel], 0, "sink was not added properly")
	require.Len(t, log.sinks[DebugLevel], 0, "sink was not added properly")
	require.Len(t, log.sinks[TraceLevel], 0, "sink was not added properly")

	log.RegisterSink(&SinkWriter{Out: os.Stdout, Formatter: formatter}, TraceLevel)
	require.Len(t, log.sinks[PanicLevel], 1, "sink was not added properly")
	require.Len(t, log.sinks[FatalLevel], 1, "sink was not added properly")
	require.Len(t, log.sinks[ErrorLevel], 1, "sink was not added properly")
	require.Len(t, log.sinks[WarnLevel], 1, "sink was not added properly")
	require.Len(t, log.sinks[InfoLevel], 1, "sink was not added properly")
	require.Len(t, log.sinks[DebugLevel], 1, "sink was not added properly")
	require.Len(t, log.sinks[TraceLevel], 1, "sink was not added properly")

	log.RegisterSink(&SinkWriter{Out: os.Stdout, Formatter: formatter}, InfoLevel)
	require.Len(t, log.sinks[PanicLevel], 2, "sink was not added properly")
	require.Len(t, log.sinks[FatalLevel], 2, "sink was not added properly")
	require.Len(t, log.sinks[ErrorLevel], 2, "sink was not added properly")
	require.Len(t, log.sinks[WarnLevel], 2, "sink was not added properly")
	require.Len(t, log.sinks[InfoLevel], 2, "sink was not added properly")
	require.Len(t, log.sinks[DebugLevel], 1, "sink was not added properly")
	require.Len(t, log.sinks[TraceLevel], 1, "sink was not added properly")

	log.Info("Congratulations!")
}

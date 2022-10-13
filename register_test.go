package logrus

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLocalhostAddAndPrint(t *testing.T) {
	log := New()
	formatter := &TextFormatter{DisableTimestamp: true, DisableColors: true}

	require.Len(t, log.sinks, 0, "sink was not added properly")
	require.False(t, log.sinks.EnabledAtLevel(PanicLevel), "sink was not added properly")
	require.False(t, log.sinks.EnabledAtLevel(FatalLevel), "sink was not added properly")
	require.False(t, log.sinks.EnabledAtLevel(ErrorLevel), "sink was not added properly")
	require.False(t, log.sinks.EnabledAtLevel(WarnLevel), "sink was not added properly")
	require.False(t, log.sinks.EnabledAtLevel(InfoLevel), "sink was not added properly")
	require.False(t, log.sinks.EnabledAtLevel(DebugLevel), "sink was not added properly")
	require.False(t, log.sinks.EnabledAtLevel(TraceLevel), "sink was not added properly")

	log.RegisterSink(NewSinkWriter(os.Stdout, formatter, InfoLevel))
	require.Len(t, log.sinks, 1, "sink was not added properly")
	require.True(t, log.sinks.EnabledAtLevel(PanicLevel), "sink was not added properly")
	require.True(t, log.sinks.EnabledAtLevel(FatalLevel), "sink was not added properly")
	require.True(t, log.sinks.EnabledAtLevel(ErrorLevel), "sink was not added properly")
	require.True(t, log.sinks.EnabledAtLevel(WarnLevel), "sink was not added properly")
	require.True(t, log.sinks.EnabledAtLevel(InfoLevel), "sink was not added properly")
	require.False(t, log.sinks.EnabledAtLevel(DebugLevel), "sink was not added properly")
	require.False(t, log.sinks.EnabledAtLevel(TraceLevel), "sink was not added properly")

	log.RegisterSink(NewSinkWriter(os.Stdout, formatter, TraceLevel))
	require.Len(t, log.sinks, 2, "sink was not added properly")
	require.True(t, log.sinks.EnabledAtLevel(PanicLevel), "sink was not added properly")
	require.True(t, log.sinks.EnabledAtLevel(FatalLevel), "sink was not added properly")
	require.True(t, log.sinks.EnabledAtLevel(ErrorLevel), "sink was not added properly")
	require.True(t, log.sinks.EnabledAtLevel(WarnLevel), "sink was not added properly")
	require.True(t, log.sinks.EnabledAtLevel(InfoLevel), "sink was not added properly")
	require.True(t, log.sinks.EnabledAtLevel(DebugLevel), "sink was not added properly")
	require.True(t, log.sinks.EnabledAtLevel(TraceLevel), "sink was not added properly")

	log.Info("Congratulations!")
}

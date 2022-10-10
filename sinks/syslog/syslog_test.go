//go:build !windows && !nacl && !plan9
// +build !windows,!nacl,!plan9

package syslog

import (
	"log/syslog"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/sirupsen/logrus"

	syslogserver "gopkg.in/mcuadros/go-syslog.v2"
)

func TestLocalhostAddAndPrint(t *testing.T) {
	channel := make(syslogserver.LogPartsChannel)
	handler := syslogserver.NewChannelHandler(channel)

	server := syslogserver.NewServer()
	server.SetFormat(syslogserver.RFC5424)
	server.SetHandler(handler)
	err := server.ListenUDP("localhost:32142")
	require.NoError(t, err)
	err = server.Boot()
	require.NoError(t, err)
	defer func() {
		err = server.Kill()
		require.NoError(t, err)
	}()

	log := logrus.New()
	formatter := &logrus.TextFormatter{DisableTimestamp: true, DisableColors: true}
	sink, err := NewSink(formatter, "udp", "localhost:32142", syslog.LOG_INFO, "")
	if err != nil {
		t.Errorf("Unable to connect to local syslog.")
	}
	log.RegisterSink(sink, logrus.TraceLevel)

	for _, level := range logrus.AllLevels {
		if !log.IsLevelEnabled(level) {
			t.Errorf("Sink was not added. The length of log.Sinks[%v]: %v", level, len(log.Hooks[level]))
		}
	}

	log.Info("Congratulations!")

	tick := time.NewTimer(100 * time.Millisecond)
	select {
	case logValue := <-channel:
		require.Equal(t, logValue["message"], "Congratulations!")
	case <-tick.C:
		return
	}
}

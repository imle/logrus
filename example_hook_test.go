//go:build !windows
// +build !windows

package logrus_test

import (
	"log/syslog"
	"os"

	"github.com/sirupsen/logrus"
	slhooks "github.com/sirupsen/logrus/sinks/syslog"
)

// An example on how to use a hook
func Example_hook() {
	formatter := &logrus.TextFormatter{
		DisableColors:    true, // remove colors
		DisableTimestamp: true, // remove timestamp from test output
	}

	var log = logrus.New()

	sink, err := slhooks.NewSink(formatter, "udp", "localhost:514", syslog.LOG_INFO, "")
	if err != nil {
		panic(err)
	}
	log.RegisterSink(sink, logrus.InfoLevel)
	log.RegisterSink(&logrus.SinkWriter{Out: os.Stdout, Formatter: formatter}, logrus.InfoLevel)

	log.WithFields(logrus.Fields{
		"animal": "walrus",
		"size":   10,
	}).Info("A group of walrus emerges from the ocean")

	log.WithFields(logrus.Fields{
		"omg":    true,
		"number": 122,
	}).Warn("The group's number increased tremendously!")

	log.WithFields(logrus.Fields{
		"omg":    true,
		"number": 100,
	}).Error("The ice breaks!")

	// Output:
	// level=info msg="A group of walrus emerges from the ocean" animal=walrus size=10
	// level=warning msg="The group's number increased tremendously!" number=122 omg=true
	// level=error msg="The ice breaks!" number=100 omg=true
}

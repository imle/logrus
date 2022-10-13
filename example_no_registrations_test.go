//go:build !windows
// +build !windows

package logrus

import (
	"os"
)

// An example on how to use a hook
func Example_noRegistrations() {
	var log = New()
	defaultSink.out = os.Stdout
	defaultSink.Formatter.(*TextFormatter).DisableTimestamp = true

	log.WithFields(Fields{
		"animal": "walrus",
		"size":   10,
	}).Info("A group of walrus emerges from the ocean")

	log.WithFields(Fields{
		"omg":    true,
		"number": 122,
	}).Warn("The group's number increased tremendously!")

	log.WithFields(Fields{
		"omg":    true,
		"number": 100,
	}).Error("The ice breaks!")

	// Output:
	// level=info msg="A group of walrus emerges from the ocean" animal=walrus size=10
	// level=warning msg="The group's number increased tremendously!" number=122 omg=true
	// level=error msg="The ice breaks!" number=100 omg=true
}

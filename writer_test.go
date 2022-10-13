package logrus_test

import (
	"log"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
)

func ExampleLogger_Writer_httpServer() {
	logger := logrus.New()
	w := logger.Writer()
	defer w.Close()

	srv := http.Server{
		// create a stdlib log.Logger that writes to
		// logrus.Logger.
		ErrorLog: log.New(w, "", 0),
	}

	if err := srv.ListenAndServe(); err != nil {
		logger.Fatal(err)
	}
}

func ExampleLogger_Writer_stdlib() {
	logger := logrus.New()
	logger.RegisterSink(logrus.NewSinkWriter(os.Stderr, &logrus.JSONFormatter{}, logrus.InfoLevel))

	// Use logrus for standard log output
	// Note that `log` here references stdlib's log
	// Not logrus imported under the name `log`.
	log.SetOutput(logger.Writer())
}

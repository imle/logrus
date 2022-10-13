package logrus_test

import (
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

func ExampleJSONFormatter_CallerPrettyfier() {
	l := logrus.New()
	l.SetReportCaller(true)
	formatter := &logrus.JSONFormatter{
		DisableTimestamp: true,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			s := strings.Split(f.Function, ".")
			funcname := s[len(s)-1]
			_, filename := path.Split(f.File)
			return funcname, filename
		},
	}
	l.RegisterSink(logrus.NewSinkWriter(os.Stdout, formatter, logrus.InfoLevel))

	l.Info("example of custom format caller")
	// Output:
	// {"file":"example_custom_caller_test.go","func":"ExampleJSONFormatter_CallerPrettyfier","level":"info","msg":"example of custom format caller"}
}

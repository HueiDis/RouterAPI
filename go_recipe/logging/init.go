package logging

import (
	"flag"
	"log"

	"github.com/onrik/logrus/filename" 
	"github.com/sirupsen/logrus"
)

var (
	logJSON = flag.Bool("log_json", false, "log in Json instead of text.")
)

func newFormatter() logrus.Formatter {
	if *logJSON {
		return new(logrus.JSONFormatter)
	}
	customFormatter := new(logrus.TextFormatter)
	customFormatter.FullTimestamp = true
	customFormatter.DisableSorting = true
	return customFormatter
}

// Init initalizes logger formatter for logrus logging.
func Init() {
	filenameHook := filename.NewHook(logrus.PanicLevel, logrus.FatalLevel, logrus.ErrorLevel, logrus.WarnLevel, logrus.DebugLevel)
	formatter := newFormatter()
	logrus.AddHook(filenameHook)
	logrus.SetFormatter(formatter)

	// Also let all external libraries use the same log formatter.
	logger := logrus.New()
	logger.AddHook(filenameHook)
	logger.Formatter = formatter
	log.SetFlags(0) // time is already handled by logrus.
	log.SetOutput(logger.Writer())
}

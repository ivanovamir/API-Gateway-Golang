package logging

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"runtime"
)

var e *logrus.Entry

type Logger struct {
	*logrus.Entry
}

func NewLogger() *Logger {
	return &Logger{e}
}

func init() {
	log := logrus.New()
	log.SetReportCaller(true)

	log.Formatter = &logrus.TextFormatter{
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			filename := path.Base(frame.File)
			return fmt.Sprintf("%s()", frame.Function), fmt.Sprintf("%s:%d", filename, frame.Line)
		},
		DisableColors: false,
		FullTimestamp: true,
	}

	err := os.MkdirAll("logs", 0755)

	if err != nil {
		panic(err)
	}

	allFile, err := os.OpenFile("logs/base.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)

	if err != nil {
		panic(err)
	}

	log.SetOutput(allFile)

	log.SetLevel(logrus.TraceLevel)

	e = logrus.NewEntry(log)
}

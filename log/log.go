package log

import (
	"log"
	"os"

	"github.com/sirupsen/logrus"
)

var (
	I logrus.Logger // Info Logger
	E logrus.Logger // Error Logger
	H logrus.Logger // HTTP Logger
)

const (
	TimeFMT string = "2006-01-02 15:04:05.0000"
)

func DefLoggers(ILog, ELog, HLog *logrus.Logger) {
	LoggerInit(ILog, logrus.InfoLevel, &logrus.TextFormatter{TimestampFormat: TimeFMT}, "log/info.log")
	LoggerInit(ELog, logrus.ErrorLevel, &logrus.TextFormatter{TimestampFormat: TimeFMT}, "log/error.log")
	LoggerInit(HLog, logrus.InfoLevel, &logrus.JSONFormatter{}, "log/http.log")
}

func LoggerInit(logger *logrus.Logger, level logrus.Level, formatter logrus.Formatter, path string) {

	logFile, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0222)
	if err != nil {
		log.Panic(err)
	}

	x := logrus.New()
	x.SetOutput(logFile)
	x.SetLevel(level)
	x.SetFormatter(formatter)

	*logger = *x
}

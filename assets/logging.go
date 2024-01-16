package assets

import (
	"io"
	"io/fs"
	"os"

	"github.com/alexandrepasc/ev-chargers-simulator-2/common"
	"github.com/alexandrepasc/ev-chargers-simulator-2/translation"
	"github.com/alexandrepasc/ev-chargers-simulator-2/translation/text"
	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

type Logging struct {
	Logger *logrus.Logger
	ToFile bool
	File   string
	Level  logrus.Level
	L      translation.Translation
}

/*
Perform the log action, setting the output and the log level defined in the structure. Will validate
that the message level sent matches the current level of the logging and if so print the log.

m	-	Map of the fields that will be used in to log the message (map[string]string)

msg	-	The message that will be logged (interface{})

sev	-	The severity level of the message that is being logged (assets.Severity)
*/
func (l Logging) Log(m map[string]string, msg interface{}, sev Severity) {
	l.setOutput()
	l.setLevel()

	switch sev {
	case Info:
		if l.Logger.GetLevel() == logrus.InfoLevel || l.Logger.GetLevel() == logrus.DebugLevel {
			var le = l.setFields(m)

			le.Info(msg)
		}

	case Error:
		if l.Logger.GetLevel() == logrus.ErrorLevel || l.Logger.GetLevel() == logrus.DebugLevel {
			var le = l.setFields(m)

			le.Error(msg)
		}

	case Warn:
		if l.Logger.GetLevel() == logrus.WarnLevel || l.Logger.GetLevel() == logrus.DebugLevel {
			var le = l.setFields(m)

			le.Warn(msg)
		}

	case Panic:
		if l.Logger.GetLevel() == logrus.PanicLevel || l.Logger.GetLevel() == logrus.DebugLevel {
			var le = l.setFields(m)

			le.Panic(msg)
		}

	case Fatal:
		if logger.GetLevel() == logrus.FatalLevel || l.Logger.GetLevel() == logrus.DebugLevel {
			var le = l.setFields(m)

			le.Fatal(msg)
		}

	case Debug:
		if l.Logger.GetLevel() == logrus.DebugLevel {
			var le = l.setFields(m)

			le.Debug(msg)
		}

	default:
		var le = l.setFields(map[string]string{"function": "log"})

		le.Error(l.L.Get(text.LogMissingSeverityLevel))
	}
}

/*
Sets the log fields to the logger and returns the log entry (*logrus.Entry)

m	-	The map of the fields and values that will be added to the log.
*/
func (l Logging) setFields(m map[string]string) (le *logrus.Entry) {
	var lf = logrus.Fields{}

	for k, v := range m {
		lf[k] = v
	}

	le = l.Logger.WithFields(lf)

	return le
}

func (l Logging) setOutput() {
	if l.ToFile {
		f, err := os.OpenFile(l.File, os.O_APPEND|os.O_CREATE|os.O_RDWR, fs.FileMode(common.FilePermissions))

		if err != nil {
			common.Log("setOutput").Fatal("can not write to log file")
		}

		mw := io.MultiWriter(os.Stdout, f)

		l.Logger.SetOutput(mw)
	} else {
		l.Logger.SetOutput(os.Stdout)
	}
}

func (l Logging) setLevel() {
	l.Logger.SetLevel(l.Level)
}

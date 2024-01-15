package ocpp16

import (
	"io"
	"io/fs"
	"os"

	"github.com/alexandrepasc/ev-chargers-simulator-2/assets"
	"github.com/alexandrepasc/ev-chargers-simulator-2/common"
	"github.com/alexandrepasc/ev-chargers-simulator-2/translation"
	"github.com/alexandrepasc/ev-chargers-simulator-2/translation/text"
	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

type logging struct {
	logger *logrus.Logger
	toFile bool
	file   string
	level  logrus.Level
	l      translation.Translation
}

/*
Perform the log action, setting the output and the log level defined in the structure. Will validate
that the message level sent matches the current level of the logging and if so print the log.

m	-	Map of the fields that will be used in to log the message (map[string]string)

msg	-	The message that will be logged (interface{})

sev	-	The severity level of the message that is being logged (assets.Severity)
*/
func (l logging) log(m map[string]string, msg interface{}, sev assets.Severity) {
	l.setOutput()
	l.setLevel()

	switch sev {
	case assets.Info:
		if l.logger.GetLevel() == logrus.InfoLevel || l.logger.GetLevel() == logrus.DebugLevel {
			var le = l.setFields(m)

			le.Info(msg)
		}

	case assets.Error:
		if l.logger.GetLevel() == logrus.ErrorLevel || l.logger.GetLevel() == logrus.DebugLevel {
			var le = l.setFields(m)

			le.Error(msg)
		}

	case assets.Warn:
		if l.logger.GetLevel() == logrus.WarnLevel || l.logger.GetLevel() == logrus.DebugLevel {
			var le = l.setFields(m)

			le.Warn(msg)
		}

	case assets.Panic:
		if l.logger.GetLevel() == logrus.PanicLevel || l.logger.GetLevel() == logrus.DebugLevel {
			var le = l.setFields(m)

			le.Panic(msg)
		}

	case assets.Fatal:
		if logger.GetLevel() == logrus.FatalLevel || l.logger.GetLevel() == logrus.DebugLevel {
			var le = l.setFields(m)

			le.Fatal(msg)
		}

	case assets.Debug:
		if l.logger.GetLevel() == logrus.DebugLevel {
			var le = l.setFields(m)

			le.Debug(msg)
		}

	default:
		var le = l.setFields(map[string]string{"function": "log"})

		le.Error(text.LogMissingSeveretyLevel)
	}
}

/*
Sets the log fields to the logger and returns the log entry (*logrus.Entry)

m	-	The map of the fields and values that will be added to the log.
*/
func (l logging) setFields(m map[string]string) (le *logrus.Entry) {
	var lf = logrus.Fields{}

	for k, v := range m {
		lf[k] = v
	}

	le = l.logger.WithFields(lf)

	return le
}

func (l logging) setOutput() {
	if l.toFile {
		f, err := os.OpenFile(l.file, os.O_APPEND|os.O_CREATE|os.O_RDWR, fs.FileMode(common.FilePermissions))

		if err != nil {
			common.Log("setOutput").Fatal("can not write to log file")
		}

		mw := io.MultiWriter(os.Stdout, f)

		l.logger.SetOutput(mw)
	} else {
		l.logger.SetOutput(os.Stdout)
	}
}

func (l logging) setLevel() {
	l.logger.SetLevel(l.level)
}

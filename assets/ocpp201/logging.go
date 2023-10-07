package ocpp201

import (
	"io"
	"io/fs"
	"os"

	"github.com/alexandrepasc/ev-chargers-simulator-2/assets"
	"github.com/alexandrepasc/ev-chargers-simulator-2/common"
	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

type logging struct {
	toFile bool
	file   string
}

func init() {
	logger = logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	logger.SetLevel(logrus.DebugLevel)
}

/**/
func (l logging) log(m map[string]string, msg interface{}, sev assets.Severity) {
	l.setOutput()

	var le = setFields(m)

	switch sev {
	case assets.Info:
		le.Info(msg)
	case assets.Error:
		le.Error(msg)
	case assets.Warn:
		le.Warn(msg)
	case assets.Panic:
		le.Panic(msg)
	case assets.Fatal:
		le.Fatal(msg)
	default:
		le.Info(msg)
	}
}

/*
Sets the log fields to the logger and returns the log entry (*logrus.Entry)

m	-	The map of the fields and values that will be added to the log.
*/
func setFields(m map[string]string) (le *logrus.Entry) {
	var lf = logrus.Fields{}

	for k, v := range m {
		lf[k] = v
	}

	le = logger.WithFields(lf)

	return le
}

func (l logging) setOutput() {
	if l.toFile {
		f, err := os.OpenFile(l.file, os.O_APPEND|os.O_CREATE|os.O_RDWR, fs.FileMode(common.FilePermissions))

		if err != nil {
			common.Log("setOutput").Fatal("can not write to log file")
		}

		mw := io.MultiWriter(os.Stdout, f)

		logger.SetOutput(mw)
	} else {
		logger.SetOutput(os.Stdout)
	}
}

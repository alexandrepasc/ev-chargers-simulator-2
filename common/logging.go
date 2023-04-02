package common

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

/*
Start a new logrus entry, setting the output and the format.
*/
func init() {
	log = logrus.New()
	log.SetOutput(os.Stdout)
	log.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	log.SetLevel(logrus.DebugLevel)
}

/*
Return a log entry with "Function" as log field.

fn	-	Function entry (string)

le	-	Log entry already created (*logrus.Entry)
*/
func Log(fn string, le *logrus.Entry) *logrus.Entry {
	le = le.WithField("Function", fn)

	return le
}

/*
Return a log entry with "Protocol" and "Model" as log field.

pn	-	Protocol entry (string)

mn	-	Model entry (string)
*/
func PreLog(pn, mn string) *logrus.Entry {
	var le = log.WithField("Protocol", pn).WithField("Model", mn)

	return le
}

func APILog(fn string) *logrus.Entry {
	var le = log.WithField("Service", "API").WithField("Function", fn)

	return le
}

/*
Print the starting log printed to the cmd.
*/
func StartLog() {
	fmt.Println("")
	fmt.Println("")
	fmt.Println("               ########   #####             ######   #########")
	fmt.Println("               #########   #####             ###    ##########              \"")
	fmt.Println("               ###          #####           ###    ###             mmm    mmm     mmmmm")
	fmt.Println("               ###           #####         ###    ###             #   \"     #     # # #")
	fmt.Println("               ##########     #####       ###    ###               \"\"\"m     #     # # #")
	fmt.Println("               ###########     #####     ###    ###               \"mmm\"   mm#mm   # # #")
	fmt.Println("               ###              #####   ###    ####")
	fmt.Println("               ###               ##### ###    #####")
	fmt.Println("               ############       #######    ###############         version: ", Version)
	fmt.Println("               #############       #####    ################")
	fmt.Println("")
	fmt.Println("")
}

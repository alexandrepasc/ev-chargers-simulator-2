package common

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func init() {
	log = logrus.New()
	log.SetOutput(os.Stdout)
	log.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	log.SetLevel(logrus.DebugLevel)
}

// fn = function name
// le = logrus entry
func Log(fn string, le *logrus.Entry) *logrus.Entry {
	le = le.WithField("Function", fn)

	return le
}

// pn = protocol name
// mn = model name
func PreLog(pn, mn string) *logrus.Entry {
	var le = log.WithField("Protocol", pn).WithField("Model", mn)

	return le
}

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

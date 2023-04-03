package common

import (
	"os"
	"path/filepath"
)

/*
It gets the executable path and add the folder, sent as argument, to the end of the path.

Returns a string.

f	-	folder to add to the executable path (string)
*/
func getThePath(f string) string {
	if f[0:1] != "/" {
		f = "/" + f
	}

	var fp = getExecPath() + f

	return fp
}

/*
Get the full path for the executable, the path doesn't finish with the spliter ('/').

Returns a string.
*/
func getExecPath() string {
	ep, err := os.Executable()

	if err != nil {
		Log("getExecPath").Panic(err)
	}

	var p = filepath.Dir(ep)

	return p
}

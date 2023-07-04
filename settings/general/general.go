package general

import (
	"encoding/json"
	"io/fs"
	"net/http"
	"os"

	"github.com/alexandrepasc/ev-chargers-simulator-2/common"
	"github.com/alexandrepasc/ev-chargers-simulator-2/flags"
	"github.com/alexandrepasc/ev-chargers-simulator-2/translation"
	"github.com/alexandrepasc/ev-chargers-simulator-2/translation/text"
	"github.com/go-playground/validator/v10"
)

/*
Package dedicated to manage reading and writing the general settings of the application. This file
has the responsibility to store the base general settings that the application uses to work, for
example the host ip address, the connection timeout, etc....
*/
func General(p string, fl *flags.Flags) Model {
	var gen Model

	ok, fp := isGeneral(p)

	if ok {
		gen = openGeneralFile(fp)

		if fl.ForceUpdate {
			gen = *updateGeneralFile(fl, &gen, fp)
		}
	} else {
		createGeneralFile(fp)

		gen = updateNewGeneralFile(fl, fp)
	}

	return gen
}

/*
Opens the general settings file, returns it (*Model), and the http status code (int). It handles
the errors and returns the message translated as the ok variable with false (bool).

p	-	The path where the general settings file is located (string)

t	-	The translation structure with the language defined (translation.Translation)
*/
// TODO: there are some duplicated code to be able to handle the api erros, this needs to be reviewed
func GetGeneralConf(p string, t translation.Translation) (ok bool, msg string, code int, g *Model) {
	var fp = p + generalFile

	b, err := os.ReadFile(fp)

	if err != nil {
		common.Log("GetGeneralConf").Error(err)
		return false, t.Get(text.GeneralErrorRead), http.StatusInternalServerError, nil
	}

	var m Model

	mErr := json.Unmarshal(b, &m)

	if mErr != nil {
		common.Log("GetGeneralConf").Error(mErr)
		return false, t.Get(text.GeneralErrorRead), http.StatusInternalServerError, nil
	}

	return true, "", http.StatusOK, &m
}

/*
Update the general settings file with the model sent. After the update will read the file and
returns the data (*Model). It handles the errors and returns the ok varial as false (bool), the
error translated message (string). In both cases returns the http status code (int).

ng	-	The model with the data to update the file (*Model)

p	-	The path where the general settings file is located (string)

t	-	The translation structure with the language defined (translation.Translation)
*/
// TODO: there are some duplicated code to be able to handle the api erros, this needs to be reviewed
func UpdateGeneralConf(ng *Model, p string, t translation.Translation) (ok bool, msg string, code int, g *Model) {
	var fp = p + generalFile

	var vErr = validator.New().Struct(ng)

	if vErr != nil {
		common.Log("UpdateGeneralConf").Error(vErr)
		return false, vErr.Error(), http.StatusBadRequest, nil
	}

	var b, err = json.MarshalIndent(ng, "", " ")

	if err != nil {
		common.Log("UpdateGeneralConf").Error(err)
		return false, t.Get(text.GeneralErrorUpdate), http.StatusInternalServerError, nil
	}

	var wErr = os.WriteFile(fp, b, fs.FileMode(common.FilePermissions))

	if wErr != nil {
		common.Log("UpdateGeneralConf").Error(wErr)
		return false, t.Get(text.GeneralErrorUpdate), http.StatusInternalServerError, nil
	}

	var gOk, gMsg, gCode, gg = GetGeneralConf(p, t)

	if !gOk {
		common.Log("UpdateGeneralConf").Error(gMsg)
		return false, gMsg, gCode, nil
	}

	return true, "", http.StatusOK, gg
}

/*
Checks if the general settings file exists in the defined path.

If exists returns true, if not returns false.

In both cases returns the file path.

p	-	Defined general path (string)
*/
func isGeneral(p string) (is bool, fp string) {
	fp = p + generalFile

	s, err := os.Stat(fp)

	if err != nil {
		common.Log("isGeneral").Warn(err)

		return false, fp
	}

	if s != nil {
		return true, fp
	}

	common.Log("isGeneral").Warn(translation.Translation{L: translation.EnGb}.Get(text.GeneralNotExist))

	return false, fp
}

/*
Open the general configuration file, unmarshal it to a structure and return the Model.

f	-	Configuration file path (string)
*/
func openGeneralFile(fp string) (g Model) {
	b, err := os.ReadFile(fp)

	if err != nil {
		common.Log("openGeneralFile").Fatal(err)
	}

	g = unmarshalGeneralJSON(b)

	return g
}

/*
Create the general configuration file.

fp	-	General configuration file with path (string)
*/
func createGeneralFile(fp string) {
	f, err := os.Create(fp)

	if err != nil {
		common.Log("createGeneralFile").Fatal(err)
	}

	f.Close()

	common.Log("createGeneralFile").Info(translation.Translation{L: translation.EnGb}.Get(text.GeneralCreated))
}

/*
Update the new general configuration file with data.

fl	-	Flags set in the execution of the application (*flags.Flags)

f	-	Full path and file name for the configuration
*/
func updateNewGeneralFile(fl *flags.Flags, fp string) Model {
	var configs = Model{
		HostIP:      common.DefSimIP,
		ConnTimeout: common.DefTimeout,
		CSAddr:      common.DefCSIP,
		CSPort:      common.DefCSPort,
		Lang:        common.DefLanguage,
		APIAddr:     common.DefAPIAddr,
		APIPort:     common.DefAPIPort,
	}

	if fl.HostAddr != "" {
		configs.HostIP = fl.HostAddr
	}

	if fl.ConnTimeout != -1 {
		configs.ConnTimeout = fl.ConnTimeout
	}

	if fl.CSAddr != "" {
		configs.CSAddr = fl.CSAddr
	}

	if fl.CSPort != "" {
		configs.CSPort = fl.CSPort
	}

	if fl.Language != "" {
		configs.Lang = fl.Language
	}

	if fl.APIAddr != "" {
		configs.APIAddr = fl.APIAddr
	}

	if fl.APIPort != "" {
		configs.APIPort = fl.APIPort
	}

	var b = marshalIndentGeneral(&configs)

	var err = os.WriteFile(fp, b, fs.FileMode(common.FilePermissions))

	if err != nil {
		common.Log("updateNewGeneralFile").Fatal(err)
	}

	common.Log("updateNewGeneralFile").Info(translation.Translation{L: translation.EnGb}.Get(text.GeneralWritten))

	return configs
}

/*
Update the general configurations file with new data.

Returns the model with the new data, in case of failure return the old data (Model).

fl	-	Flags set in the execution of the application (*flags.Flags)

c	-	Model with the data read from the file (Model)

fp	-	Full path and file name for the general configurations (string)
*/
func updateGeneralFile(fl *flags.Flags, c *Model, fp string) *Model {
	var nc = Model{}

	if fl.HostAddr != "" {
		nc.HostIP = fl.HostAddr
	} else {
		nc.HostIP = c.HostIP
	}

	if fl.ConnTimeout != -1 {
		nc.ConnTimeout = fl.ConnTimeout
	} else {
		nc.ConnTimeout = c.ConnTimeout
	}

	if fl.CSAddr != "" {
		nc.CSAddr = fl.CSAddr
	} else {
		nc.CSAddr = c.CSAddr
	}

	if fl.CSPort != "" {
		nc.CSPort = fl.CSPort
	} else {
		nc.CSPort = c.CSPort
	}

	if fl.Language != "" {
		nc.Lang = fl.Language
	} else {
		nc.Lang = c.Lang
	}

	if fl.APIAddr != "" {
		nc.APIAddr = fl.APIAddr
	} else {
		nc.APIAddr = c.APIAddr
	}

	if fl.APIPort != "" {
		nc.APIPort = fl.APIPort
	} else {
		nc.APIPort = c.APIPort
	}

	var b = marshalIndentGeneral(&nc)

	var err = os.WriteFile(fp, b, fs.FileMode(common.FilePermissions))

	if err != nil {
		common.Log("updateGeneralFile").Error(err)

		return c
	}

	common.Log("updateGeneralFile").Info(translation.Translation{L: translation.EnGb}.Get(text.GeneralWritten))

	return &nc
}

/*
Convert general json file (in bytes) to the Model and returns it.

b	-	Json file converted into bytes ([]byte)
*/
func unmarshalGeneralJSON(b []byte) (g Model) {
	err := json.Unmarshal(b, &g)

	if err != nil {
		common.Log("unmarshalGeneralJSON").Fatal(err)
	}

	return g
}

/*
Deserialize the Model to json in an byte array.

Returns []byte.

c	-	The model with the file data (Model)
*/
func marshalIndentGeneral(c *Model) []byte {
	b, err := json.MarshalIndent(c, "", " ")

	if err != nil {
		common.Log("marshalIndent").Fatal(err)
	}

	return b
}

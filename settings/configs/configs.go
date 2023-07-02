package configs

import (
	"encoding/json"
	"fmt"
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
Package dedicated to manage reading and writing the configs file. This file has the responsibility
to store the base configurations that the application uses to work, for example the folder paths
to the general and simulators settings.
*/
func Configs(fl *flags.Flags) Model {
	var conf Model

	ok, fp := isConfigs()

	if ok {
		conf = openConfigsFile(fp)

		if fl.ForceUpdate {
			conf = updateConfigsFile(fl, conf, fp)
		}

		fmt.Println(conf)
	} else {
		createConfigsFile(fp)

		conf = updateNewConfigsFile(fl, fp)
	}

	return conf
}

/*
Read the configs file, get the simulators configuration path and set that to the SimPathModel. The
model will be returned to the api. This will handle the errors and return them to the api module.
It returns false (boolean) in case something fails, the message (string) with the error message,
the http status code (int), and the response model.

t	-	The translation structure with the language defined (translation.Translation)
*/
func GetConfigsSimPath(t translation.Translation) (ok bool, msg string, code int, s *SimPathModel) {
	var p = common.GetThePath(configFile)

	var b, rErr = os.ReadFile(p)

	if rErr != nil {
		common.Log("GetConfigsSimPath").Error(rErr)
		return false, t.Get(text.ConfigsErrorRead), http.StatusInternalServerError, nil
	}

	var m Model

	var mErr = json.Unmarshal(b, &m)

	if mErr != nil {
		common.Log("GetConfigsSimPath").Error(mErr)
		return false, t.Get(text.ConfigsErrorRead), http.StatusInternalServerError, nil
	}

	var rs = SimPathModel{
		SimulatorsConfigFolder: m.SimulatorsConfigFolder,
	}

	return true, "", http.StatusOK, &rs
}

/*
Update the configuration file with the new simulators path, reads the existing configs file
update the model with the new value, write the file with the updated information, and returns
the status code and the model in case no error was detected. In case ocurres an error it will
return false (boolean), the error message (string), the status code for the error (int), and
the model as nil (*SimPathModel).

ns	-	The information sent by the user using the api (*SimPathModel)

t	-	The translation structure with the language defined (translation.Translation)
*/
func UpdateConfigsSimPath(ns *SimPathModel, t translation.Translation) (ok bool, msg string, code int, s *SimPathModel) {
	var p = common.GetThePath(configFile)

	var vErr = validator.New().Struct(ns)

	if vErr != nil {
		common.Log("UpdateConfigsSimPath").Error(vErr)
		return false, vErr.Error(), http.StatusBadRequest, nil
	}

	// Read the configuration file
	var rb, rErr = os.ReadFile(p)

	if rErr != nil {
		common.Log("UpdateConfigsSimPath").Error(rErr)
		return false, t.Get(text.ConfigsErrorUpdate), http.StatusInternalServerError, nil
	}

	var m Model

	var mErr = json.Unmarshal(rb, &m)

	if mErr != nil {
		common.Log("UpdateConfigsSimPath").Error(mErr)
		return false, t.Get(text.ConfigsErrorUpdate), http.StatusInternalServerError, nil
	}

	m.SimulatorsConfigFolder = ns.SimulatorsConfigFolder

	var wb, miErr = json.MarshalIndent(m, "", " ")

	if miErr != nil {
		common.Log("UpdateConfigsSimPath").Error(miErr)
		return false, t.Get(text.ConfigsErrorUpdate), http.StatusInternalServerError, nil
	}

	// Write the configuration file with the new data
	var wErr = os.WriteFile(p, wb, fs.FileMode(common.FilePermissions))

	if wErr != nil {
		common.Log("UpdateConfigsSimPath").Error(wErr)
		return false, t.Get(text.ConfigsErrorUpdate), http.StatusInternalServerError, nil
	}

	s.SimulatorsConfigFolder = m.SimulatorsConfigFolder

	return true, "", http.StatusOK, s
}

/*
Checks if the configuration file exists or not.

If it exist return true and the file path.

If not logs a warning, returns false and the file path.
*/
func isConfigs() (is bool, c string) {
	c = common.GetThePath(configFile)

	s, err := os.Stat(c)

	if err != nil {
		common.Log("isConfigs").Warn(err)

		return false, c
	}

	if s != nil {
		return true, c
	}

	common.Log("isConfigs").Warn(translation.Translation{L: translation.EnGb}.Get(text.ConfigsNotExist))

	return false, c
}

/*
Open the configuration file, unmarshal it to a structure and return the Model.

f	-	Configuration file path (string)
*/
func openConfigsFile(f string) Model {
	bv, err := os.ReadFile(f)

	if err != nil {
		common.Log("openConfigFile").Fatal(err)
	}

	var c = unmarshalConfigsJSON(bv)

	return c
}

/*
Convert configs json file (in bytes) to the Model and returns it.

b	-	Json file converted into bytes ([]byte)
*/
func unmarshalConfigsJSON(b []byte) Model {
	var c Model

	err := json.Unmarshal(b, &c)

	if err != nil {
		common.Log("unmarshalConfigs").Fatal(err)
	}

	return c
}

/*
Create the configs file.

p	-	Configs file with path (string)
*/
func createConfigsFile(p string) {
	f, err := os.Create(p)

	if err != nil {
		common.Log("createConfigsFile").Panic(err)
	}

	f.Close()

	common.Log("createConfigsFile").Info(translation.Translation{L: translation.EnGb}.Get(text.ConfigsCreated))
}

/*
Update the new configuration file with data.

fl	-	Flags set in the execution of the application (*flags.Flags)

f	-	Full path and file name for the configuration
*/
func updateNewConfigsFile(fl *flags.Flags, f string) Model {
	var configs = Model{
		GeneralConfigFolder:    common.DefGSPath,
		SimulatorsConfigFolder: common.DefSCPath,
	}

	if fl.GCFolder != "" {
		configs.GeneralConfigFolder = fl.GCFolder
	}

	if fl.SCFolder != "" {
		configs.SimulatorsConfigFolder = fl.SCFolder
	}

	var b = marshalIndent(configs)

	var err = os.WriteFile(f, b, fs.FileMode(common.FilePermissions))

	if err != nil {
		common.Log("updateNewConfigsFile").Fatal(err)
	}

	common.Log("updateNewConfigsFile").Info(translation.Translation{L: translation.EnGb}.Get(text.ConfigsWritten))

	return configs
}

/*
Update the configuration file with new data.

Returns the model with the new data, in case of failure return the old data (Model).

fl	-	Flags set in the execution of the application (*flags.Flags)

c	-	Model with the data read from the file (Model)

fp	-	Full path and file name for the configuration (string)
*/
func updateConfigsFile(fl *flags.Flags, c Model, fp string) Model {
	var nc = Model{}

	if fl.GCFolder != "" {
		nc.GeneralConfigFolder = fl.GCFolder
	} else {
		nc.GeneralConfigFolder = c.GeneralConfigFolder
	}

	if fl.SCFolder != "" {
		nc.SimulatorsConfigFolder = fl.SCFolder
	} else {
		nc.SimulatorsConfigFolder = c.SimulatorsConfigFolder
	}

	var b = marshalIndent(nc)

	var err = os.WriteFile(fp, b, fs.FileMode(common.FilePermissions))

	if err != nil {
		common.Log("updateConfigsFile").Error(err)

		return c
	}

	common.Log("updateConfigsFile").Info(translation.Translation{L: translation.EnGb}.Get(text.ConfigsWritten))

	return nc
}

/*
Deserialize the Model to json in an byte array.

Returns []byte.

c	-	The model with the file data (Model)
*/
func marshalIndent(c Model) []byte {
	b, err := json.MarshalIndent(c, "", " ")

	if err != nil {
		common.Log("marshalIndent").Fatal(err)
	}

	return b
}

package configs

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"

	"github.com/alexandrepasc/ev-chargers-simulator-2/common"
	"github.com/alexandrepasc/ev-chargers-simulator-2/flags"
)

func Configs(fl *flags.Flags) {
	ok, fp := isConfigs()

	if ok {
		var conf = openConfigsFile(fp)

		if fl.ForceUpdate {
			conf = updateConfigsFile(fl, conf, fp)
		}

		fmt.Println(conf)
	} else {
		createConfigsFile(fp)

		updateNewConfigsFile(fl, fp)
	}
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

	common.Log("isConfigs").Warn("Configurations file doesn't exist.")

	return false, c
}

/*
Open the configuration file, unmarshal it to a structure and return the configsModel structure.

f	-	Configuration file path (string)
*/
func openConfigsFile(f string) configsModel {
	bv, err := os.ReadFile(f)

	if err != nil {
		common.Log("openConfigFile").Warn(err)
	}

	var c = unmarshalConfigsJSON(bv)

	return c
}

/*
Convert configs json file (in bytes) to the configsModel and returns it.

b	-	Json file converted into bytes ([]byte)
*/
func unmarshalConfigsJSON(b []byte) configsModel {
	var c configsModel

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

	common.Log("createConfigsFile").Info("Created configuration file")
}

/*
Update the new configuration file with data.

fl	-	Flags set in the execution of the application (*flags.Flags)

f	-	Full path and file name for the configuration
*/
func updateNewConfigsFile(fl *flags.Flags, f string) configsModel {
	var configs = configsModel{
		GeneralConfigFolder:    common.DefGSPath,
		SimulatorsConfigFolder: common.DefSCPath,
	}

	if !isGCFolderDefault(fl.GCFolder) {
		configs.GeneralConfigFolder = fl.GCFolder
	}

	if !isSCFolderDefault(fl.SCFolder) {
		configs.SimulatorsConfigFolder = fl.SCFolder
	}

	var b = marshalIndent(configs)

	var err = os.WriteFile(f, b, fs.FileMode(common.FilePermissions))

	if err != nil {
		common.Log("updateNewConfigsFile").Fatal(err)
	}

	common.Log("updateNewConfigsFile").Info("Configurations written to file")

	return configs
}

/*
Update the configuration file with new data.

Returns the model with the new data, in case of failure return the old data (configsModel).

fl	-	Flags set in the execution of the application (*flags.Flags)

c	-	Model with the data read from the file (configsModel)

fp	-	Full path and file name for the configuration (string)
*/
func updateConfigsFile(fl *flags.Flags, c configsModel, fp string) configsModel {
	var nc = configsModel{}

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

	common.Log("updateConfigsFile")

	return nc
}

/*
Check if the general configurations folder is default.

v	-	Value of the configurations folder (string)
*/
func isGCFolderDefault(v string) bool {
	return v == common.DefGSPath
}

/*
Check if the simulators configurations folder is default.

v	-	Value of the configurations folder (string)
*/
func isSCFolderDefault(v string) bool {
	return v == common.DefSCPath
}

/*
Deserialize the configsModel to json in an byte array.

Returns []byte.

c	-	The model with the file data (configsModel)
*/
func marshalIndent(c configsModel) []byte {
	b, err := json.MarshalIndent(c, "", " ")

	if err != nil {
		common.Log("marshalIndent").Fatal(err)
	}

	return b
}

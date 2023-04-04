package settings

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/alexandrepasc/ev-chargers-simulator-2/common"
	"github.com/alexandrepasc/ev-chargers-simulator-2/flags"
)

func Settings(fl *flags.Flags) {
	ok, fp := isConfigs()

	if ok {
		var conf = openConfigFile(fp)

		fmt.Println(conf)
	} else {
		createConfigsFile(fp)

		updateNewConfigsFile(fl, fp)
	}
}

func SetDefaultConfigs() {}

func SetConfigs() {}

/*
Checks if the configuration file exists or not.

If it exist return true and the file path.

If not logs a warning, returns false and the file path.
*/
func isConfigs() (is bool, c string) {
	c = common.GetThePath(configFile)

	s, err := os.Stat(c)

	if err != nil {
		return false, c
	}

	if s != nil {
		return true, c
	}

	common.Log("isConfigs").Warn("Configuration file doesn't exist.")

	return false, c
}

/*
Open the configuration file, unmarshal it to a structure and return the configsModel structure.

f	-	Configuration file path (string)
*/
func openConfigFile(f string) configsModel {
	bv, err := os.ReadFile(f)

	if err != nil {
		common.Log("openConfigFile").Warn(err)
	}

	var i = unmarshalJSON(bv)

	config, ok := i.(configsModel)

	if !ok {
		common.Log("openConfigFile").Panic("Unmarshal configs file error")
	}

	return config
}

/*
Convert a json file (in bytes) to an interface and returns it.

b	-	Json file converted into bytes ([]byte)
*/
func unmarshalJSON(b []byte) interface{} {
	var p interface{}

	err := json.Unmarshal(b, &p)

	if err != nil {
		common.Log("unmarshalConfigs").Panic(err)
	}

	return p
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
}

func updateNewConfigsFile(fl *flags.Flags, f string) {
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

	b, err := json.MarshalIndent(configs, "", " ")

	if err != nil {
		common.Log("updateNewConfigsFile").Fatal(err)
	}

	const fs = 0o600

	nok := os.WriteFile(f, b, fs)

	if nok != nil {
		common.Log("updateNewConfigsFile").Fatal(err)
	}
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

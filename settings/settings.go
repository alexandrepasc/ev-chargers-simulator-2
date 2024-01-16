package settings

import (
	"io/fs"
	"os"

	"github.com/alexandrepasc/ev-chargers-simulator-2/common"
	"github.com/alexandrepasc/ev-chargers-simulator-2/flags"
	"github.com/alexandrepasc/ev-chargers-simulator-2/settings/configs"
	"github.com/alexandrepasc/ev-chargers-simulator-2/settings/general"
	"github.com/alexandrepasc/ev-chargers-simulator-2/translation"
	"github.com/alexandrepasc/ev-chargers-simulator-2/translation/text"
)

// TODO: evaluate if in an update to the configs should only create the new folders, or remove the older ones and create the new
/*
Handles the logic to create and edit the application settings.

Returns the configurations and the general models with the data set in the files (configs.Model, general.Model)

fl	-	The flags that the user set when executing the application (*flags.Flags)
*/
func Settings(fl *flags.Flags) (c configs.Model, g general.Model) {
	var confs = configs.Configs(fl)

	checkConfigs(confs)

	var gen = general.General(confs.GeneralConfigFolder, fl)

	return confs, gen
}

/*
Check if the configs folders exist, if not will create them.

configs	-	Configs model with the data to be tested (configs.Model)
*/
func checkConfigs(confs configs.Model) {
	if !isFolder(confs.GeneralConfigFolder) {
		createFolder(confs.GeneralConfigFolder)
	}

	if !isFolder(confs.SimulatorsConfigFolder) {
		createFolder(confs.SimulatorsConfigFolder)
	}

	if !isFolder(confs.SimulatorsConfigFolder + common.DefMCFolder) {
		createFolder(confs.SimulatorsConfigFolder + common.DefMCFolder)
	}

	if !isFolder(confs.LogsConfigFolder) {
		createFolder(confs.LogsConfigFolder)
	}
}

/*
Checks if the folder exit.

p	-	Folder path to be checked (string)
*/
func isFolder(p string) bool {
	s, err := os.Stat(p)

	if err != nil {
		common.Log("isFolder").Warn(err)

		return false
	}

	if s != nil {
		return true
	}

	common.Log("isFolder").Warn(translation.Translation{L: translation.EnGb}.Get(text.FolderNotExist) + p)

	return false
}

/*
Create the folder and stops the application if the creation fail.

p	-	Folder path to be created (string)
*/
func createFolder(p string) {
	err := os.MkdirAll(p, fs.FileMode(common.FolderPermissions))

	if err != nil {
		common.Log("createFolder").Fatal(err)
	}

	common.Log("createFolder").Info(translation.Translation{L: translation.EnGb}.Get(text.FolderCreated) + p)
}

//nolint:gocritic,nolintlint
package configs_test

import (
	"os"
	"testing"

	"github.com/alexandrepasc/ev-chargers-simulator-2/common"
	"github.com/alexandrepasc/ev-chargers-simulator-2/flags"
	"github.com/alexandrepasc/ev-chargers-simulator-2/settings/configs"
	"github.com/alexandrepasc/ev-chargers-simulator-2/translation"
	"github.com/alexandrepasc/ev-chargers-simulator-2/translation/text"
	"github.com/stretchr/testify/assert"
)

const (
	file            string = "/configs.json"
	generalFolder   string = "/general"
	simulatorFolder string = "/simulator"
	logsFolder      string = "/logs"
)

func TestConfigs(t *testing.T) {
	fl := flags.Flags{
		ForceUpdate: false,
		HostAddr:    "",
		ConnTimeout: -1,
		CSAddr:      "",
		CSPort:      "",
		GCFolder:    "",
		SCFolder:    "",
	}

	e := configs.Model{
		GeneralConfigFolder:    common.DefGSPath,
		SimulatorsConfigFolder: common.DefSCPath,
		LogsConfigFolder:       common.DefLogsPath,
	}

	a := configs.Configs(&fl)

	s, err := os.Stat(file)

	assert.NotNil(t, err)

	assert.Nil(t, s)

	assert.Equal(t, e, a)

	os.Remove(common.GetThePath(file))

	configsNoFileWithFlags(t)

	os.Remove(common.GetThePath(file))

	updateConfigsNoForce(t)

	os.Remove(common.GetThePath(file))

	updateConfigsForce(t)

	os.Remove(common.GetThePath(file))

	getConfigsSimPath(t)

	os.Remove(common.GetThePath(file))

	getConfigsSimPathFailRead(t)

	os.Remove(common.GetThePath(file))

	updateConfigsSimPath(t)

	os.Remove(common.GetThePath(file))

	updateConfigsSimPathWithEndSlash(t)

	os.Remove(common.GetThePath(file))

	updateConfigsSimPathFailRead(t)
}

func configsNoFileWithFlags(t *testing.T) {
	t.Helper()

	fl := flags.Flags{
		ForceUpdate: false,
		HostAddr:    "",
		ConnTimeout: -1,
		CSAddr:      "",
		CSPort:      "",
		GCFolder:    generalFolder,
		SCFolder:    simulatorFolder,
		LogFolder:   logsFolder,
	}

	e := configs.Model{
		GeneralConfigFolder:    generalFolder,
		SimulatorsConfigFolder: simulatorFolder,
		LogsConfigFolder:       logsFolder,
	}

	a := configs.Configs(&fl)

	s, err := os.Stat(file)

	assert.NotNil(t, err)

	assert.Nil(t, s)

	assert.Equal(t, e, a)
}

func updateConfigsNoForce(t *testing.T) {
	t.Helper()

	fl := flags.Flags{
		ForceUpdate: false,
		HostAddr:    "",
		ConnTimeout: -1,
		CSAddr:      "",
		CSPort:      "",
		GCFolder:    "",
		SCFolder:    "",
		LogFolder:   "",
	}

	e := configs.Model{
		GeneralConfigFolder:    common.DefGSPath,
		SimulatorsConfigFolder: common.DefSCPath,
		LogsConfigFolder:       common.DefLogsPath,
	}

	configs.Configs(&fl)

	fl.GCFolder = generalFolder
	fl.SCFolder = simulatorFolder

	a := configs.Configs(&fl)

	assert.Equal(t, e, a)
}

func updateConfigsForce(t *testing.T) {
	t.Helper()

	fl := flags.Flags{
		ForceUpdate: false,
		HostAddr:    "",
		ConnTimeout: -1,
		CSAddr:      "",
		CSPort:      "",
		GCFolder:    "",
		SCFolder:    "",
	}

	e := configs.Model{
		GeneralConfigFolder:    generalFolder,
		SimulatorsConfigFolder: simulatorFolder,
		LogsConfigFolder:       logsFolder,
	}

	configs.Configs(&fl)

	fl.GCFolder = generalFolder
	fl.SCFolder = simulatorFolder
	fl.LogFolder = logsFolder
	fl.ForceUpdate = true

	a := configs.Configs(&fl)

	assert.Equal(t, e, a)
}

func getConfigsSimPath(t *testing.T) {
	t.Helper()

	var fl = flags.Flags{
		ForceUpdate: true,
		HostAddr:    "",
		ConnTimeout: -1,
		CSAddr:      "",
		CSPort:      "",
		GCFolder:    "",
		SCFolder:    simulatorFolder,
	}

	configs.Configs(&fl)

	var e = configs.PathsModel{
		SimulatorsConfigFolder: simulatorFolder,
	}

	var ok, msg, code, a = configs.GetConfigPaths(translation.Translation{L: translation.PtPt})

	assert.True(t, ok)

	assert.Empty(t, msg)

	assert.Equal(t, 200, code)

	assert.Equal(t, e.SimulatorsConfigFolder, a.SimulatorsConfigFolder)
}

func getConfigsSimPathFailRead(t *testing.T) {
	t.Helper()

	var ok, msg, code, a = configs.GetConfigPaths(translation.Translation{L: translation.PtPt})

	assert.False(t, ok)

	assert.Equal(t, translation.Translation{L: translation.PtPt}.Get(text.ConfigsErrorRead), msg)

	assert.Equal(t, 500, code)

	assert.Nil(t, a)
}

func updateConfigsSimPath(t *testing.T) {
	t.Helper()

	var fl = flags.Flags{
		ForceUpdate: true,
		HostAddr:    "",
		ConnTimeout: -1,
		CSAddr:      "",
		CSPort:      "",
		GCFolder:    "",
		SCFolder:    simulatorFolder,
	}

	configs.Configs(&fl)

	var e = configs.PathsModel{
		SimulatorsConfigFolder: "/test",
		LoggingConfigFolder:    "/logTest",
	}

	var ok, msg, code, a = configs.UpdateConfigPaths(&e, translation.Translation{L: translation.PtPt})

	assert.True(t, ok)

	assert.Empty(t, msg)

	assert.Equal(t, 200, code)

	assert.Equal(t, e.SimulatorsConfigFolder, a.SimulatorsConfigFolder)

	assert.Equal(t, e.LoggingConfigFolder, a.LoggingConfigFolder)

	ok, msg, code, a = configs.GetConfigPaths(translation.Translation{L: translation.PtPt})

	assert.True(t, ok)

	assert.Empty(t, msg)

	assert.Equal(t, 200, code)

	assert.Equal(t, e.SimulatorsConfigFolder, a.SimulatorsConfigFolder)
}

func updateConfigsSimPathWithEndSlash(t *testing.T) {
	t.Helper()

	var fl = flags.Flags{
		ForceUpdate: true,
		HostAddr:    "",
		ConnTimeout: -1,
		CSAddr:      "",
		CSPort:      "",
		GCFolder:    "",
		SCFolder:    simulatorFolder,
	}

	configs.Configs(&fl)

	var rb = configs.PathsModel{
		SimulatorsConfigFolder: "/test/",
		LoggingConfigFolder:    "/log/",
	}

	const esc = "/test"

	const elp = "/log"

	var ok, msg, code, a = configs.UpdateConfigPaths(&rb, translation.Translation{L: translation.PtPt})

	assert.True(t, ok)

	assert.Empty(t, msg)

	assert.Equal(t, 200, code)

	assert.Equal(t, esc, a.SimulatorsConfigFolder)

	assert.Equal(t, elp, a.LoggingConfigFolder)

	ok, msg, code, a = configs.GetConfigPaths(translation.Translation{L: translation.PtPt})

	assert.True(t, ok)

	assert.Empty(t, msg)

	assert.Equal(t, 200, code)

	assert.Equal(t, esc, a.SimulatorsConfigFolder)

	assert.Equal(t, elp, a.LoggingConfigFolder)
}

func updateConfigsSimPathFailRead(t *testing.T) {
	t.Helper()

	var e = configs.PathsModel{
		SimulatorsConfigFolder: "/test",
		LoggingConfigFolder:    "/log",
	}

	var ok, msg, code, a = configs.UpdateConfigPaths(&e, translation.Translation{L: translation.PtPt})

	assert.False(t, ok)

	assert.Equal(t, translation.Translation{L: translation.PtPt}.Get(text.ConfigsErrorRead), msg)

	assert.Equal(t, 500, code)

	assert.Nil(t, a)
}

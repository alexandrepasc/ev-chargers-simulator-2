//nolint:gocritic,nolintlint,errcheck
package general_test

import (
	"io/fs"
	"os"
	"testing"

	"github.com/alexandrepasc/ev-chargers-simulator-2/common"
	"github.com/alexandrepasc/ev-chargers-simulator-2/flags"
	"github.com/alexandrepasc/ev-chargers-simulator-2/settings/general"
	"github.com/stretchr/testify/assert"
)

const (
	file string = "/general.json"
)

func TestGeneral(t *testing.T) {
	os.Mkdir(common.DefGSPath, fs.FileMode(common.FolderPermissions))

	fl := flags.Flags{
		ForceUpdate: false,
		HostAddr:    "",
		ConnTimeout: -1,
		CSAddr:      "",
		CSPort:      "",
		GCFolder:    "",
		SCFolder:    "",
		Language:    "",
		APIAddr:     "",
		APIPort:     "",
	}

	e := general.Model{
		HostIP:      common.DefSimIP,
		ConnTimeout: common.DefTimeout,
		CSAddr:      common.DefCSIP,
		CSPort:      common.DefCSPort,
		Lang:        common.DefLanguage,
		APIAddr:     common.DefAPIAddr,
		APIPort:     common.DefAPIPort,
	}

	a := general.General(common.DefGSPath, &fl)

	s, err := os.Stat(common.DefGSPath + file)

	assert.Nil(t, err)

	assert.NotNil(t, s)

	assert.Equal(t, e, a)

	os.Remove(common.DefGSPath + file)

	generalNoFileWithFlags(t)

	os.Remove(common.DefGSPath + file)

	updateGeneralNoForce(t)

	os.Remove(common.DefGSPath + file)

	updateGeneralForce(t)
}

func generalNoFileWithFlags(t *testing.T) {
	t.Helper()

	fl := flags.Flags{
		ForceUpdate: false,
		HostAddr:    "123.123.123.123",
		ConnTimeout: 100,
		CSAddr:      "234.234.234.234",
		CSPort:      "444",
		GCFolder:    "",
		SCFolder:    "",
		Language:    "pt-PT",
		APIAddr:     "123.123.123.123",
		APIPort:     "999",
	}

	e := general.Model{
		HostIP:      "123.123.123.123",
		ConnTimeout: 100,
		CSAddr:      "234.234.234.234",
		CSPort:      "444",
		Lang:        "pt-PT",
		APIAddr:     "123.123.123.123",
		APIPort:     "999",
	}

	a := general.General(common.DefGSPath, &fl)

	s, err := os.Stat(common.DefGSPath + file)

	assert.Nil(t, err)

	assert.NotNil(t, s)

	assert.Equal(t, e, a)
}

func updateGeneralNoForce(t *testing.T) {
	t.Helper()

	fl := flags.Flags{
		ForceUpdate: false,
		HostAddr:    "",
		ConnTimeout: -1,
		CSAddr:      "",
		CSPort:      "",
		GCFolder:    "",
		SCFolder:    "",
		Language:    "",
		APIAddr:     "",
		APIPort:     "",
	}

	general.General(common.DefGSPath, &fl)

	fl.HostAddr = "123.123.123.123"
	fl.ConnTimeout = 400
	fl.CSAddr = "234.234.234.234"
	fl.CSPort = "444"
	fl.Language = "pt-PT"
	fl.APIAddr = "123.123.123.123"
	fl.APIPort = "999"

	e := general.Model{
		HostIP:      "123.123.123.123",
		ConnTimeout: 400,
		CSAddr:      "234.234.234.234",
		CSPort:      "444",
		Lang:        "pt-PT",
		APIAddr:     "123.123.123.123",
		APIPort:     "999",
	}

	a := general.General(common.DefGSPath, &fl)

	s, err := os.Stat(common.DefGSPath + file)

	assert.Nil(t, err)

	assert.NotNil(t, s)

	assert.NotEqual(t, e, a)
}

func updateGeneralForce(t *testing.T) {
	t.Helper()

	fl := flags.Flags{
		ForceUpdate: false,
		HostAddr:    "",
		ConnTimeout: -1,
		CSAddr:      "",
		CSPort:      "",
		GCFolder:    "",
		SCFolder:    "",
		Language:    "",
		APIAddr:     "",
		APIPort:     "",
	}

	general.General(common.DefGSPath, &fl)

	fl.ForceUpdate = true
	fl.HostAddr = "123.123.123.123"
	fl.ConnTimeout = 400
	fl.CSAddr = "234.234.234.234"
	fl.CSPort = "444"
	fl.Language = "pt-PT"
	fl.APIAddr = "123.123.123.123"
	fl.APIPort = "999"

	e := general.Model{
		HostIP:      "123.123.123.123",
		ConnTimeout: 400,
		CSAddr:      "234.234.234.234",
		CSPort:      "444",
		Lang:        "pt-PT",
		APIAddr:     "123.123.123.123",
		APIPort:     "999",
	}

	a := general.General(common.DefGSPath, &fl)

	s, err := os.Stat(common.DefGSPath + file)

	assert.Nil(t, err)

	assert.NotNil(t, s)

	assert.Equal(t, e, a)
}

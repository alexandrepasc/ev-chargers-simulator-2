//nolint:gocritic,nolintlint,errcheck,goconst
package general_test

import (
	"io/fs"
	"os"
	"testing"

	"github.com/alexandrepasc/ev-chargers-simulator-2/common"
	"github.com/alexandrepasc/ev-chargers-simulator-2/flags"
	"github.com/alexandrepasc/ev-chargers-simulator-2/settings/general"
	"github.com/alexandrepasc/ev-chargers-simulator-2/translation"
	"github.com/alexandrepasc/ev-chargers-simulator-2/translation/text"
	"github.com/stretchr/testify/assert"
)

const (
	file        string = "/general.json"
	defGSFolder string = "/settings"
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

func TestGetGeneral(t *testing.T) {
	tmp := t.TempDir()

	os.Mkdir(tmp+defGSFolder, fs.FileMode(common.FolderPermissions))

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

	general.General(tmp+defGSFolder, &fl)

	ok, msg, code, g := general.GetGeneralConf(tmp+defGSFolder, translation.Translation{L: translation.EnGb})

	e := general.Model{
		HostIP:      common.DefSimIP,
		ConnTimeout: common.DefTimeout,
		CSAddr:      common.DefCSIP,
		CSPort:      common.DefCSPort,
		Lang:        common.DefLanguage,
		APIAddr:     common.DefAPIAddr,
		APIPort:     common.DefAPIPort,
	}

	assert.True(t, ok)

	assert.Equal(t, "", msg)

	assert.Equal(t, 200, code)

	assert.Equal(t, e.HostIP, g.HostIP)

	assert.Equal(t, e.ConnTimeout, g.ConnTimeout)

	assert.Equal(t, e.CSAddr, g.CSAddr)

	assert.Equal(t, e.CSPort, g.CSPort)

	assert.Equal(t, e.Lang, g.Lang)

	assert.Equal(t, e.APIAddr, g.APIAddr)

	assert.Equal(t, e.APIPort, g.APIPort)
}

func TestGetGeneralFailRead(t *testing.T) {
	ok, msg, code, g := general.GetGeneralConf(defGSFolder, translation.Translation{L: translation.EnGb})

	assert.False(t, ok)

	assert.Equal(t, translation.Translation{L: translation.EnGb}.Get(text.GeneralErrorRead), msg)

	assert.Equal(t, 500, code)

	assert.Nil(t, g)
}

func TestUpdateGeneral(t *testing.T) {
	tmp := t.TempDir()

	os.Mkdir(tmp+defGSFolder, fs.FileMode(common.FolderPermissions))

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

	general.General(tmp+defGSFolder, &fl)

	e := general.Model{
		HostIP:      "HostIP",
		ConnTimeout: 999,
		CSAddr:      "CSAddr",
		CSPort:      "CSPort",
		Lang:        "Lang",
		APIAddr:     "APIAddr",
		APIPort:     "APIPort",
	}

	ok, msg, code, a := general.UpdateGeneralConf(&e, tmp+defGSFolder, translation.Translation{L: translation.EnGb})

	assert.True(t, ok)

	assert.Empty(t, msg)

	assert.Equal(t, 200, code)

	assert.Equal(t, e.HostIP, a.HostIP)

	assert.Equal(t, e.ConnTimeout, a.ConnTimeout)

	assert.Equal(t, e.CSAddr, a.CSAddr)

	assert.Equal(t, e.CSPort, a.CSPort)

	assert.Equal(t, e.Lang, a.Lang)

	assert.Equal(t, e.APIAddr, a.APIAddr)

	assert.Equal(t, e.APIPort, a.APIPort)
}

func TestUpdateGeneralFailUpdate(t *testing.T) {
	tmp := t.TempDir()

	os.Mkdir(tmp+defGSFolder, fs.FileMode(common.FolderPermissions))

	e := general.Model{
		HostIP:      "HostIP",
		ConnTimeout: 999,
		CSAddr:      "CSAddr",
		CSPort:      "CSPort",
		Lang:        "Lang",
		APIAddr:     "APIAddr",
		APIPort:     "APIPort",
	}

	ok, msg, code, a := general.UpdateGeneralConf(&e, defGSFolder, translation.Translation{L: translation.EnGb})

	assert.False(t, ok)

	assert.Equal(t, translation.Translation{L: translation.EnGb}.Get(text.GeneralErrorUpdate), msg)

	assert.Equal(t, 500, code)

	assert.Nil(t, a)
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

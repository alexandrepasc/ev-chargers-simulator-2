package config_test

import (
	"bytes"
	"encoding/json"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/alexandrepasc/ev-chargers-simulator-2/api/config"
	"github.com/alexandrepasc/ev-chargers-simulator-2/api/errors"
	"github.com/alexandrepasc/ev-chargers-simulator-2/common"
	"github.com/alexandrepasc/ev-chargers-simulator-2/flags"
	"github.com/alexandrepasc/ev-chargers-simulator-2/settings/general"
	"github.com/alexandrepasc/ev-chargers-simulator-2/translation"
	"github.com/alexandrepasc/ev-chargers-simulator-2/translation/text"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type method string

const (
	defGSFolder string = "/settings"
	ep          string = "/configs/general"
	mGet        method = "GET"
	mPut        method = "PUT"
	// mPost       method = "POST"
	// mDelete     method = "DELETE"
)

func TestGetGeneral(t *testing.T) {
	r, w, c := before(t)

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

	general.General(c.Path, &fl)

	e := general.Model{
		HostIP:      common.DefSimIP,
		ConnTimeout: common.DefTimeout,
		CSAddr:      common.DefCSIP,
		CSPort:      common.DefCSPort,
		Lang:        common.DefLanguage,
		APIAddr:     common.DefAPIAddr,
		APIPort:     common.DefAPIPort,
	}

	c.Configs(r)

	req, _ := http.NewRequest(string(mGet), ep, http.NoBody)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	a := general.Model{}
	json.Unmarshal(w.Body.Bytes(), &a) //nolint:errcheck // because test

	assert.Equal(t, e.HostIP, a.HostIP)

	assert.Equal(t, e.ConnTimeout, a.ConnTimeout)

	assert.Equal(t, e.CSAddr, a.CSAddr)

	assert.Equal(t, e.CSPort, a.CSPort)

	assert.Equal(t, e.Lang, a.Lang)

	assert.Equal(t, e.APIAddr, a.APIAddr)

	assert.Equal(t, e.APIPort, a.APIPort)
}

func TestGetGeneralNoFile(t *testing.T) {
	r, w, c := before(t)

	c.Configs(r)

	req, _ := http.NewRequest(string(mGet), ep, http.NoBody)

	r.ServeHTTP(w, req)

	e := errors.ErroMsg{
		Message: c.Lang.Get(text.GeneralErrorRead),
	}

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	a := errors.ErroMsg{}
	json.Unmarshal(w.Body.Bytes(), &a) //nolint:errcheck // because test

	assert.Equal(t, e.Message, a.Message)
}

func TestGetGeneralNoFolder(t *testing.T) {
	r, w, c := before(t)

	os.Remove(c.Path)

	c.Configs(r)

	req, _ := http.NewRequest(string(mGet), ep, http.NoBody)

	r.ServeHTTP(w, req)

	e := errors.ErroMsg{
		Message: c.Lang.Get(text.GeneralErrorRead),
	}

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	a := errors.ErroMsg{}
	json.Unmarshal(w.Body.Bytes(), &a) //nolint:errcheck // because test

	assert.Equal(t, e.Message, a.Message)
}

func TestPutGeneral(t *testing.T) {
	r, w, c := before(t)

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

	general.General(c.Path, &fl)

	e := general.Model{
		HostIP:      "test",
		ConnTimeout: 123,
		CSAddr:      "test",
		CSPort:      "test",
		Lang:        "pt-pt",
		APIAddr:     "test",
		APIPort:     "test",
	}

	j, _ := json.Marshal(e)
	b := bytes.NewReader(j)

	c.Configs(r)

	req, _ := http.NewRequest(string(mPut), ep, b)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	a := general.Model{}
	json.Unmarshal(w.Body.Bytes(), &a) //nolint:errcheck // because test

	assert.Equal(t, e.HostIP, a.HostIP)

	assert.Equal(t, e.ConnTimeout, a.ConnTimeout)

	assert.Equal(t, e.CSAddr, a.CSAddr)

	assert.Equal(t, e.CSPort, a.CSPort)

	assert.Equal(t, e.Lang, a.Lang)

	assert.Equal(t, e.APIAddr, a.APIAddr)

	assert.Equal(t, e.APIPort, a.APIPort)
}

func TestPutGeneralBadRequestBody(t *testing.T) {
	r, w, c := before(t)

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

	general.General(c.Path, &fl)

	rb := `{
		"wrong": "yes"
	}`

	j, _ := json.Marshal(rb)
	b := bytes.NewReader(j)

	c.Configs(r)

	req, _ := http.NewRequest(string(mPut), ep, b)

	r.ServeHTTP(w, req)

	e := errors.ErroMsg{
		Message: c.Lang.Get(text.RequestBodyDoesntMatch),
	}

	assert.Equal(t, http.StatusBadRequest, w.Code)

	a := errors.ErroMsg{}
	json.Unmarshal(w.Body.Bytes(), &a) //nolint:errcheck // because test

	assert.Equal(t, e.Message, a.Message)
}

func TestPutGeneralNoFile(t *testing.T) {
	r, w, c := before(t)

	e := general.Model{
		HostIP:      "test",
		ConnTimeout: 123,
		CSAddr:      "test",
		CSPort:      "test",
		Lang:        "pt-pt",
		APIAddr:     "test",
		APIPort:     "test",
	}

	j, _ := json.Marshal(e)
	b := bytes.NewReader(j)

	c.Configs(r)

	req, _ := http.NewRequest(string(mPut), ep, b)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	a := general.Model{}
	json.Unmarshal(w.Body.Bytes(), &a) //nolint:errcheck // because test

	assert.Equal(t, e.HostIP, a.HostIP)

	assert.Equal(t, e.ConnTimeout, a.ConnTimeout)

	assert.Equal(t, e.CSAddr, a.CSAddr)

	assert.Equal(t, e.CSPort, a.CSPort)

	assert.Equal(t, e.Lang, a.Lang)

	assert.Equal(t, e.APIAddr, a.APIAddr)

	assert.Equal(t, e.APIPort, a.APIPort)
}

func TestPutGeneralNoFolder(t *testing.T) {
	r, w, c := before(t)

	os.Remove(c.Path)

	e := errors.ErroMsg{
		Message: c.Lang.Get(text.GeneralErrorUpdate),
	}

	rb := general.Model{
		HostIP:      "test",
		ConnTimeout: 123,
		CSAddr:      "test",
		CSPort:      "test",
		Lang:        "pt-pt",
		APIAddr:     "test",
		APIPort:     "test",
	}

	j, _ := json.Marshal(rb)
	b := bytes.NewReader(j)

	c.Configs(r)

	req, _ := http.NewRequest(string(mPut), ep, b)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	a := errors.ErroMsg{}
	json.Unmarshal(w.Body.Bytes(), &a) //nolint:errcheck // because test

	assert.Equal(t, e.Message, a.Message)
}

func before(t *testing.T) (r *gin.Engine, w *httptest.ResponseRecorder, c config.Configs) {
	t.Helper()

	r = gin.Default()

	w = httptest.NewRecorder()

	tmp := t.TempDir()

	os.Mkdir(tmp+defGSFolder, fs.FileMode(common.FolderPermissions)) //nolint:errcheck // because test

	c = config.Configs{
		Lang: translation.Translation{L: translation.PtPt},
		Path: tmp + defGSFolder,
	}

	return r, w, c
}

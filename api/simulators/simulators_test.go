//nolint:errcheck //because they are tests
package simulators_test

import (
	"encoding/json"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/alexandrepasc/ev-chargers-simulator-2/api/simulators"
	"github.com/alexandrepasc/ev-chargers-simulator-2/common"
	"github.com/alexandrepasc/ev-chargers-simulator-2/simulator"
	"github.com/alexandrepasc/ev-chargers-simulator-2/simulator/model"
	"github.com/alexandrepasc/ev-chargers-simulator-2/translation"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

const (
	ep       string = "/simulators"
	modelsEp string = ep + "/models"
)

func TestGetSimulators(t *testing.T) {
	r := gin.Default()

	al := []simulator.Asset{
		{
			Name:        "a_name",
			Type:        simulator.Evc,
			Protocol:    simulator.Ocpp201,
			Model:       "a",
			Phases:      simulator.One,
			CurrentType: simulator.Ac,
			Evses:       []simulator.Evse{},
		},
		{
			Name:        "b_name",
			Type:        simulator.Evc,
			Protocol:    simulator.Ocpp16,
			Model:       "b",
			Phases:      simulator.One,
			CurrentType: simulator.Ac,
			Evses:       []simulator.Evse{},
		},
	}
	ml := []model.OcppModel{}

	sim := before(t, al, ml)

	s := simulators.Simulators{
		Sim: sim,
	}

	s.Simulators(r)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", ep, http.NoBody)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	a := simulators.GetSimulatorsTemp{}

	json.Unmarshal(w.Body.Bytes(), &a)

	assert.Equal(t, int64(2), a.Total)

	for i := 0; i < len(al); i++ {
		assert.Equal(t, string(al[i].Type), a.Assets.([]interface{})[i].(map[string]interface{})["type"])

		assert.Equal(t, string(al[i].Protocol), a.Assets.([]interface{})[i].(map[string]interface{})["protocol"])

		assert.Equal(t, al[i].Model, a.Assets.([]interface{})[i].(map[string]interface{})["model"])
	}
}

func TestGetSimulatorsNoAssets(t *testing.T) {
	r := gin.Default()

	al := []simulator.Asset{}
	ml := []model.OcppModel{}

	sim := before(t, al, ml)

	s := simulators.Simulators{
		Sim: sim,
	}

	s.Simulators(r)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", ep, http.NoBody)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	a := simulators.GetSimulatorsTemp{}

	json.Unmarshal(w.Body.Bytes(), &a)

	assert.Equal(t, int64(0), a.Total)

	assert.Nil(t, a.Assets)
}

func TestGetSimModels(t *testing.T) {
	r := gin.Default()

	al := []simulator.Asset{}
	ml := []model.OcppModel{
		{
			SerialNumb: "123qwe",
			Model:      "name",
			Vendor:     "vendor",
			FwVersion:  "0.0.0.0",
		},
		{
			SerialNumb: "456asd",
			Model:      "nop",
			Vendor:     "nop",
			FwVersion:  "9.9.9.9",
		},
	}

	sim := before(t, al, ml)

	sim.Al = al
	sim.Oml = ml

	s := simulators.Simulators{
		Sim: sim,
	}

	s.Simulators(r)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", modelsEp, http.NoBody)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	a := simulators.GetModelsTemp{}

	json.Unmarshal(w.Body.Bytes(), &a)

	assert.Equal(t, int64(2), a.Total)

	for i := 0; i < len(ml); i++ {
		assert.Equal(t, ml[i].SerialNumb, a.Models.([]interface{})[i].(map[string]interface{})["serialNumb"])

		assert.Equal(t, ml[i].Model, a.Models.([]interface{})[i].(map[string]interface{})["model"])

		assert.Equal(t, ml[i].Vendor, a.Models.([]interface{})[i].(map[string]interface{})["vendor"])

		assert.Equal(t, ml[i].FwVersion, a.Models.([]interface{})[i].(map[string]interface{})["fwVersion"])
	}
}

func TestGetSimModelsNoModels(t *testing.T) {
	r := gin.Default()

	al := []simulator.Asset{}
	ml := []model.OcppModel{}

	sim := before(t, al, ml)

	sim.Al = al
	sim.Oml = ml

	s := simulators.Simulators{
		Sim: sim,
	}

	s.Simulators(r)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", modelsEp, http.NoBody)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	a := simulators.GetModelsTemp{}

	json.Unmarshal(w.Body.Bytes(), &a)

	assert.Equal(t, int64(0), a.Total)

	assert.Nil(t, a.Models)
}

// func TestPostSimulators(t *testing.T) {
// 	r := gin.Default()

// 	al := []simulator.Asset{}
// 	ml := []model.OcppModel{}

// 	s := simulators.Simulators{
// 		Al:  al,
// 		Oml: ml,
// 	}

// 	s.Simulators(r)

// 	w := httptest.NewRecorder()
// }

func before(t *testing.T, al []simulator.Asset, oml []model.OcppModel) simulator.Simulator {
	t.Helper()

	tmp := t.TempDir()

	os.Mkdir(tmp+"/simConf", fs.FileMode(common.FolderPermissions))

	os.Mkdir(tmp+"/simConf"+common.DefMCFolder, fs.FileMode(common.FolderPermissions))

	generateConfFiles(tmp, al, oml)

	l := translation.Translation{
		L: translation.Translation{}.GetKey("en-GB"),
	}

	s := simulator.Simulator{
		Scp: tmp + "/simConf",
		L:   l,
	}

	return s
}

func generateConfFiles(tmp string, al []simulator.Asset, oml []model.OcppModel) {
	for i := 0; i < len(al); i++ {
		p := tmp + "/simConf" + "/" + al[i].Name + ".json"

		f, _ := os.Create(p)
		f.Close()

		b, _ := json.MarshalIndent(al[i], "", " ")

		os.WriteFile(p, b, fs.FileMode(common.FilePermissions))
	}

	for i := 0; i < len(oml); i++ {
		p := tmp + "/simConf" + common.DefMCFolder + "/model" + strconv.Itoa(i) + "_ocpp.json"

		f, _ := os.Create(p)
		f.Close()

		b, _ := json.MarshalIndent(oml[i], "", " ")

		os.WriteFile(p, b, fs.FileMode(common.FilePermissions))
	}
}

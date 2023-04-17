//nolint:errcheck,dupl //because they are tests
package simulators_test

import (
	"bytes"
	"encoding/json"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/alexandrepasc/ev-chargers-simulator-2/api/errors"
	"github.com/alexandrepasc/ev-chargers-simulator-2/api/simulators"
	"github.com/alexandrepasc/ev-chargers-simulator-2/common"
	"github.com/alexandrepasc/ev-chargers-simulator-2/simulator"
	"github.com/alexandrepasc/ev-chargers-simulator-2/simulator/model"
	"github.com/alexandrepasc/ev-chargers-simulator-2/translation"
	"github.com/alexandrepasc/ev-chargers-simulator-2/translation/text"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

const (
	ep       string = "/simulators"
	modelsEp string = ep + "/models"
)

func TestGetSimulators(t *testing.T) {
	r := gin.Default()

	al := []*simulator.Asset{
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

	sim, _ := before(t, al, ml)

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

	al := []*simulator.Asset{}
	ml := []model.OcppModel{}

	sim, _ := before(t, al, ml)

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

	al := []*simulator.Asset{}
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

	sim, _ := before(t, al, ml)

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

	al := []*simulator.Asset{}
	ml := []model.OcppModel{}

	sim, _ := before(t, al, ml)

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

func TestPostSimulators(t *testing.T) {
	r := gin.Default()

	al := []*simulator.Asset{}
	ml := []model.OcppModel{}

	sim, tmp := before(t, al, ml)

	sim.Al = al
	sim.Oml = ml

	s := simulators.Simulators{
		Sim: sim,
	}

	body := simulator.Asset{
		Name:          "name",
		Type:          simulator.Evc,
		Protocol:      simulator.Modbus,
		StartCharging: true,
		Phases:        simulator.One,
		CurrentType:   simulator.Ac,
		Evses: []simulator.Evse{
			{
				ID: 1,
				Connectors: []simulator.Connector{
					{
						ID: 1,
					},
				},
			},
		},
	}

	j, _ := json.Marshal(body)
	b := bytes.NewReader(j)

	s.Simulators(r)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("POST", ep, b)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	a := simulator.Asset{}
	json.Unmarshal(w.Body.Bytes(), &a)

	_, err := uuid.Parse(a.SimID.String())
	assert.Nil(t, err)

	assert.Equal(t, body.Name, a.Name)

	assert.Equal(t, body.Type, a.Type)

	assert.Equal(t, body.Protocol, a.Protocol)

	assert.Equal(t, body.Phases, a.Phases)

	assert.Equal(t, body.CurrentType, a.CurrentType)

	assert.Equal(t, body.Evses, a.Evses)

	assert.FileExists(t, tmp+"/simConf/"+body.Name+".json")
}

func TestPostSimulatorsRequiredFields(t *testing.T) {
	r := gin.Default()

	al := []*simulator.Asset{}
	ml := []model.OcppModel{}

	sim, _ := before(t, al, ml)

	sim.Al = al
	sim.Oml = ml

	s := simulators.Simulators{
		Sim: sim,
	}

	bn, bt, bp, bst, bph, bc, be := postSimulatorsRequiredFields()

	// name
	j, _ := json.Marshal(bn)
	b := bytes.NewReader(j)

	s.Simulators(r)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("POST", ep, b)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	// type
	j, _ = json.Marshal(bt)
	b = bytes.NewReader(j)

	req, _ = http.NewRequest("POST", ep, b)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	// protocol
	j, _ = json.Marshal(bp)
	b = bytes.NewReader(j)

	req, _ = http.NewRequest("POST", ep, b)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	// start charging
	j, _ = json.Marshal(bst)
	b = bytes.NewReader(j)

	req, _ = http.NewRequest("POST", ep, b)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	// phases
	j, _ = json.Marshal(bph)
	b = bytes.NewReader(j)

	req, _ = http.NewRequest("POST", ep, b)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	// current type
	j, _ = json.Marshal(bc)
	b = bytes.NewReader(j)

	req, _ = http.NewRequest("POST", ep, b)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	// evses
	j, _ = json.Marshal(be)
	b = bytes.NewReader(j)

	req, _ = http.NewRequest("POST", ep, b)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPostSimulatorsBadRequestBody(t *testing.T) {
	r := gin.Default()

	al := []*simulator.Asset{}
	ml := []model.OcppModel{}

	sim, _ := before(t, al, ml)

	sim.Al = al
	sim.Oml = ml

	s := simulators.Simulators{
		Sim:  sim,
		Lang: translation.Translation{L: translation.EnGb},
	}

	j, _ := json.Marshal("{nothing: atall}")
	b := bytes.NewReader(j)

	s.Simulators(r)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("POST", ep, b)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	a := errors.ErroMsg{}
	json.Unmarshal(w.Body.Bytes(), &a)

	assert.Equal(t, translation.Translation{L: translation.EnGb}.Get(text.RequestBodyDoesntMatch), a.Message)
}

func TestNotAblePostSimulatorsSameName(t *testing.T) {
	r := gin.Default()

	al := []*simulator.Asset{
		{
			Name: "test1",
		},
	}
	ml := []model.OcppModel{}

	sim, _ := before(t, al, ml)

	sim.Al = al
	sim.Oml = ml

	s := simulators.Simulators{
		Sim:  sim,
		Lang: translation.Translation{L: translation.EnGb},
	}

	body := simulator.Asset{
		Name:          "test1",
		Type:          simulator.Evc,
		Protocol:      simulator.Modbus,
		StartCharging: true,
		Phases:        simulator.One,
		CurrentType:   simulator.Ac,
		Evses: []simulator.Evse{
			{
				ID: 1,
				Connectors: []simulator.Connector{
					{
						ID: 1,
					},
				},
			},
		},
	}

	j, _ := json.Marshal(body)
	b := bytes.NewReader(j)

	s.Simulators(r)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("POST", ep, b)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	a := errors.ErroMsg{}
	json.Unmarshal(w.Body.Bytes(), &a)

	assert.Equal(t, translation.Translation{L: translation.EnGb}.Get(text.CreateSimConfFileNameExists), a.Message)
}

func TestPutSimulators(t *testing.T) {
	r := gin.Default()

	al := []*simulator.Asset{
		{
			SimID:         uuid.New(),
			Name:          "test1",
			Type:          simulator.Evc,
			Protocol:      simulator.Ocpp201,
			StartCharging: true,
			Phases:        simulator.One,
			CurrentType:   simulator.Dc,
			Evses:         []simulator.Evse{},
		},
	}
	ml := []model.OcppModel{}

	sim, _ := before(t, al, ml)

	sim.Al = al
	sim.Oml = ml

	s := simulators.Simulators{
		Sim:  sim,
		Lang: translation.Translation{L: translation.EnGb},
	}

	e := simulator.Asset{
		Type:          simulator.Pm,
		Protocol:      simulator.Ocpp16,
		StartCharging: false,
		Phases:        simulator.Three,
		CurrentType:   simulator.Ac,
		Evses:         []simulator.Evse{},
	}

	j, _ := json.Marshal(e)
	b := bytes.NewReader(j)

	s.Simulators(r)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("PUT", ep+"/"+al[0].SimID.String(), b)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	a := simulator.Asset{}
	json.Unmarshal(w.Body.Bytes(), &a)

	assert.Equal(t, al[0].SimID, a.SimID)

	assert.Equal(t, al[0].Name, a.Name)

	assert.Equal(t, e.Type, a.Type)

	assert.Equal(t, e.Protocol, a.Protocol)

	assert.Equal(t, e.StartCharging, a.StartCharging)

	assert.Equal(t, e.Phases, a.Phases)

	assert.Equal(t, e.CurrentType, a.CurrentType)
}

func TestPutSimulatorsWrongID(t *testing.T) {
	r := gin.Default()

	al := []*simulator.Asset{
		{
			SimID:         uuid.New(),
			Name:          "test1",
			Type:          simulator.Evc,
			Protocol:      simulator.Ocpp201,
			StartCharging: true,
			Phases:        simulator.One,
			CurrentType:   simulator.Dc,
			Evses:         []simulator.Evse{},
		},
	}
	ml := []model.OcppModel{}

	sim, _ := before(t, al, ml)

	sim.Al = al
	sim.Oml = ml

	s := simulators.Simulators{
		Sim:  sim,
		Lang: translation.Translation{L: translation.EnGb},
	}

	e := simulator.Asset{
		Type:          simulator.Pm,
		Protocol:      simulator.Ocpp16,
		StartCharging: false,
		Phases:        simulator.Three,
		CurrentType:   simulator.Ac,
		Evses:         []simulator.Evse{},
	}

	j, _ := json.Marshal(e)
	b := bytes.NewReader(j)

	s.Simulators(r)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("PUT", ep+"/aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa0", b)

	r.ServeHTTP(w, req)

	a := errors.ErroMsg{}
	json.Unmarshal(w.Body.Bytes(), &a)

	assert.Equal(t, http.StatusNotFound, w.Code)

	assert.Equal(t, sim.L.Get(text.UpdateSimConfFileNotFoud), a.Message)
}

func TestPutSimulatorsInvalidID(t *testing.T) {
	r := gin.Default()

	al := []*simulator.Asset{
		{
			SimID:         uuid.New(),
			Name:          "test1",
			Type:          simulator.Evc,
			Protocol:      simulator.Ocpp201,
			StartCharging: true,
			Phases:        simulator.One,
			CurrentType:   simulator.Dc,
			Evses:         []simulator.Evse{},
		},
	}
	ml := []model.OcppModel{}

	sim, _ := before(t, al, ml)

	sim.Al = al
	sim.Oml = ml

	s := simulators.Simulators{
		Sim:  sim,
		Lang: translation.Translation{L: translation.EnGb},
	}

	e := simulator.Asset{
		Type:          simulator.Pm,
		Protocol:      simulator.Ocpp16,
		StartCharging: false,
		Phases:        simulator.Three,
		CurrentType:   simulator.Ac,
		Evses:         []simulator.Evse{},
	}

	j, _ := json.Marshal(e)
	b := bytes.NewReader(j)

	s.Simulators(r)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("PUT", ep+"/asd", b)

	r.ServeHTTP(w, req)

	a := errors.ErroMsg{}
	json.Unmarshal(w.Body.Bytes(), &a)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	assert.Equal(t, sim.L.Get(text.UUIDParsingError), a.Message)
}

func TestPutSimulatorsChangeNameID(t *testing.T) {
	r := gin.Default()

	al := []*simulator.Asset{
		{
			SimID:         uuid.New(),
			Name:          "test1",
			Type:          simulator.Evc,
			Protocol:      simulator.Ocpp201,
			StartCharging: true,
			Phases:        simulator.One,
			CurrentType:   simulator.Dc,
			Evses:         []simulator.Evse{},
		},
	}
	ml := []model.OcppModel{}

	sim, _ := before(t, al, ml)

	sim.Al = al
	sim.Oml = ml

	s := simulators.Simulators{
		Sim:  sim,
		Lang: translation.Translation{L: translation.EnGb},
	}

	eid, _ := uuid.Parse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa9")
	e := simulator.Asset{
		SimID:         eid,
		Name:          "testing",
		Type:          simulator.Pm,
		Protocol:      simulator.Ocpp16,
		StartCharging: false,
		Phases:        simulator.Three,
		CurrentType:   simulator.Ac,
		Evses:         []simulator.Evse{},
	}

	j, _ := json.Marshal(e)
	b := bytes.NewReader(j)

	s.Simulators(r)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("PUT", ep+"/"+al[0].SimID.String(), b)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	a := simulator.Asset{}
	json.Unmarshal(w.Body.Bytes(), &a)

	assert.Equal(t, al[0].SimID, a.SimID)

	assert.Equal(t, al[0].Name, a.Name)

	assert.Equal(t, e.Type, a.Type)

	assert.Equal(t, e.Protocol, a.Protocol)

	assert.Equal(t, e.StartCharging, a.StartCharging)

	assert.Equal(t, e.Phases, a.Phases)

	assert.Equal(t, e.CurrentType, a.CurrentType)
}

func TestPutSimulatorsBadRequestBody(t *testing.T) {
	r := gin.Default()

	al := []*simulator.Asset{
		{
			SimID:         uuid.New(),
			Name:          "test1",
			Type:          simulator.Evc,
			Protocol:      simulator.Ocpp201,
			StartCharging: true,
			Phases:        simulator.One,
			CurrentType:   simulator.Dc,
			Evses:         []simulator.Evse{},
		},
	}
	ml := []model.OcppModel{}

	sim, _ := before(t, al, ml)

	sim.Al = al
	sim.Oml = ml

	s := simulators.Simulators{
		Sim:  sim,
		Lang: translation.Translation{L: translation.EnGb},
	}

	e := simulator.Asset{
		Type: simulator.Pm,
	}

	j, _ := json.Marshal(e)
	b := bytes.NewReader(j)

	s.Simulators(r)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("PUT", ep+"/"+al[0].SimID.String(), b)

	r.ServeHTTP(w, req)

	a := errors.ErroMsg{}
	json.Unmarshal(w.Body.Bytes(), &a)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	assert.Equal(t, sim.L.Get(text.RequestBodyDoesntMatch), a.Message)
}

func before(t *testing.T, al []*simulator.Asset, oml []model.OcppModel) (s simulator.Simulator, tmp string) {
	t.Helper()

	tmp = t.TempDir()

	os.Mkdir(tmp+"/simConf", fs.FileMode(common.FolderPermissions))

	os.Mkdir(tmp+"/simConf"+common.DefMCFolder, fs.FileMode(common.FolderPermissions))

	generateConfFiles(tmp, al, oml)

	l := translation.Translation{
		L: translation.Translation{}.GetKey("en-GB"),
	}

	s = simulator.Simulator{
		Scp: tmp + "/simConf",
		L:   l,
	}

	return s, tmp
}

func generateConfFiles(tmp string, al []*simulator.Asset, oml []model.OcppModel) {
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

func postSimulatorsRequiredFields() (bn, bt, bp, bst, bph, bc, be simulator.Asset) { //nolint:gocritic // because tests
	// name
	bn = simulator.Asset{
		Type:          simulator.Evc,
		Protocol:      simulator.Modbus,
		StartCharging: true,
		Phases:        simulator.One,
		CurrentType:   simulator.Ac,
		Evses:         []simulator.Evse{},
	}

	// type
	bt = simulator.Asset{
		Name:          "name",
		Protocol:      simulator.Modbus,
		StartCharging: true,
		Phases:        simulator.One,
		CurrentType:   simulator.Ac,
		Evses:         []simulator.Evse{},
	}

	// protocol
	bp = simulator.Asset{
		Name:          "name",
		Type:          simulator.Evc,
		StartCharging: true,
		Phases:        simulator.One,
		CurrentType:   simulator.Ac,
		Evses:         []simulator.Evse{},
	}

	// start charging
	bst = simulator.Asset{
		Name:        "name",
		Type:        simulator.Evc,
		Protocol:    simulator.Modbus,
		Phases:      simulator.One,
		CurrentType: simulator.Ac,
		Evses:       []simulator.Evse{},
	}

	// phases
	bph = simulator.Asset{
		Name:          "name",
		Type:          simulator.Evc,
		Protocol:      simulator.Modbus,
		StartCharging: true,
		CurrentType:   simulator.Ac,
		Evses:         []simulator.Evse{},
	}

	// current type
	bc = simulator.Asset{
		Name:          "name",
		Type:          simulator.Evc,
		Protocol:      simulator.Modbus,
		StartCharging: true,
		Phases:        simulator.One,
		Evses:         []simulator.Evse{},
	}

	// evses
	be = simulator.Asset{
		Name:          "name",
		Type:          simulator.Evc,
		Protocol:      simulator.Modbus,
		StartCharging: true,
		Phases:        simulator.One,
		CurrentType:   simulator.Ac,
	}

	return bn, bt, bp, bst, bph, bc, be
}

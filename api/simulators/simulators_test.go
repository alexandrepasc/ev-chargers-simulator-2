//nolint:errcheck,dupl //because they are tests
package simulators_test

import (
	"bytes"
	"encoding/json"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"os"
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

type method string

const (
	ep       string = "/simulators"
	modelsEp string = ep + "/models"
	mGet     method = "GET"
	mPost    method = "POST"
	mPut     method = "PUT"
	mDelete  method = "DELETE"
)

func TestGetSimulators(t *testing.T) {
	r := gin.Default()

	mID0, _ := uuid.Parse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa0")
	mID1, _ := uuid.Parse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa1")

	al := []*simulator.Asset{
		{
			Name:        "a_name",
			Type:        simulator.Evc,
			Protocol:    simulator.Ocpp201,
			Model:       mID0,
			Phases:      simulator.One,
			CurrentType: simulator.Ac,
			Evses:       []simulator.Evse{},
		},
		{
			Name:        "b_name",
			Type:        simulator.Evc,
			Protocol:    simulator.Ocpp16,
			Model:       mID1,
			Phases:      simulator.One,
			CurrentType: simulator.Ac,
			Evses:       []simulator.Evse{},
		},
	}
	ml := []*model.Struct{}

	sim, _, _ := before(t, al, ml)

	s := simulators.Simulators{
		Sim: sim,
	}

	s.Simulators(r)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest(string(mGet), ep, http.NoBody)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	a := simulators.GetSimulatorsTemp{}

	json.Unmarshal(w.Body.Bytes(), &a)

	assert.Equal(t, int64(2), a.Total)

	for i := 0; i < len(al); i++ {
		assert.Equal(t, string(al[i].Type), a.Assets.([]interface{})[i].(map[string]interface{})["type"])

		assert.Equal(t, string(al[i].Protocol), a.Assets.([]interface{})[i].(map[string]interface{})["protocol"])

		assert.Equal(t, al[i].Model.String(), a.Assets.([]interface{})[i].(map[string]interface{})["model"])
	}
}

func TestGetSimulatorsNoAssets(t *testing.T) {
	r := gin.Default()

	al := []*simulator.Asset{}
	ml := []*model.Struct{}

	sim, _, _ := before(t, al, ml)

	s := simulators.Simulators{
		Sim: sim,
	}

	s.Simulators(r)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest(string(mGet), ep, http.NoBody)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	a := simulators.GetSimulatorsTemp{}

	json.Unmarshal(w.Body.Bytes(), &a)

	assert.Equal(t, int64(0), a.Total)

	assert.Nil(t, a.Assets)
}

func TestPostSimulators(t *testing.T) {
	r := gin.Default()

	al := []*simulator.Asset{}
	ml := []*model.Struct{}

	sim, _, tmp := before(t, al, ml)

	sim.Al = al

	s := simulators.Simulators{
		Sim: sim,
	}

	body := simulator.Asset{
		Name:            "name",
		Type:            simulator.Evc,
		Protocol:        simulator.Modbus,
		StartCharging:   true,
		Phases:          simulator.One,
		CurrentType:     simulator.Ac,
		AuthorizeRemote: true,
		AuthList:        true,
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

	req, _ := http.NewRequest(string(mPost), ep, b)

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
	ml := []*model.Struct{}

	sim, _, _ := before(t, al, ml)

	sim.Al = al

	s := simulators.Simulators{
		Sim: sim,
	}

	bn, bt, bp, bst, bph, bc, be := postSimulatorsRequiredFields()

	// name
	j, _ := json.Marshal(bn)
	b := bytes.NewReader(j)

	s.Simulators(r)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest(string(mPost), ep, b)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	// type
	j, _ = json.Marshal(bt)
	b = bytes.NewReader(j)

	req, _ = http.NewRequest(string(mPost), ep, b)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	// protocol
	j, _ = json.Marshal(bp)
	b = bytes.NewReader(j)

	req, _ = http.NewRequest(string(mPost), ep, b)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	// start charging
	j, _ = json.Marshal(bst)
	b = bytes.NewReader(j)

	req, _ = http.NewRequest(string(mPost), ep, b)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	// phases
	j, _ = json.Marshal(bph)
	b = bytes.NewReader(j)

	req, _ = http.NewRequest(string(mPost), ep, b)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	// current type
	j, _ = json.Marshal(bc)
	b = bytes.NewReader(j)

	req, _ = http.NewRequest(string(mPost), ep, b)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	// evses
	j, _ = json.Marshal(be)
	b = bytes.NewReader(j)

	req, _ = http.NewRequest(string(mPost), ep, b)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPostSimulatorsBadRequestBody(t *testing.T) {
	r := gin.Default()

	al := []*simulator.Asset{}
	ml := []*model.Struct{}

	sim, _, _ := before(t, al, ml)

	sim.Al = al

	s := simulators.Simulators{
		Sim:  sim,
		Lang: translation.Translation{L: translation.EnGb},
	}

	j, _ := json.Marshal("{nothing: atall}")
	b := bytes.NewReader(j)

	s.Simulators(r)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest(string(mPost), ep, b)

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
	ml := []*model.Struct{}

	sim, _, _ := before(t, al, ml)

	sim.Al = al

	s := simulators.Simulators{
		Sim:  sim,
		Lang: translation.Translation{L: translation.EnGb},
	}

	body := simulator.Asset{
		Name:            "test1",
		Type:            simulator.Evc,
		Protocol:        simulator.Modbus,
		StartCharging:   true,
		Phases:          simulator.One,
		CurrentType:     simulator.Ac,
		AuthorizeRemote: true,
		AuthList:        true,
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

	req, _ := http.NewRequest(string(mPost), ep, b)

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
	ml := []*model.Struct{}

	sim, _, _ := before(t, al, ml)

	sim.Al = al

	s := simulators.Simulators{
		Sim:  sim,
		Lang: translation.Translation{L: translation.EnGb},
	}

	e := simulator.Asset{
		Type:            simulator.Pm,
		Protocol:        simulator.Ocpp16,
		StartCharging:   false,
		Phases:          simulator.Three,
		CurrentType:     simulator.Ac,
		AuthorizeRemote: true,
		AuthList:        true,
		Evses:           []simulator.Evse{},
	}

	j, _ := json.Marshal(e)
	b := bytes.NewReader(j)

	s.Simulators(r)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest(string(mPut), ep+"/"+al[0].SimID.String(), b)

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
	ml := []*model.Struct{}

	sim, _, _ := before(t, al, ml)

	sim.Al = al

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

	req, _ := http.NewRequest(string(mPut), ep+"/aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa0", b)

	r.ServeHTTP(w, req)

	a := errors.ErroMsg{}
	json.Unmarshal(w.Body.Bytes(), &a)

	assert.Equal(t, http.StatusNotFound, w.Code)

	assert.Equal(t, sim.L.Get(text.UpdateSimConfFileNotFound), a.Message)
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
	ml := []*model.Struct{}

	sim, _, _ := before(t, al, ml)

	sim.Al = al

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

	req, _ := http.NewRequest(string(mPut), ep+"/asd", b)

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
	ml := []*model.Struct{}

	sim, _, _ := before(t, al, ml)

	sim.Al = al

	s := simulators.Simulators{
		Sim:  sim,
		Lang: translation.Translation{L: translation.EnGb},
	}

	eid, _ := uuid.Parse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa9")
	e := simulator.Asset{
		SimID:           eid,
		Name:            "testing",
		Type:            simulator.Pm,
		Protocol:        simulator.Ocpp16,
		StartCharging:   false,
		Phases:          simulator.Three,
		CurrentType:     simulator.Ac,
		AuthorizeRemote: true,
		AuthList:        true,
		Evses:           []simulator.Evse{},
	}

	j, _ := json.Marshal(e)
	b := bytes.NewReader(j)

	s.Simulators(r)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest(string(mPut), ep+"/"+al[0].SimID.String(), b)

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
	ml := []*model.Struct{}

	sim, _, _ := before(t, al, ml)

	sim.Al = al

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

	req, _ := http.NewRequest(string(mPut), ep+"/"+al[0].SimID.String(), b)

	r.ServeHTTP(w, req)

	a := errors.ErroMsg{}
	json.Unmarshal(w.Body.Bytes(), &a)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	assert.Equal(t, sim.L.Get(text.RequestBodyDoesntMatch), a.Message)
}

func TestDeleteSimulators(t *testing.T) {
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
		{
			SimID:         uuid.New(),
			Name:          "test2",
			Type:          simulator.Evc,
			Protocol:      simulator.Ocpp201,
			StartCharging: true,
			Phases:        simulator.One,
			CurrentType:   simulator.Dc,
			Evses:         []simulator.Evse{},
		},
	}
	ml := []*model.Struct{}

	sim, _, _ := before(t, al, ml)

	sim.Al = al

	s := simulators.Simulators{
		Sim:  sim,
		Lang: translation.Translation{L: translation.EnGb},
	}

	s.Simulators(r)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest(string(mDelete), ep+"/"+al[0].SimID.String(), http.NoBody)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)

	assert.Empty(t, w.Body.String())

	req, _ = http.NewRequest(string(mGet), ep, http.NoBody)

	r.ServeHTTP(w, req)

	a := simulators.GetSimulatorsTemp{}

	json.Unmarshal(w.Body.Bytes(), &a)

	assert.Equal(t, int64(1), a.Total)

	assert.Equal(t, al[1].SimID.String(), a.Assets.([]interface{})[0].(map[string]interface{})["simId"])

	assert.Equal(t, al[1].Name, a.Assets.([]interface{})[0].(map[string]interface{})["name"])
}

func TestDeleteSimulatorsWrongId(t *testing.T) {
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
	ml := []*model.Struct{}

	sim, _, _ := before(t, al, ml)

	sim.Al = al

	s := simulators.Simulators{
		Sim:  sim,
		Lang: translation.Translation{L: translation.EnGb},
	}

	s.Simulators(r)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest(string(mDelete), ep+"/aaaaaaaa-aaaa-bbbb-aaaa-aaaaaaaaaaa0", http.NoBody)

	r.ServeHTTP(w, req)

	a := errors.ErroMsg{}
	json.Unmarshal(w.Body.Bytes(), &a)

	assert.Equal(t, http.StatusNotFound, w.Code)

	assert.Equal(t, s.Lang.Get(text.DeleteSimConfFileNotFound), a.Message)

	w = httptest.NewRecorder()

	req, _ = http.NewRequest(string(mGet), ep, http.NoBody)

	r.ServeHTTP(w, req)

	ag := simulators.GetSimulatorsTemp{}

	json.Unmarshal(w.Body.Bytes(), &ag)

	assert.Equal(t, int64(1), ag.Total)

	assert.Equal(t, al[0].SimID.String(), ag.Assets.([]interface{})[0].(map[string]interface{})["simId"])

	assert.Equal(t, al[0].Name, ag.Assets.([]interface{})[0].(map[string]interface{})["name"])
}

func TestDeleteSimulatorsInvalidId(t *testing.T) {
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
	ml := []*model.Struct{}

	sim, _, _ := before(t, al, ml)

	sim.Al = al

	s := simulators.Simulators{
		Sim:  sim,
		Lang: translation.Translation{L: translation.EnGb},
	}

	s.Simulators(r)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest(string(mDelete), ep+"/asd", http.NoBody)

	r.ServeHTTP(w, req)

	a := errors.ErroMsg{}
	json.Unmarshal(w.Body.Bytes(), &a)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	assert.Equal(t, s.Lang.Get(text.UUIDParsingError), a.Message)

	w = httptest.NewRecorder()

	req, _ = http.NewRequest(string(mGet), ep, http.NoBody)

	r.ServeHTTP(w, req)

	ag := simulators.GetSimulatorsTemp{}

	json.Unmarshal(w.Body.Bytes(), &ag)

	assert.Equal(t, int64(1), ag.Total)

	assert.Equal(t, al[0].SimID.String(), ag.Assets.([]interface{})[0].(map[string]interface{})["simId"])

	assert.Equal(t, al[0].Name, ag.Assets.([]interface{})[0].(map[string]interface{})["name"])
}

func TestGetSimModels(t *testing.T) {
	r := gin.Default()

	id1, _ := uuid.Parse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa1")
	id2, _ := uuid.Parse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa2")

	al := []*simulator.Asset{}
	ml := []*model.Struct{
		{
			ID:   id1,
			Name: "mod1",
			Type: model.Modbus,
		},
		{
			ID:   id2,
			Name: "mod2",
			Type: model.Ocpp,
		},
	}

	_, mod, _ := before(t, al, ml)

	s := simulators.Simulators{
		Mod: mod,
	}

	s.Simulators(r)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest(string(mGet), modelsEp, http.NoBody)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	a := simulators.GetModelsTemp{}

	json.Unmarshal(w.Body.Bytes(), &a)

	assert.Equal(t, int64(2), a.Total)

	for i := 0; i < len(ml); i++ {
		assert.Equal(t, ml[i].ID.String(), a.Models.([]interface{})[i].(map[string]interface{})["id"])

		assert.Equal(t, ml[i].Name, a.Models.([]interface{})[i].(map[string]interface{})["name"])

		assert.Equal(t, string(ml[i].Type), a.Models.([]interface{})[i].(map[string]interface{})["type"])
	}
}

func TestGetSimModelsNoModels(t *testing.T) {
	r := gin.Default()

	al := []*simulator.Asset{}
	ml := []*model.Struct{}

	_, mod, _ := before(t, al, ml)

	s := simulators.Simulators{
		Mod: mod,
	}

	s.Simulators(r)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest(string(mGet), modelsEp, http.NoBody)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	a := simulators.GetModelsTemp{}

	json.Unmarshal(w.Body.Bytes(), &a)

	assert.Equal(t, int64(0), a.Total)

	assert.Nil(t, a.Models)
}

func TestPostModels(t *testing.T) {
	r := gin.Default()

	_, mod, tmp := before(t, []*simulator.Asset{}, []*model.Struct{})

	s := simulators.Simulators{
		Mod:  mod,
		Lang: mod.L,
	}

	e := model.Struct{
		Name: "name",
		Type: model.Modbus,
		Ocpp: model.OcppModel{
			SerialNumb: "serial",
		},
	}

	j, _ := json.Marshal(e)
	b := bytes.NewReader(j)

	s.Simulators(r)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest(string(mPost), modelsEp, b)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	a := model.Struct{}
	json.Unmarshal(w.Body.Bytes(), &a)

	_, err := uuid.Parse(a.ID.String())
	assert.Nil(t, err)

	assert.Equal(t, e.Name, a.Name)

	assert.Equal(t, e.Type, a.Type)

	assert.Equal(t, e.Ocpp.SerialNumb, a.Ocpp.SerialNumb)

	assert.FileExists(t, tmp+"/simConf/models/"+e.Name+".json")
}

func TestPostModelsRequiredFields(t *testing.T) {
	r := gin.Default()

	_, mod, _ := before(t, []*simulator.Asset{}, []*model.Struct{})

	s := simulators.Simulators{
		Mod:  mod,
		Lang: mod.L,
	}

	bn, bt := postModelsRequiredFields()

	s.Simulators(r)

	w := httptest.NewRecorder()

	// name
	j, _ := json.Marshal(bn)
	b := bytes.NewReader(j)

	req, _ := http.NewRequest(string(mPost), modelsEp, b)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	// type
	j, _ = json.Marshal(bt)
	b = bytes.NewReader(j)

	req, _ = http.NewRequest(string(mPost), modelsEp, b)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPostModelsBadRequestBody(t *testing.T) {
	r := gin.Default()

	_, mod, _ := before(t, []*simulator.Asset{}, []*model.Struct{})

	s := simulators.Simulators{
		Mod:  mod,
		Lang: mod.L,
	}

	j, _ := json.Marshal("{nothing: atall}")
	b := bytes.NewReader(j)

	s.Simulators(r)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest(string(mPost), modelsEp, b)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	a := errors.ErroMsg{}
	json.Unmarshal(w.Body.Bytes(), &a)

	assert.Equal(t, translation.Translation{L: translation.EnGb}.Get(text.RequestBodyDoesntMatch), a.Message)
}

func TestNotAblePostModelsSameName(t *testing.T) {
	r := gin.Default()

	id, _ := uuid.Parse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa1")
	ml := []*model.Struct{
		{
			ID:   id,
			Name: "name",
			Type: model.Modbus,
			Ocpp: model.OcppModel{
				SerialNumb: "serial",
			},
		},
	}

	sim, mod, _ := before(t, []*simulator.Asset{}, ml)

	s := simulators.Simulators{
		Sim:  sim,
		Mod:  mod,
		Lang: translation.Translation{L: translation.EnGb},
	}

	e := model.Struct{
		Name: "name",
		Type: model.Modbus,
		Ocpp: model.OcppModel{
			SerialNumb: "serial",
		},
	}

	j, _ := json.Marshal(e)
	b := bytes.NewReader(j)

	s.Simulators(r)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest(string(mPost), modelsEp, b)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)

	a := errors.ErroMsg{}
	json.Unmarshal(w.Body.Bytes(), &a)

	assert.Equal(t, translation.Translation{L: translation.EnGb}.Get(text.CreateModelConfFileNameExists), a.Message)
}

func TestPutModels(t *testing.T) {
	r := gin.Default()

	id, _ := uuid.Parse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa0")
	ol := []*model.Struct{
		{
			ID:   id,
			Name: "name",
			Type: model.Modbus,
			Ocpp: model.OcppModel{
				SerialNumb: "poi",
			},
		},
	}

	_, mod, _ := before(t, []*simulator.Asset{}, ol)

	s := simulators.Simulators{
		Mod:  mod,
		Lang: mod.L,
	}

	e := model.Struct{
		Type: model.Ocpp,
		Ocpp: model.OcppModel{
			SerialNumb: "asd",
		},
	}

	j, _ := json.Marshal(e)
	b := bytes.NewReader(j)

	s.Simulators(r)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest(string(mPut), modelsEp+"/aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa0", b)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	a := model.Struct{}
	json.Unmarshal(w.Body.Bytes(), &a)

	assert.Equal(t, e.Type, a.Type)

	assert.Equal(t, e.Ocpp.SerialNumb, a.Ocpp.SerialNumb)
}

func TestPutModelsWrongId(t *testing.T) {
	r := gin.Default()

	id, _ := uuid.Parse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa0")
	ol := []*model.Struct{
		{
			ID:   id,
			Name: "name",
			Type: model.Modbus,
			Ocpp: model.OcppModel{
				SerialNumb: "poi",
			},
		},
	}

	_, mod, _ := before(t, []*simulator.Asset{}, ol)

	s := simulators.Simulators{
		Mod:  mod,
		Lang: mod.L,
	}

	e := model.Struct{
		Type: model.Ocpp,
		Ocpp: model.OcppModel{
			SerialNumb: "asd",
		},
	}

	j, _ := json.Marshal(e)
	b := bytes.NewReader(j)

	s.Simulators(r)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest(string(mPut), modelsEp+"/aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa", b)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	a := errors.ErroMsg{}
	json.Unmarshal(w.Body.Bytes(), &a)

	assert.Equal(t, s.Lang.Get(text.UpdateModelConfFileNotFound), a.Message)
}

func TestPutModelsInvalidId(t *testing.T) {
	r := gin.Default()

	id, _ := uuid.Parse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa0")
	ol := []*model.Struct{
		{
			ID:   id,
			Name: "name",
			Type: model.Modbus,
			Ocpp: model.OcppModel{
				SerialNumb: "poi",
			},
		},
	}

	_, mod, _ := before(t, []*simulator.Asset{}, ol)

	s := simulators.Simulators{
		Mod:  mod,
		Lang: mod.L,
	}

	e := model.Struct{
		Type: model.Ocpp,
		Ocpp: model.OcppModel{
			SerialNumb: "asd",
		},
	}

	j, _ := json.Marshal(e)
	b := bytes.NewReader(j)

	s.Simulators(r)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest(string(mPut), modelsEp+"/asd", b)

	r.ServeHTTP(w, req)

	a := errors.ErroMsg{}
	json.Unmarshal(w.Body.Bytes(), &a)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	assert.Equal(t, s.Lang.Get(text.UUIDParsingError), a.Message)
}

func TestPutModelsChangeName(t *testing.T) {
	r := gin.Default()

	id, _ := uuid.Parse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa0")
	ol := []*model.Struct{
		{
			ID:   id,
			Name: "name",
			Type: model.Modbus,
			Ocpp: model.OcppModel{
				SerialNumb: "poi",
			},
		},
	}

	_, mod, _ := before(t, []*simulator.Asset{}, ol)

	s := simulators.Simulators{
		Mod:  mod,
		Lang: mod.L,
	}

	e := model.Struct{
		Name: "change",
		Type: model.Ocpp,
		Ocpp: model.OcppModel{
			SerialNumb: "asd",
		},
	}

	j, _ := json.Marshal(e)
	b := bytes.NewReader(j)

	s.Simulators(r)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest(string(mPut), modelsEp+"/aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa0", b)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	a := model.Struct{}
	json.Unmarshal(w.Body.Bytes(), &a)

	assert.NotEqual(t, e.Name, a.Name)

	assert.Equal(t, e.Type, a.Type)

	assert.Equal(t, e.Ocpp.SerialNumb, a.Ocpp.SerialNumb)
}

func TestPutModelsInvalidBody(t *testing.T) {
	r := gin.Default()

	id, _ := uuid.Parse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa0")
	ol := []*model.Struct{
		{
			ID:   id,
			Name: "name",
			Type: model.Modbus,
			Ocpp: model.OcppModel{
				SerialNumb: "poi",
			},
		},
	}

	_, mod, _ := before(t, []*simulator.Asset{}, ol)

	s := simulators.Simulators{
		Mod:  mod,
		Lang: mod.L,
	}

	e := model.Struct{
		Ocpp: model.OcppModel{
			SerialNumb: "asd",
		},
	}

	j, _ := json.Marshal(e)
	b := bytes.NewReader(j)

	s.Simulators(r)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest(string(mPut), modelsEp+"/aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa0", b)

	r.ServeHTTP(w, req)

	a := errors.ErroMsg{}
	json.Unmarshal(w.Body.Bytes(), &a)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	assert.Equal(t, s.Lang.Get(text.RequestBodyDoesntMatch), a.Message)
}

func TestDeleteModels(t *testing.T) {
	r := gin.Default()

	id, _ := uuid.Parse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa0")
	id2, _ := uuid.Parse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa1")
	ol := []*model.Struct{
		{
			ID:   id,
			Name: "name",
			Type: model.Modbus,
			Ocpp: model.OcppModel{
				SerialNumb: "poi",
			},
		},
		{
			ID:   id2,
			Name: "name2",
			Type: model.Modbus,
			Ocpp: model.OcppModel{
				SerialNumb: "poi2",
			},
		},
	}

	_, mod, _ := before(t, []*simulator.Asset{}, ol)

	s := simulators.Simulators{
		Mod:  mod,
		Lang: mod.L,
	}

	s.Simulators(r)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest(string(mDelete), modelsEp+"/aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa0", http.NoBody)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)

	assert.Empty(t, w.Body.String())

	req, _ = http.NewRequest(string(mGet), modelsEp, http.NoBody)

	r.ServeHTTP(w, req)

	a := simulators.GetModelsTemp{}

	json.Unmarshal(w.Body.Bytes(), &a)

	assert.Equal(t, int64(1), a.Total)

	assert.Equal(t, ol[1].ID.String(), a.Models.([]interface{})[0].(map[string]interface{})["id"])

	assert.Equal(t, ol[1].Name, a.Models.([]interface{})[0].(map[string]interface{})["name"])
}

func TestDeleteModelsWrongId(t *testing.T) {
	r := gin.Default()

	id, _ := uuid.Parse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa0")
	ol := []*model.Struct{
		{
			ID:   id,
			Name: "name",
			Type: model.Modbus,
			Ocpp: model.OcppModel{
				SerialNumb: "poi",
			},
		},
	}

	_, mod, _ := before(t, []*simulator.Asset{}, ol)

	s := simulators.Simulators{
		Mod:  mod,
		Lang: mod.L,
	}

	s.Simulators(r)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest(string(mDelete), modelsEp+"/aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaab", http.NoBody)

	r.ServeHTTP(w, req)

	a := errors.ErroMsg{}
	json.Unmarshal(w.Body.Bytes(), &a)

	assert.Equal(t, http.StatusNotFound, w.Code)

	assert.Equal(t, s.Lang.Get(text.DeleteModelConfFileNotFound), a.Message)
}

func TestDeleteModelsInvalidId(t *testing.T) {
	r := gin.Default()

	id, _ := uuid.Parse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa0")
	ol := []*model.Struct{
		{
			ID:   id,
			Name: "name",
			Type: model.Modbus,
			Ocpp: model.OcppModel{
				SerialNumb: "poi",
			},
		},
	}

	_, mod, _ := before(t, []*simulator.Asset{}, ol)

	s := simulators.Simulators{
		Mod:  mod,
		Lang: mod.L,
	}

	s.Simulators(r)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest(string(mDelete), modelsEp+"/aaa", http.NoBody)

	r.ServeHTTP(w, req)

	a := errors.ErroMsg{}
	json.Unmarshal(w.Body.Bytes(), &a)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	assert.Equal(t, s.Lang.Get(text.UUIDParsingError), a.Message)
}

// TODO: Add tests to the run and stop endpoints

func before(t *testing.T, al []*simulator.Asset, ml []*model.Struct) (s simulator.Simulator, m model.Model, tmp string) {
	t.Helper()

	tmp = t.TempDir()

	os.Mkdir(tmp+"/simConf", fs.FileMode(common.FolderPermissions))

	os.Mkdir(tmp+"/simConf"+common.DefMCFolder, fs.FileMode(common.FolderPermissions))

	generateConfFiles(tmp, al, ml)

	l := translation.Translation{
		L: translation.Translation{}.GetKey("en-GB"),
	}

	s = simulator.Simulator{
		Scp: tmp + "/simConf",
		L:   l,
	}

	m = model.Model{
		Scp: tmp + "/simConf",
		L:   l,
	}

	return s, m, tmp
}

func generateConfFiles(tmp string, al []*simulator.Asset, ml []*model.Struct) {
	for i := 0; i < len(al); i++ {
		p := tmp + "/simConf" + "/" + al[i].Name + ".json"

		f, _ := os.Create(p)
		f.Close()

		b, _ := json.MarshalIndent(al[i], "", " ")

		os.WriteFile(p, b, fs.FileMode(common.FilePermissions))
	}

	for i := 0; i < len(ml); i++ {
		p := tmp + "/simConf" + common.DefMCFolder + "/" + ml[i].Name + ".json"

		f, _ := os.Create(p)
		f.Close()

		b, _ := json.MarshalIndent(ml[i], "", " ")

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

func postModelsRequiredFields() (bn, bt model.Struct) {
	// name
	bn = model.Struct{
		Type: model.Modbus,
		Ocpp: model.OcppModel{
			SerialNumb: "serial",
		},
	}

	// type
	bt = model.Struct{
		Name: "name",
		Ocpp: model.OcppModel{
			SerialNumb: "serial",
		},
	}

	return bn, bt
}

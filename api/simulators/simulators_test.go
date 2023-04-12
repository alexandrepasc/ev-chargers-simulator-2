//nolint:errcheck //because they are tests
package simulators_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alexandrepasc/ev-chargers-simulator-2/api/simulators"
	"github.com/alexandrepasc/ev-chargers-simulator-2/simulator"
	"github.com/alexandrepasc/ev-chargers-simulator-2/simulator/model"
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
			Type:     simulator.Evc,
			Protocol: simulator.Ocpp201,
			Model:    "a",
		},
		{
			Type:     simulator.Evc,
			Protocol: simulator.Ocpp16,
			Model:    "b",
		},
	}

	s := simulators.Simulators{
		Al: al,
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

	s := simulators.Simulators{
		Al: al,
	}

	s.Simulators(r)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", ep, http.NoBody)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	a := simulators.GetSimulatorsTemp{}

	json.Unmarshal(w.Body.Bytes(), &a)

	assert.Equal(t, int64(0), a.Total)

	assert.Equal(t, 0, len(a.Assets.([]interface{})))
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

	s := simulators.Simulators{
		Al:  al,
		Oml: ml,
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

	s := simulators.Simulators{
		Al:  al,
		Oml: ml,
	}

	s.Simulators(r)

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", modelsEp, http.NoBody)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	a := simulators.GetModelsTemp{}

	json.Unmarshal(w.Body.Bytes(), &a)

	assert.Equal(t, int64(0), a.Total)

	assert.Equal(t, 0, len(a.Models.([]interface{})))
}

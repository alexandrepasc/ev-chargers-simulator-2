//nolint:errcheck // because this is a test file
package simulator_test

import (
	"encoding/json"
	"io/fs"
	"os"
	"strconv"
	"testing"

	"github.com/alexandrepasc/ev-chargers-simulator-2/common"
	"github.com/alexandrepasc/ev-chargers-simulator-2/simulator"
	"github.com/alexandrepasc/ev-chargers-simulator-2/simulator/model"
	"github.com/alexandrepasc/ev-chargers-simulator-2/translation"
	"github.com/stretchr/testify/assert"
)

const defSCFolder string = "/simConf"

func TestGetSimulatorData(t *testing.T) {
	tmp := t.TempDir()

	createFolders(tmp)

	s := before(tmp)

	generateAssetConfFiles(1, tmp)

	generateModelConfFiles(1, tmp)

	al, ml := s.GetSimsConfs()

	assert.Equal(t, 1, len(al))

	assert.Equal(t, 1, len(ml))
}

func TestGetSimulatorDataNoFiles(t *testing.T) {
	tmp := t.TempDir()

	createFolders(tmp)

	s := before(tmp)

	al, ml := s.GetSimsConfs()

	assert.Equal(t, 0, len(al))

	assert.Equal(t, 0, len(ml))
}

func TestSimulatorModelWrongType(t *testing.T) {
	tmp := t.TempDir()

	createFolders(tmp)

	s := before(tmp)

	generateAssetConfFiles(1, tmp)

	generateInvalidModelConfFile(tmp)

	al, ml := s.GetSimsConfs()

	assert.Equal(t, 1, len(al))

	assert.Equal(t, 0, len(ml))
}

func TestCreateSimulator(t *testing.T) {
	tmp := t.TempDir()

	createFolders(tmp)

	s := before(tmp)

	sim := simulator.Asset{
		Name:          "sim1",
		Type:          simulator.Evc,
		Protocol:      simulator.Ocpp201,
		Model:         "model",
		StartCharging: true,
		Phases:        simulator.One,
		CurrentType:   simulator.Ac,
		Evses: []simulator.Evse{
			{
				ID: 1,
				Connectors: []simulator.Connector{
					{
						ID: 1,
						Data: []simulator.Data{
							{
								Duration:      20,
								ChargingState: 2,
								ErrorCode:     0,
								PowerFactor:   900,
							},
						},
					},
				},
			},
		},
	}

	s.CreateSimConf(&sim)

	nf := tmp + defSCFolder + "/" + sim.Name + ".json"

	assert.FileExists(t, nf)

	a := readFile(nf)

	assert.NotEmpty(t, a.SimID)

	assert.Equal(t, sim.Name, a.Name)

	assert.Equal(t, sim.Type, a.Type)

	assert.Equal(t, sim, a)
}

func TestCreateEmptySimulator(t *testing.T) {
	tmp := t.TempDir()

	createFolders(tmp)

	s := before(tmp)

	sim := simulator.Asset{}

	s.CreateSimConf(&sim)

	nf := tmp + defSCFolder + "/" + sim.Name + ".json"

	assert.NoFileExists(t, nf)
}

func TestCreateSimulatorRequiredFields(t *testing.T) {
	tmp := t.TempDir()

	createFolders(tmp)

	s := before(tmp)

	// name
	sim := simulator.Asset{
		Type:        simulator.Evc,
		Protocol:    simulator.Modbus,
		Model:       "asd",
		Phases:      simulator.One,
		CurrentType: simulator.Ac,
		Evses:       []simulator.Evse{},
	}

	ok, _, _ := s.CreateSimConf(&sim)

	nf := tmp + defSCFolder + "/" + sim.Name + ".json"

	assert.NoFileExists(t, nf)

	assert.False(t, ok)

	// type
	sim = simulator.Asset{
		Name:        "name",
		Protocol:    simulator.Modbus,
		Model:       "asd",
		Phases:      simulator.One,
		CurrentType: simulator.Ac,
		Evses:       []simulator.Evse{},
	}

	ok, _, _ = s.CreateSimConf(&sim)

	nf = tmp + defSCFolder + "/" + sim.Name + ".json"

	assert.NoFileExists(t, nf)

	assert.False(t, ok)

	// protocol
	sim = simulator.Asset{
		Name:        "name",
		Type:        simulator.Evc,
		Model:       "asd",
		Phases:      simulator.One,
		CurrentType: simulator.Ac,
		Evses:       []simulator.Evse{},
	}

	ok, _, _ = s.CreateSimConf(&sim)

	nf = tmp + defSCFolder + "/" + sim.Name + ".json"

	assert.NoFileExists(t, nf)

	assert.False(t, ok)

	// phases
	sim = simulator.Asset{
		Name:        "name",
		Type:        simulator.Evc,
		Protocol:    simulator.Modbus,
		Model:       "asd",
		CurrentType: simulator.Ac,
		Evses:       []simulator.Evse{},
	}

	ok, _, _ = s.CreateSimConf(&sim)

	nf = tmp + defSCFolder + "/" + sim.Name + ".json"

	assert.NoFileExists(t, nf)

	assert.False(t, ok)

	// current type
	sim = simulator.Asset{
		Name:     "name",
		Type:     simulator.Evc,
		Protocol: simulator.Modbus,
		Model:    "asd",
		Phases:   simulator.One,
		Evses:    []simulator.Evse{},
	}

	ok, _, _ = s.CreateSimConf(&sim)

	nf = tmp + defSCFolder + "/" + sim.Name + ".json"

	assert.NoFileExists(t, nf)

	assert.False(t, ok)

	// evses
	sim = simulator.Asset{
		Name:        "name",
		Type:        simulator.Evc,
		Protocol:    simulator.Modbus,
		Model:       "asd",
		Phases:      simulator.One,
		CurrentType: simulator.Ac,
	}

	ok, _, _ = s.CreateSimConf(&sim)

	nf = tmp + defSCFolder + "/" + sim.Name + ".json"

	assert.NoFileExists(t, nf)

	assert.False(t, ok)
}

func before(tmp string) simulator.Simulator {
	t := translation.Translation{
		L: translation.Translation{}.GetKey("en-GB"),
	}

	s := simulator.Simulator{
		Scp: tmp + defSCFolder,
		L:   t,
	}

	return s
}

func createFolders(tmp string) {
	os.Mkdir(tmp+defSCFolder, fs.FileMode(common.FolderPermissions))

	os.Mkdir(tmp+defSCFolder+common.DefMCFolder, fs.FileMode(common.FolderPermissions))
}

func generateAssetConfFiles(n int64, tmp string) {
	for i := int64(0); i < n; i++ {
		a := simulator.Asset{
			Type:     simulator.Evc,
			Protocol: simulator.Ocpp201,
			Phases:   1,
		}

		p := tmp + defSCFolder + "/sim" + strconv.FormatInt(i, 10) + ".json"

		f, _ := os.Create(p)
		f.Close()

		b, _ := json.MarshalIndent(a, "", " ")

		os.WriteFile(p, b, fs.FileMode(common.FilePermissions))
	}
}

func generateModelConfFiles(n int64, tmp string) {
	for i := int64(0); i < n; i++ {
		m := model.OcppModel{
			SerialNumb: "123456789",
			Model:      "model",
		}

		p := tmp + defSCFolder + common.DefMCFolder + "/model" + strconv.FormatInt(i, 10) + "_ocpp.json"

		f, _ := os.Create(p)
		f.Close()

		b, _ := json.MarshalIndent(m, "", " ")

		os.WriteFile(p, b, fs.FileMode(common.FilePermissions))
	}
}

func generateInvalidModelConfFile(tmp string) {
	m := model.OcppModel{
		SerialNumb: "123456789",
		Model:      "model",
	}

	p := tmp + defSCFolder + common.DefMCFolder + "/model_invalid.json"

	f, _ := os.Create(p)
	f.Close()

	b, _ := json.MarshalIndent(m, "", " ")

	os.WriteFile(p, b, fs.FileMode(common.FilePermissions))
}

func readFile(p string) (a simulator.Asset) {
	b, _ := os.ReadFile(p)

	json.Unmarshal(b, &a)

	return a
}

//nolint:errcheck // because this is a test file
package simulator_test

import (
	"encoding/json"
	"io/fs"
	"net/http"
	"os"
	"strconv"
	"testing"

	"github.com/alexandrepasc/ev-chargers-simulator-2/common"
	"github.com/alexandrepasc/ev-chargers-simulator-2/simulator"
	"github.com/alexandrepasc/ev-chargers-simulator-2/translation"
	"github.com/alexandrepasc/ev-chargers-simulator-2/translation/text"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

const defSCFolder string = "/simConf"

func TestGetSimulatorData(t *testing.T) {
	tmp := t.TempDir()

	createFolders(tmp)

	s := before(tmp)

	generateAssetConfFiles(1, tmp)

	al := s.GetSimulators()

	assert.Equal(t, 1, len(al))
}

func TestGetSimulatorDataNoFiles(t *testing.T) {
	tmp := t.TempDir()

	createFolders(tmp)

	s := before(tmp)

	al := s.GetSimulators()

	assert.Equal(t, 0, len(al))
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

func TestCanNotCreateSimulatorSameName(t *testing.T) {
	tmp := t.TempDir()

	createFolders(tmp)

	s := before(tmp)

	generateAssetConfFiles(1, tmp)

	sim := simulator.Asset{
		Name:          "sim0",
		Type:          simulator.Pm,
		Protocol:      simulator.Ocpp16,
		Model:         "model",
		StartCharging: true,
		Phases:        simulator.Three,
		CurrentType:   simulator.Dc,
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

	ok, _, _ := s.CreateSimConf(&sim)

	assert.False(t, ok)

	nf := tmp + defSCFolder + "/" + sim.Name + ".json"

	assert.FileExists(t, nf)

	a := readFile(nf)

	eid, _ := uuid.Parse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa0")

	assert.Equal(t, eid, a.SimID)

	assert.NotEqual(t, sim.Type, a.Type)

	assert.NotEqual(t, sim.Protocol, a.Protocol)

	assert.NotEqual(t, sim.Phases, a.Phases)

	assert.NotEqual(t, sim, a)
}

func TestUpdateSimulator(t *testing.T) {
	tmp := t.TempDir()

	createFolders(tmp)

	s := before(tmp)

	generateAssetConfFiles(1, tmp)

	id, _ := uuid.Parse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa0")

	a := simulator.Asset{
		Type:        simulator.Pm,
		Protocol:    simulator.Ocpp16,
		Phases:      simulator.Three,
		CurrentType: simulator.Dc,
		Evses:       []simulator.Evse{},
	}

	ab, _, c, aa := s.UpdateSimConf(id, &a)

	assert.True(t, ab)

	assert.Equal(t, http.StatusOK, c)

	assert.Equal(t, id, aa.SimID)

	assert.Equal(t, "sim0", aa.Name)

	assert.Equal(t, simulator.Pm, aa.Type)

	assert.Equal(t, simulator.Ocpp16, aa.Protocol)

	assert.Equal(t, simulator.Three, aa.Phases)

	p := tmp + defSCFolder + "/sim0.json"

	af := readFile(p)

	assert.Equal(t, id, af.SimID)

	assert.Equal(t, "sim0", af.Name)

	assert.Equal(t, simulator.Pm, af.Type)

	assert.Equal(t, simulator.Ocpp16, af.Protocol)

	assert.Equal(t, simulator.Three, af.Phases)
}

func TestUpdateSimulatorWithWrongId(t *testing.T) {
	tmp := t.TempDir()

	createFolders(tmp)

	s := before(tmp)

	generateAssetConfFiles(1, tmp)

	id, _ := uuid.Parse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa9")

	a := simulator.Asset{
		Type:     simulator.Pm,
		Protocol: simulator.Ocpp16,
		Phases:   simulator.Three,
	}

	ab, am, c, aa := s.UpdateSimConf(id, &a)

	assert.False(t, ab)

	assert.Equal(t, s.L.Get(text.UpdateSimConfFileNotFound), am)

	assert.Equal(t, http.StatusNotFound, c)

	assert.Empty(t, aa)

	p := tmp + defSCFolder + "/sim0.json"

	af := readFile(p)

	eid, _ := uuid.Parse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa0")

	assert.Equal(t, eid, af.SimID)

	assert.Equal(t, "sim0", af.Name)

	assert.Equal(t, simulator.Evc, af.Type)

	assert.Equal(t, simulator.Ocpp201, af.Protocol)

	assert.Equal(t, simulator.One, af.Phases)
}

func TestUpdateSimulatorEmpty(t *testing.T) {
	tmp := t.TempDir()

	createFolders(tmp)

	s := before(tmp)

	generateAssetConfFiles(1, tmp)

	id, _ := uuid.Parse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa0")

	a := simulator.Asset{}

	ab, am, c, _ := s.UpdateSimConf(id, &a)

	assert.False(t, ab)

	assert.Equal(t, http.StatusBadRequest, c)

	assert.Equal(t, s.L.Get(text.RequestBodyDoesntMatch), am)

	p := tmp + defSCFolder + "/sim0.json"

	af := readFile(p)

	eid, _ := uuid.Parse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa0")

	assert.Equal(t, eid, af.SimID)

	assert.Equal(t, "sim0", af.Name)

	assert.Equal(t, simulator.Evc, af.Type)

	assert.Equal(t, simulator.Ocpp201, af.Protocol)

	assert.Equal(t, simulator.One, af.Phases)
}

func TestUpdateSimulatorChangeName(t *testing.T) {
	tmp := t.TempDir()

	createFolders(tmp)

	s := before(tmp)

	generateAssetConfFiles(1, tmp)

	id, _ := uuid.Parse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa0")

	a := simulator.Asset{
		Name:        "asdasd",
		Type:        simulator.Pm,
		Protocol:    simulator.Ocpp16,
		Phases:      simulator.Three,
		CurrentType: simulator.Dc,
		Evses:       []simulator.Evse{},
	}

	ab, am, c, aa := s.UpdateSimConf(id, &a)

	assert.True(t, ab)

	assert.Equal(t, "", am)

	assert.Equal(t, http.StatusOK, c)

	assert.Equal(t, id, aa.SimID)

	assert.Equal(t, "sim0", aa.Name)

	assert.Equal(t, simulator.Pm, aa.Type)

	assert.Equal(t, simulator.Ocpp16, aa.Protocol)

	assert.Equal(t, simulator.Three, aa.Phases)

	p := tmp + defSCFolder + "/sim0.json"

	af := readFile(p)

	eid, _ := uuid.Parse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa0")

	assert.Equal(t, eid, af.SimID)

	assert.Equal(t, "sim0", af.Name)

	assert.Equal(t, simulator.Pm, af.Type)

	assert.Equal(t, simulator.Ocpp16, af.Protocol)

	assert.Equal(t, simulator.Three, af.Phases)
}

func TestDeleteSimulator(t *testing.T) {
	tmp := t.TempDir()

	createFolders(tmp)

	s := before(tmp)

	generateAssetConfFiles(2, tmp)

	id, _ := uuid.Parse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa0")

	ab, am, ac := s.DeleteSimConf(id)

	p0 := tmp + defSCFolder + "/sim0.json"
	p1 := tmp + defSCFolder + "/sim1.json"

	assert.True(t, ab)

	assert.Empty(t, am)

	assert.Equal(t, http.StatusNoContent, ac)

	assert.NoFileExists(t, p0)

	assert.FileExists(t, p1)
}

func TestDeleteSimulatorWrongId(t *testing.T) {
	tmp := t.TempDir()

	createFolders(tmp)

	s := before(tmp)

	generateAssetConfFiles(2, tmp)

	id, _ := uuid.Parse("aaaaaaaa-aaaa-aaaa-bbbb-aaaaaaaaaaa0")

	ab, am, ac := s.DeleteSimConf(id)

	p0 := tmp + defSCFolder + "/sim0.json"
	p1 := tmp + defSCFolder + "/sim1.json"

	assert.False(t, ab)

	assert.Equal(t, translation.Translation{L: translation.EnGb}.Get(text.DeleteSimConfFileNotFound), am)

	assert.Equal(t, http.StatusNotFound, ac)

	assert.FileExists(t, p0)

	assert.FileExists(t, p1)
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

func generateAssetConfFiles(n int64, tmp string) { //nolint:unparam,nolintlint // because is a test
	for i := int64(0); i < n; i++ {
		id, _ := uuid.Parse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa" + strconv.FormatInt(i, 10))

		a := simulator.Asset{
			SimID:    id,
			Name:     "sim" + strconv.FormatInt(i, 10),
			Type:     simulator.Evc,
			Protocol: simulator.Ocpp201,
			Phases:   simulator.One,
		}

		p := tmp + defSCFolder + "/sim" + strconv.FormatInt(i, 10) + ".json"

		f, _ := os.Create(p)
		f.Close()

		b, _ := json.MarshalIndent(a, "", " ")

		os.WriteFile(p, b, fs.FileMode(common.FilePermissions))
	}
}

func readFile(p string) (a simulator.Asset) {
	b, _ := os.ReadFile(p)

	json.Unmarshal(b, &a)

	return a
}

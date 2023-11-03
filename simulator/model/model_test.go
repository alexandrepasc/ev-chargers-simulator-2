//nolint:errcheck // because this is a test file
package model_test

import (
	"encoding/json"
	"io/fs"
	"net/http"
	"os"
	"strconv"
	"testing"

	"github.com/alexandrepasc/ev-chargers-simulator-2/common"
	"github.com/alexandrepasc/ev-chargers-simulator-2/simulator/model"
	"github.com/alexandrepasc/ev-chargers-simulator-2/translation"
	"github.com/alexandrepasc/ev-chargers-simulator-2/translation/text"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

const (
	defSCFolder string = "/simConf"
	defFolder   string = "/simConf" + common.DefMCFolder
	mod0        string = "/model0.json"
)

func TestGetModelData(t *testing.T) {
	tmp := t.TempDir()

	createFolders(tmp)

	m := before(tmp)

	generateModelDefConfFiles(1, tmp)

	ml := m.GetModels()

	assert.Equal(t, 1, len(ml))
}

func TestGetModelDataNoFiles(t *testing.T) {
	tmp := t.TempDir()

	createFolders(tmp)

	m := before(tmp)

	ml := m.GetModels()

	assert.Equal(t, 0, len(ml))
}

func TestGetModelInvalidData(t *testing.T) {
	tmp := t.TempDir()

	createFolders(tmp)

	m := before(tmp)

	generateInvalidModelConfFile(tmp)

	ml := m.GetModels()

	assert.Equal(t, 0, len(ml))
}

func TestCreateModelReturn(t *testing.T) {
	tmp := t.TempDir()

	createFolders(tmp)

	m := before(tmp)

	e := model.Struct{
		Name: "model",
		Type: model.Ocpp,
		Ocpp: model.OcppModel{
			SerialNumb: "123-asd-zxc",
		},
	}

	ab, as, ac, a := m.CreateModel(&e)

	assert.True(t, ab)

	assert.Empty(t, as)

	assert.Equal(t, http.StatusCreated, ac)

	assert.NotEmpty(t, a.ID)

	assert.Equal(t, e.Name, a.Name)

	assert.Equal(t, e.Type, a.Type)

	assert.Equal(t, e.Ocpp.SerialNumb, a.Ocpp.SerialNumb)
}

func TestCreateModelFile(t *testing.T) {
	tmp := t.TempDir()

	createFolders(tmp)

	m := before(tmp)

	e := model.Struct{
		Name: "model",
		Type: model.Ocpp,
		Ocpp: model.OcppModel{
			SerialNumb: "123-asd-zxc",
		},
	}

	m.CreateModel(&e)

	nf := tmp + defFolder + "/model.json"

	assert.FileExists(t, nf)

	a := readFile(nf)

	assert.NotEmpty(t, a.ID)

	assert.Equal(t, e.Name, a.Name)

	assert.Equal(t, e.Type, a.Type)

	assert.Equal(t, e.Ocpp.SerialNumb, a.Ocpp.SerialNumb)
}

func TestCreateEmptyModel(t *testing.T) {
	tmp := t.TempDir()

	createFolders(tmp)

	m := before(tmp)

	e := model.Struct{}

	ab, as, ac, _ := m.CreateModel(&e)

	nf := tmp + defFolder + "/model.json"

	assert.NoFileExists(t, nf)

	assert.False(t, ab)

	assert.NotEmpty(t, as)

	assert.Equal(t, http.StatusBadRequest, ac)
}

func TestCreateModelRequiredFields(t *testing.T) {
	tmp := t.TempDir()

	createFolders(tmp)

	m := before(tmp)

	// name
	e := model.Struct{
		Type: model.Ocpp,
		Ocpp: model.OcppModel{
			SerialNumb: "123-asd-zxc",
		},
	}

	ab, as, ac, _ := m.CreateModel(&e)

	assert.False(t, ab)

	assert.Equal(t, m.L.Get(text.RequestBodyDoesntMatch), as)

	assert.Equal(t, http.StatusBadRequest, ac)

	// type
	e = model.Struct{
		Name: "model",
		Ocpp: model.OcppModel{
			SerialNumb: "123-asd-zxc",
		},
	}

	ab, as, ac, _ = m.CreateModel(&e)

	assert.False(t, ab)

	assert.Equal(t, m.L.Get(text.RequestBodyDoesntMatch), as)

	assert.Equal(t, http.StatusBadRequest, ac)
}

func TestCanNotCreateModelSameName(t *testing.T) {
	tmp := t.TempDir()

	createFolders(tmp)

	m := before(tmp)

	generateModelDefConfFiles(1, tmp)

	e := model.Struct{
		Name: "model0",
		Type: model.Ocpp,
		Ocpp: model.OcppModel{
			SerialNumb: "123-asd-zxc",
		},
	}

	ab, as, ac, a := m.CreateModel(&e)

	assert.False(t, ab)

	// TODO: need to create translation to this error
	assert.Equal(t, m.L.Get(text.CreateModelConfFileNameExists), as)

	assert.Equal(t, http.StatusConflict, ac)

	assert.Empty(t, a)
}

func TestUpdateModel(t *testing.T) {
	tmp := t.TempDir()

	createFolders(tmp)

	m := before(tmp)

	generateModelDefConfFiles(1, tmp)

	id, _ := uuid.Parse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa0")

	e := model.Struct{
		Name: "model0",
		Type: model.Modbus,
		Ocpp: model.OcppModel{
			SerialNumb: "123",
		},
	}

	ab, as, ai, a := m.UpdateModel(id, &e)

	assert.True(t, ab)

	assert.Empty(t, as)

	assert.Equal(t, http.StatusOK, ai)

	assert.Equal(t, e.ID, a.ID)

	assert.Equal(t, e.Name, a.Name)

	assert.Equal(t, e.Type, e.Type)

	assert.Equal(t, e.Ocpp.SerialNumb, a.Ocpp.SerialNumb)
}

func TestUpdateModelFile(t *testing.T) {
	tmp := t.TempDir()

	createFolders(tmp)

	m := before(tmp)

	generateModelDefConfFiles(1, tmp)

	id, _ := uuid.Parse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa0")

	e := model.Struct{
		Name: "model0",
		Type: model.Modbus,
		Ocpp: model.OcppModel{
			SerialNumb: "123",
		},
	}

	m.UpdateModel(id, &e)

	nf := tmp + defFolder + mod0

	assert.FileExists(t, nf)

	a := readFile(nf)

	assert.Equal(t, id, a.ID)

	assert.Equal(t, e.Name, a.Name)

	assert.Equal(t, e.Type, a.Type)

	assert.Equal(t, e.Ocpp.SerialNumb, a.Ocpp.SerialNumb)
}

func TestUpdateModelWrongId(t *testing.T) {
	tmp := t.TempDir()

	createFolders(tmp)

	m := before(tmp)

	generateModelDefConfFiles(1, tmp)

	id, _ := uuid.Parse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaab")

	e := model.Struct{
		Name: "model0",
		Type: model.Modbus,
		Ocpp: model.OcppModel{
			SerialNumb: "123",
		},
	}

	ab, as, ai, a := m.UpdateModel(id, &e)

	assert.False(t, ab)

	assert.Equal(t, m.L.Get(text.UpdateModelConfFileNotFound), as)

	assert.Equal(t, http.StatusNotFound, ai)

	assert.Empty(t, a)
}

func TestUpdateModelEmpty(t *testing.T) {
	tmp := t.TempDir()

	createFolders(tmp)

	m := before(tmp)

	generateModelDefConfFiles(1, tmp)

	id, _ := uuid.Parse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa0")

	e := model.Struct{}

	ab, as, ai, a := m.UpdateModel(id, &e)

	assert.False(t, ab)

	assert.Equal(t, m.L.Get(text.RequestBodyDoesntMatch), as)

	assert.Equal(t, http.StatusBadRequest, ai)

	assert.Empty(t, a)
}

func TestUpdateModelChangeName(t *testing.T) {
	tmp := t.TempDir()

	createFolders(tmp)

	m := before(tmp)

	generateModelDefConfFiles(1, tmp)

	id, _ := uuid.Parse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa0")

	e := model.Struct{
		Name: "nop",
		Type: model.Modbus,
		Ocpp: model.OcppModel{
			SerialNumb: "123",
		},
	}

	ep := e

	ab, as, ai, a := m.UpdateModel(id, &ep)

	assert.True(t, ab)

	assert.Empty(t, as)

	assert.Equal(t, http.StatusOK, ai)

	assert.NotEqual(t, e.Name, a.Name)

	assert.Equal(t, e.Type, a.Type)

	assert.Equal(t, e.Ocpp.SerialNumb, a.Ocpp.SerialNumb)
}

func TestDeleteModel(t *testing.T) {
	tmp := t.TempDir()

	createFolders(tmp)

	m := before(tmp)

	generateModelDefConfFiles(2, tmp)

	id, _ := uuid.Parse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa0")

	p0 := tmp + defFolder + mod0
	p1 := tmp + defFolder + "/model1.json"

	ab, as, ai := m.DeleteModel(id)

	assert.True(t, ab)

	assert.Empty(t, as)

	assert.Equal(t, http.StatusNoContent, ai)

	assert.NoFileExists(t, p0)

	assert.FileExists(t, p1)
}

func TestDeleteModelWrongId(t *testing.T) {
	tmp := t.TempDir()

	createFolders(tmp)

	m := before(tmp)

	generateModelDefConfFiles(2, tmp)

	id, _ := uuid.Parse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaab")

	p0 := tmp + defFolder + mod0
	p1 := tmp + defFolder + "/model1.json"

	ab, as, ai := m.DeleteModel(id)

	assert.False(t, ab)

	assert.Equal(t, m.L.Get(text.DeleteModelConfFileNotFound), as)

	assert.Equal(t, http.StatusNotFound, ai)

	assert.FileExists(t, p0)

	assert.FileExists(t, p1)
}

func before(tmp string) model.Model {
	t := translation.Translation{
		L: translation.Translation{}.GetKey("en-GB"),
	}

	m := model.Model{
		Scp: tmp + defSCFolder,
		L:   t,
	}

	return m
}

func generateModelDefConfFiles(n int64, tmp string) {
	for i := int64(0); i < n; i++ {
		id, _ := uuid.Parse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa" + strconv.FormatInt(i, 10))

		m := model.Struct{
			ID:   id,
			Name: "model" + strconv.FormatInt(i, 10),
			Type: model.Ocpp,
		}

		p := tmp + defFolder + "/model" + strconv.FormatInt(i, 10) + ".json"

		f, _ := os.Create(p)
		f.Close()

		b, _ := json.MarshalIndent(m, "", " ")

		os.WriteFile(p, b, fs.FileMode(common.FilePermissions))
	}
}

func generateInvalidModelConfFile(tmp string) {
	id, _ := uuid.Parse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa1")

	m := model.Struct{
		ID:   id,
		Name: "invalid",
	}

	p := tmp + defFolder + "/model_invalid.json"

	f, _ := os.Create(p)
	f.Close()

	b, _ := json.MarshalIndent(m, "", " ")

	os.WriteFile(p, b, fs.FileMode(common.FilePermissions))
}

func createFolders(tmp string) {
	os.Mkdir(tmp+defSCFolder, fs.FileMode(common.FolderPermissions))

	os.Mkdir(tmp+defFolder, fs.FileMode(common.FolderPermissions))
}

func readFile(p string) (m model.Struct) {
	b, _ := os.ReadFile(p)

	json.Unmarshal(b, &m)

	return m
}

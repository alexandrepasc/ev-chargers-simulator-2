//nolint:errcheck // because this is a test file
package model_test

import (
	"encoding/json"
	"io/fs"
	"os"
	"strconv"
	"testing"

	"github.com/alexandrepasc/ev-chargers-simulator-2/common"
	"github.com/alexandrepasc/ev-chargers-simulator-2/simulator/model"
	"github.com/alexandrepasc/ev-chargers-simulator-2/translation"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

const (
	defSCFolder string = "/simConf"
	defFolder   string = "/simConf" + common.DefMCFolder
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
			Name: "mod" + strconv.FormatInt(i, 10),
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

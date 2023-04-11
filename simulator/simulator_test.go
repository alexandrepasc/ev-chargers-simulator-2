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

func TestGetSimulatorData(t *testing.T) {
	createFolders()

	s := before()

	generateAssetConfFiles(1)

	generateModelConfFiles(1)

	al, ml := s.GetSimsConfs()

	assert.Equal(t, 1, len(al))

	assert.Equal(t, 1, len(ml))

	after()
}

func TestGetSimulatorDataNoFiles(t *testing.T) {
	createFolders()

	s := before()

	al, ml := s.GetSimsConfs()

	assert.Equal(t, 0, len(al))

	assert.Equal(t, 0, len(ml))

	after()
}

func TestSimulatorModelWrongType(t *testing.T) {
	createFolders()

	s := before()

	generateAssetConfFiles(1)

	generateInvalidModelConfFile()

	al, ml := s.GetSimsConfs()

	assert.Equal(t, 1, len(al))

	assert.Equal(t, 0, len(ml))

	after()
}

func before() simulator.Simulator {
	t := translation.Translation{
		L: translation.Translation{}.GetKey("en-GB"),
	}

	s := simulator.Simulator{
		Scp: common.DefSCPath,
		L:   t,
	}

	return s
}

func after() {
	os.RemoveAll(common.DefSCPath)
}

func createFolders() {
	os.Mkdir(common.DefSCPath, fs.FileMode(common.FolderPermissions))

	os.Mkdir(common.DefSCPath+common.DefMCFolder, fs.FileMode(common.FolderPermissions))
}

func generateAssetConfFiles(n int64) {
	for i := int64(0); i < n; i++ {
		a := simulator.Asset{
			Type:     simulator.Evc,
			Protocol: simulator.Ocpp201,
			Phases:   1,
		}

		p := common.DefSCPath + "/sim" + strconv.FormatInt(i, 10) + ".json"

		f, _ := os.Create(p)
		f.Close()

		b, _ := json.MarshalIndent(a, "", " ")

		os.WriteFile(p, b, fs.FileMode(common.FilePermissions))
	}
}

func generateModelConfFiles(n int64) {
	for i := int64(0); i < n; i++ {
		m := model.OcppModel{
			SerialNumb: "123456789",
			Model:      "model",
		}

		p := common.DefSCPath + common.DefMCFolder + "/model" + strconv.FormatInt(i, 10) + "_ocpp.json"

		f, _ := os.Create(p)
		f.Close()

		b, _ := json.MarshalIndent(m, "", " ")

		os.WriteFile(p, b, fs.FileMode(common.FilePermissions))
	}
}

func generateInvalidModelConfFile() {
	m := model.OcppModel{
		SerialNumb: "123456789",
		Model:      "model",
	}

	p := common.DefSCPath + common.DefMCFolder + "/model_invalid.json"

	f, _ := os.Create(p)
	f.Close()

	b, _ := json.MarshalIndent(m, "", " ")

	os.WriteFile(p, b, fs.FileMode(common.FilePermissions))
}

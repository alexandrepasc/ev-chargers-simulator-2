package model

import (
	"encoding/json"
	"io/fs"
	"os"

	"github.com/alexandrepasc/ev-chargers-simulator-2/common"
	"github.com/alexandrepasc/ev-chargers-simulator-2/translation"
	"github.com/alexandrepasc/ev-chargers-simulator-2/translation/text"
	"github.com/go-playground/validator/v10"
)

type Model struct {
	Scp string                  // Simulator configuration path
	L   translation.Translation // Translation language setting
	Ml  []*Struct               // Model structure list
}

/*
Get the model configuration files, reads them and convert them into an Struct model array.

Returns the array of the Struct with the data from the files ([]*Struct).
*/
func (m *Model) GetModels() (ml []*Struct) {
	var f = m.getModelsCofigFiles()

	ml = m.readModelsConfigFiles(f)

	m.Ml = ml

	return ml
}

/*
Get the list of files that are stored in the simulator model folder.

Returns an array of the file system entries stored in the folder ([]fs.DirEntry)
*/
func (m *Model) getModelsCofigFiles() []fs.DirEntry {
	common.Log("getModelsCofigFiles").Info(m.L.Get(text.GetModelsConfsFiles))

	f, err := os.ReadDir(m.Scp + common.DefMCFolder)

	if err != nil {
		common.Log("getModelsCofigFiles").Fatal(err)
	}

	return f
}

// TODO: complete the logic after implementing the modbus logic
/*
Read the models files, unmarshal each to the Model struct, and return and asset array ([]Struct).

f	-	File system entries list ([]fs.DirEntry)
*/
func (m *Model) readModelsConfigFiles(f []fs.DirEntry) (ml []*Struct) {
	common.Log("readModelsConfigFiles").Info(m.L.Get(text.ReadModelsConfsFiles))

	var p = m.Scp + common.DefMCFolder

	for _, entry := range f {
		if entry.IsDir() {
			continue
		}

		bv, err := os.ReadFile(p + "/" + entry.Name())

		if err != nil {
			common.Log("readModelsConfigFiles").Error(err)
		}

		m, ok := unmarshalModelJSON(bv)

		if ok {
			vErr := validator.New().Struct(m)

			if vErr != nil {
				common.Log("readModelsConfigFile").Error(err)

				continue
			}
		}

		ml = append(ml, &m)
	}

	return ml
}

/*
Convert configs json file (in bytes) to the Struct and returns it.

Besides the Model structure it returns a boolean false in case a reading fails.

b	-	Json file converted into bytes ([]byte)
*/
func unmarshalModelJSON(b []byte) (m Struct, ok bool) {
	err := json.Unmarshal(b, &m)

	if err != nil {
		common.Log("unmarshalModelJSON").Error(err)

		return m, false
	}

	return m, true
}

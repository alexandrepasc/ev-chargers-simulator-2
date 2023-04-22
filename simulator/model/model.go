package model

import (
	"encoding/json"
	"io/fs"
	"net/http"
	"os"

	"github.com/alexandrepasc/ev-chargers-simulator-2/common"
	"github.com/alexandrepasc/ev-chargers-simulator-2/translation"
	"github.com/alexandrepasc/ev-chargers-simulator-2/translation/text"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
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
Receive the model Struct, generates a new uuid, checks if the name does not exists, with
the Name field creates the file, and writes the data.

It returns true (bool), an empty message (string), the html code (int), and the new model with the
uuid (*Struct).

In case it fails returns false, the error message, the http code, and an empty struct.

s	-	Model struct with the new model data (*Struct)
*/
func (m *Model) CreateModel(s *Struct) (ok bool, msg string, code int, ns *Struct) {
	var id = uuid.New()

	s.ID = id

	vErr := validator.New().Struct(s)

	if vErr != nil {
		common.Log("CreateModel").Error(vErr)

		return false, vErr.Error(), http.StatusBadRequest, ns
	}

	var sl = m.GetModels()

	var nm = true

	for _, i := range sl {
		if i.Name == s.Name {
			nm = false
		}
	}

	if !nm {
		return false, m.L.Get(text.CreateModelConfFileNameExists), http.StatusConflict, ns
	}

	ok, f := generateFile(m.Scp+common.DefMCFolder, s.Name, m.L)

	if !ok {
		return ok, m.L.Get(text.CreateModelConfFileError), http.StatusBadRequest, ns
	}

	if !writeFile(f, s) {
		return false, m.L.Get(text.CreateModelConfFileError), http.StatusBadRequest, ns
	}

	return true, "", http.StatusCreated, s
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
Create a new json file with the configuration of a simulator.

It returns true and the file (*os.File), if not returns false and nil.

p	-	Path where the file will be created (string).

n	-	Name for the file (string).
*/
func generateFile(p, n string, l translation.Translation) (ok bool, f *os.File) {
	f, err := os.Create(p + "/" + n + ".json")

	if err != nil {
		common.Log("generateFile").Error(err)

		defer f.Close()

		return false, nil
	}

	common.Log("generateFile").Info(l.Get(text.CreateModelConfFile))

	return true, f
}

/*
Write to file the data for the new model.

Returns true if nothing fails, and false if it fails.

f	-	The file created to store the new model configurations (*os.File)

a	-	Model structure with the data to store to file (Struct)
*/
func writeFile(f *os.File, s *Struct) bool {
	ok, b := marshalModelToJSON(s)

	if !ok {
		return false
	}

	_, err := f.Write(b)

	defer f.Close()

	if err != nil {
		common.Log("writeFile").Error(err)

		return false
	}

	return true
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

/*
Marshal the model structure to byte.

Returns true and the model in bytes, if fails returns false and nil.

s	- Model structure (*Struct)
*/
func marshalModelToJSON(s *Struct) (ok bool, b []byte) {
	b, err := json.MarshalIndent(s, "", " ")

	if err != nil {
		common.Log("marshalModelToJSON").Error(err)

		return false, nil
	}

	return true, b
}

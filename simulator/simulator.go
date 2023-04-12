package simulator

import (
	"encoding/json"
	"io/fs"
	"os"
	"strings"

	"github.com/alexandrepasc/ev-chargers-simulator-2/common"
	"github.com/alexandrepasc/ev-chargers-simulator-2/simulator/model"
	"github.com/alexandrepasc/ev-chargers-simulator-2/translation"
	"github.com/alexandrepasc/ev-chargers-simulator-2/translation/text"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type Simulator struct {
	Scp string                  // simulator configuration path
	L   translation.Translation // Translation language setting
	Al  []Asset
	Oml []model.OcppModel
}

/*
Get the simulators configuration files and the model type files, read them and convert to structures.

Returns the lists of structures related to the two types of files.
*/
func (s *Simulator) GetSimsConfs() (al []Asset, oml []model.OcppModel) {
	var f = s.getSimConfigFIles()

	al = s.readSimConfigFiles(f)

	f = s.getModelsCofigFiles()

	oml = s.readModelsConfigFiles(f)

	s.Al = al

	s.Oml = oml

	return al, oml
}

/**/
func (s *Simulator) CreateSimConf(a *Asset) (ok bool, msg string, na *Asset) {
	id := uuid.New()

	a.SimID = id

	err := validator.New().Struct(a)

	if err != nil {
		common.Log("CreateSimConf").Error(err)

		return false, err.Error(), na
	}

	ok, f := generateFile(s.Scp, a.Name, s.L)

	if !ok {
		return ok, s.L.Get(text.CreateSimConfFileError), na
	}

	if !writeFile(f, a) {
		return false, s.L.Get(text.CreateSimConfFileError), na
	}

	return true, "", a
}

/*
Get the list of files that are stored in the simulator configurations folder.

Returns an array of the file system entries stored in the folder ([]fs.DirEntry)
*/
func (s *Simulator) getSimConfigFIles() []fs.DirEntry {
	common.Log("getSimConfigFIles").Info(s.L.Get(text.GetSimsConfsFiles))

	f, err := os.ReadDir(s.Scp)

	if err != nil {
		common.Log("getSimConfigFIles").Fatal(err)
	}

	return f
}

/*
Read the configuration files, unmarshal each to the Asset struct, and return and asset array ([]Asset).

f	-	File system entries list ([]fs.DirEntry)
*/
func (s *Simulator) readSimConfigFiles(f []fs.DirEntry) (as []Asset) {
	common.Log("readSimConfigFiles").Info(s.L.Get(text.ReadSimsConfsFiles))

	for _, entry := range f {
		if !entry.IsDir() {
			bv, err := os.ReadFile(s.Scp + "/" + entry.Name())

			if err != nil {
				common.Log("readSimConfigFiles").Error(err)
			}

			a, e := unmarshalAssetJSON(bv)

			if e {
				as = append(as, a)
			}
		}
	}

	return as
}

/*
Convert configs json file (in bytes) to the Asset and returns it.

Besides the Asset structure it returns a boolean false in case a reading fails.

b	-	Json file converted into bytes ([]byte)
*/
func unmarshalAssetJSON(b []byte) (a Asset, e bool) {
	err := json.Unmarshal(b, &a)

	if err != nil {
		common.Log("unmarshalAssetJSON").Error(err)

		return a, false
	}

	return a, true
}

/*
Get the list of files that are stored in the simulator model folder.

Returns an array of the file system entries stored in the folder ([]fs.DirEntry)
*/
func (s *Simulator) getModelsCofigFiles() []fs.DirEntry {
	common.Log("getModelsCofigFiles").Info(s.L.Get(text.GetModelsConfsFiles))

	f, err := os.ReadDir(s.Scp + common.DefMCFolder)

	if err != nil {
		common.Log("getModelsCofigFiles").Fatal(err)
	}

	return f
}

// TODO: complete the logic after implementing the modbus logic
/*
Read the models files, unmarshal each to the OcppModel struct, and return and asset array ([]OcppModel).

f	-	File system entries list ([]fs.DirEntry)
*/
func (s *Simulator) readModelsConfigFiles(f []fs.DirEntry) (oms []model.OcppModel) {
	common.Log("readModelsConfigFiles").Info(s.L.Get(text.ReadModelsConfsFiles))

	var p = s.Scp + common.DefMCFolder

	for _, entry := range f {
		if !entry.IsDir() {
			if strings.Contains(entry.Name(), "_") {
				ok, mt := checkModel(entry.Name())

				if ok {
					bv, err := os.ReadFile(p + "/" + entry.Name())

					if err != nil {
						common.Log("readModelsConfigFiles").Error(err)
					}

					switch mt {
					case string(model.Ocpp):
						m, e := unmarshalOcppJSON(bv)

						if e {
							oms = append(oms, m)
						}

					case string(model.Modbus):
					}
				}
			}
		}
	}

	return oms
}

/*
Convert configs json file (in bytes) to the OcppModel and returns it.

Besides the OcppModel structure it returns a boolean false in case a reading fails.

b	-	Json file converted into bytes ([]byte)
*/
func unmarshalOcppJSON(b []byte) (a model.OcppModel, e bool) {
	err := json.Unmarshal(b, &a)

	if err != nil {
		common.Log("unmarshalOcppJSON").Error(err)

		return a, false
	}

	return a, true
}

/*
Check if the file name match any of the model types supported.

At the moment it supports Ocpp and Modbus.

It returns a boolean var with true if it matches and the model type as string.

In case it doesn't match returns false and an empty string.

n	-	The file name to check (string)
*/
func checkModel(n string) (ok bool, mt string) {
	var aux = strings.Split(n, "_")[1]

	var s = strings.Split(aux, ".")[0]

	if strings.EqualFold(s, string(model.Ocpp)) {
		return true, string(model.Ocpp)
	}

	if strings.EqualFold(s, string(model.Modbus)) {
		return true, string(model.Modbus)
	}

	return false, ""
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

	common.Log("generateFile").Info(l.Get(text.CreateSimConfFile))

	return true, f
}

/*
Write to file the data for the new simulator.

Returns true if nothing fails, and false if it fails.

f	-	The file created to store the new simulator configurations (*os.File)

a	-	Asset structure with the data to store to file (Asset)
*/
func writeFile(f *os.File, a *Asset) bool {
	ok, b := marshalAssetToJSON(a)

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
Marshal the asset structure to byte.

Returns true and the asset in bytes, if fails returns false and nil.

a	- Asset structure (Asset)
*/
func marshalAssetToJSON(a *Asset) (ok bool, b []byte) {
	b, err := json.MarshalIndent(a, "", " ")

	if err != nil {
		common.Log("marshalAssetToJson").Error(err)

		return false, nil
	}

	return true, b
}

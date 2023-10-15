package simulator

import (
	"encoding/json"
	"io/fs"
	"net/http"
	"os"
	"strings"

	"github.com/alexandrepasc/ev-chargers-simulator-2/common"
	"github.com/alexandrepasc/ev-chargers-simulator-2/translation"
	"github.com/alexandrepasc/ev-chargers-simulator-2/translation/text"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type Simulator struct {
	Scp string                  // simulator configuration path
	L   translation.Translation // Translation language setting
	Al  []*Asset
}

/*
Get the simulator configuration files, reads them and convert them into an Asset model array.

Returns the array of the Asset with the data from the files ([]*Asset).
*/
func (s *Simulator) GetSimulators() (al []*Asset) {
	var f = s.getSimConfigFIles()

	al = s.readSimConfigFiles(f)

	s.Al = al

	return al
}

/*
Receive the Asset structure, generates the uuid for the new simulator, and with the Name field
creates the file to store the new configurations.

It will return a boolean as true if the process goes well and false otherwise (bool).

Will return an error message in case it fails, or an empty message (string).

The last returned value will be the Asset structure with the uuid, or empty one in fail case.

a	-	The new asset configurations (Asset)
*/
func (s *Simulator) CreateSimConf(a *Asset) (ok bool, msg string, na *Asset) {
	id := uuid.New()

	a.SimID = id

	err := validator.New().Struct(a)

	if err != nil {
		common.Log("CreateSimConf").Error(err)

		return false, err.Error(), &Asset{}
	}

	if a.Type == Evc {
		var o, m = validateEvcFields(a, s.L)

		if !o {
			return false, m, &Asset{}
		}
	}

	if strings.Contains(string(a.Protocol), "ocpp") {
		if a.CPId == "" {
			return false, s.L.Get(text.OcppMissingCPIdError), &Asset{}
		}
	}

	if a.Protocol == Ocpp16 {
		if len(a.Evses) > 1 {
			return false, s.L.Get(text.Ocpp16SimConfFileMoreEvse), &Asset{}
		}
	}

	if (a.TLS || a.BasicAuth) && a.Model.String() == "00000000-0000-0000-0000-000000000000" {
		return false, s.L.Get(text.MissingModelError), &Asset{}
	}

	var al = s.GetSimulators()

	var nm = true

	for _, i := range al {
		if i.Name == a.Name {
			nm = false
		}
	}

	if !nm {
		return false, s.L.Get(text.CreateSimConfFileNameExists), &Asset{}
	}

	ok, f := generateFile(s.Scp, a.Name, s.L)

	if !ok {
		return ok, s.L.Get(text.CreateSimConfFileError), &Asset{}
	}

	if !writeFile(f, a) {
		return false, s.L.Get(text.CreateSimConfFileError), &Asset{}
	}

	return true, "", a
}

/*
Receive the uuid of the simulator that will be updated and the changes. This will replace the
configurations entirely, it will not modify the only one value.

Returns true (bool), an empty string (string), the http code (int), and the new configuration
(*Asset) in case of success.

Will return false, the error message, the http code, and an empty Asset.

id	-	Simulator identifier (uuid.UUID)

a	-	Asset structure with the new configuration (*Asset)
*/
func (s *Simulator) UpdateSimConf(id uuid.UUID, a *Asset) (ok bool, msg string, code int, na *Asset) {
	var al = s.GetSimulators()

	for _, i := range al {
		if i.SimID != id {
			continue
		}

		na = a

		na.SimID = i.SimID

		na.Name = i.Name

		e := validator.New().Struct(a)

		if e != nil {
			common.Log("UpdateSimConf").Error(e)

			return false, s.L.Get(text.RequestBodyDoesntMatch), http.StatusBadRequest, &Asset{}
		}

		if strings.Contains(string(a.Protocol), "ocpp") {
			if a.CPId == "" {
				return false, s.L.Get(text.OcppMissingCPIdError), http.StatusBadRequest, &Asset{}
			}
		}

		if a.Type == Evc {
			var o, m = validateEvcFields(a, s.L)

			if !o {
				return false, m, http.StatusBadRequest, &Asset{}
			}
		}

		if a.Protocol == Ocpp16 {
			if len(a.Evses) > 1 {
				return false, s.L.Get(text.Ocpp16SimConfFileMoreEvse), http.StatusBadRequest, &Asset{}
			}
		}

		if (a.TLS || a.BasicAuth) && a.Model.String() == "00000000-0000-0000-0000-000000000000" {
			return false, s.L.Get(text.MissingModelError), http.StatusBadRequest, &Asset{}
		}

		common.Log("UpdateSimConf").Info(s.L.Get(text.UpdateSimConfFile))

		_, b := marshalAssetToJSON(na)

		err := os.WriteFile(s.Scp+"/"+na.Name+".json", b, fs.FileMode(common.FilePermissions))

		if err != nil {
			common.Log("UpdateSimConf").Error(err)

			return false, s.L.Get(text.WriteSimConfFileError), http.StatusTeapot, &Asset{}
		}

		return true, "", http.StatusOK, na
	}

	return false, s.L.Get(text.UpdateSimConfFileNotFound), http.StatusNotFound, &Asset{}
}

/*
Receive the uuid of the simulator that will be deleted. This will remove the configuration file
that matches the id.

Returns true (bool), an empty string (string) and the http code (int) in case of success.

Will return false, the error message, and the http code.

id	-	Simulator identifier (uuid.UUID)
*/
func (s *Simulator) DeleteSimConf(id uuid.UUID) (ok bool, msg string, code int) {
	var al = s.GetSimulators()

	for _, i := range al {
		if i.SimID != id {
			continue
		}

		err := os.Remove(s.Scp + "/" + i.Name + ".json")

		if err != nil {
			common.Log("DeleteSimConf").Error(err)

			return false, s.L.Get(text.DeleteSimConfFileError), http.StatusBadRequest
		}

		common.Log("DeleteSimConf").Info(s.L.Get(text.DeleteSimConfFile))

		return true, "", http.StatusNoContent
	}

	return false, s.L.Get(text.DeleteSimConfFileNotFound), http.StatusNotFound
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

// TODO: add struct validator after unmarshal
/*
Read the configuration files, unmarshal each to the Asset struct, and return and asset array ([]Asset).

f	-	File system entries list ([]fs.DirEntry)
*/
func (s *Simulator) readSimConfigFiles(f []fs.DirEntry) (as []*Asset) {
	common.Log("readSimConfigFiles").Info(s.L.Get(text.ReadSimsConfsFiles))

	for _, entry := range f {
		if !entry.IsDir() {
			bv, err := os.ReadFile(s.Scp + "/" + entry.Name())

			if err != nil {
				common.Log("readSimConfigFiles").Error(err)
			}

			a, e := unmarshalAssetJSON(bv)

			if e {
				as = append(as, &a)
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

/*
Execute the validation to the required properties of the evc asset type. Returns true (bool) and
an empty string if every validation pass, and returns false and the error message if it fails.

a	-	Asset structure (*Asset)

l	-	Translation module (translation.Translation)
*/
func validateEvcFields(a *Asset, l translation.Translation) (ok bool, msg string) {
	if !containsInt64(PhasesList, int64(a.Phases)) {
		return false, l.Get(text.EvcMissingPhasesError)
	}

	if !containsString(CurrentTypeList, string(a.CurrentType)) {
		return false, l.Get(text.EvcMissingCurrentTypeError)
	}

	if len(a.Evses) == 0 {
		return false, l.Get(text.EvcMissingEvsesError)
	}

	return true, ""
}

/*
Runs the list of int64 and checks if the value matches any of the list items, if it match return
true (bool) if not returns false.

il	-	List of data to validate ([]int64)

c	-	Item to be checked if exists in the list (int64)
*/
func containsInt64(il []int64, v int64) bool {
	for _, i := range il {
		if i == v {
			return true
		}
	}

	return false
}

/*
Runs the list of string and checks if the value matches any of the list items, if it match return
true (bool) if not returns false.

sl	-	List of data to validate ([]string)

c	-	Item to be checked if exists in the list (string)
*/
func containsString(sl []string, v string) bool {
	for _, s := range sl {
		if s == v {
			return true
		}
	}

	return false
}

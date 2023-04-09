package translation

import (
	"github.com/alexandrepasc/ev-chargers-simulator-2/common"
	"github.com/alexandrepasc/ev-chargers-simulator-2/translation/text"
)

/*
Structure used to set the language.
*/
type Translation struct {
	L Local // Localization of the language (Local)
}

/*
Will use the language set in the structure and a key to retrieve the text from the map.

Returns a string with the text associated to the key. In case no key is set log fatal.

k	-	key identifying the text to be retrieved (text.Key)
*/
func (t Translation) Get(k text.Key) string {
	if t.L == "" {
		common.Log("Get").Fatal(location[EnGb].(map[text.Key]string)[text.LocalizationNotSet])
	}

	return location[t.L].(map[text.Key]string)[k]
}

/*
List the languages supported by the application.

Returns an array with the keys of the language location ([]string)
*/
func (t Translation) List() []string {
	keys := make([]string, 0, len(location))

	for k := range location {
		keys = append(keys, string(k))
	}

	return keys
}

/*
Get the location key using the string associated to it, this is useful when getting the string from
the general settings and converting to the location key.

It returns the matching key, in case there is no matching key it will return the default one.

s	-	The string from the general settings file (string)
*/
func (t Translation) GetKey(s string) Local {
	for k := range location {
		if string(k) == s {
			return k
		}
	}

	return Local(common.DefLanguage)
}

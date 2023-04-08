package translation

import "github.com/alexandrepasc/ev-chargers-simulator-2/translation/text"

/*
Structure used to set the language.
*/
type Translation struct {
	L Local // Localization of the language (Local)
}

/*
Will use the language set in the structure and a key to retrieve the text from the map.

Returns a string with the text associated to the key.

k	-	key identifying the text to be retrieved (text.Key)
*/
func (t Translation) Get(k text.Key) string {
	return location[t.L].(map[text.Key]string)[k]
}

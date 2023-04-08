package translation

import "github.com/alexandrepasc/ev-chargers-simulator-2/translation/text"

type Local string // Type to be used in the localization

const (
	EnGb Local = "en-GB" // en-GB language localization
	PtPt Local = "pt-PT" // pt-PT language localization
 )

/*
Maps the language with the map of each languages, where the texts are retrieved
*/
var location = map[Local]interface{}{
	EnGb: text.EnGB,
	PtPt: text.PtPT,
}

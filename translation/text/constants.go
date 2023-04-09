package text

type Key string // Type used in the get translation text

/*
List of the keys that are mapped to text.
*/
const (
	LocalizationNotSet Key = "localizationNotSet" // Localization not set when retreiving text key
	RouteNotFound      Key = "routeNotFound"      // API message returned when the route requested is not mapped
	MethodNotAllowed   Key = "methodNotAllowed"   // API message returned when the method requested doesn't match the route
)

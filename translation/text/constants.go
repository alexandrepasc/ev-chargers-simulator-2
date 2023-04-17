package text

type Key string // Type used in the get translation text

/*
List of the keys that are mapped to text.
*/
const (
	FolderNotExist              Key = "FolderNotExist"             // Folder doesn't exist error
	FolderCreated               Key = "FolderCreated"              // Info folder created
	ConfigsNotExist             Key = "ConfigsNotExist"            // Configuration file doesn't exist
	ConfigsCreated              Key = "ConfigsCreated"             // Configuration file created
	ConfigsWritten              Key = "ConfigsWritten"             // Configurations saved
	GeneralNotExist             Key = "GeneralNotExist"            // General settings file doesn't exist
	GeneralCreated              Key = "GeneralCreated"             // General settings file created
	GeneralWritten              Key = "GeneralWritten"             // General settings saved
	LocalizationNotSet          Key = "LocalizationNotSet"         // Localization not set when retreiving text key
	RouteNotFound               Key = "RouteNotFound"              // API message returned when the route requested is not mapped
	MethodNotAllowed            Key = "MethodNotAllowed"           // API message returned when the method requested doesn't match the route
	GetSimsConfsFiles           Key = "GetSimsConfsFiles"          // Get the simulator configuration files
	ReadSimsConfsFiles          Key = "ReadSimsConfsFiles"         // Read the simulator configuration files
	GetModelsConfsFiles         Key = "GetModelsConfsFiles"        // Get the models configuration files
	ReadModelsConfsFiles        Key = "ReadModelsConfsFiles"       // Read the models configuration files
	CreateSimConfFile           Key = "CreateSimConfFile"          // Create a new simulator configuration file
	CreateSimConfFileError      Key = "CreateSimConfFileError"     // Unable to create the simulator configuration file
	CreateSimConfFileNameExists Key = "CreateSimConfFileNameExist" // The name of the new asset already exists
	UpdateSimConfFileNotFound   Key = "UpdateSimConfFileNotFound"  // Not found the configuration file
	UpdateSimConfFile           Key = "UpdateSimConfFile"          // Update simulator configuration file
	DeleteSimConfFileNotFound   Key = "DeleteSimConfFileNotFound"  // Configuration ID not found
	OpenSimConfFileError        Key = "OpenSimConfFileError"       // Error opening the simulator configuration file
	WriteSimConfFileError       Key = "WriteSimConfFileError"      // Error writing the simulator configuration file
	InternalServerError         Key = "InternalServerError"        // API message returned when something breaks
	RequestBodyDoesntMatch      Key = "RequestBodyDoesntMatch"     // API message returned when the request body doesn't bind to the model
	UUIDParsingError            Key = "UUIDParsingError"           // API message returned when couldn't parse the id to uuid
)

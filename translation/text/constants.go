package text

type Key string // Type used in the get translation text

/*
List of the keys that are mapped to text.
*/
const (
	FolderNotExist                Key = "FolderNotExist"                // Folder doesn't exist error
	FolderCreated                 Key = "FolderCreated"                 // Info folder created
	ConfigsNotExist               Key = "ConfigsNotExist"               // Configuration file doesn't exist
	ConfigsCreated                Key = "ConfigsCreated"                // Configuration file created
	ConfigsWritten                Key = "ConfigsWritten"                // Configurations saved
	ConfigsErrorRead              Key = "ConfigsErrorRead"              // Configurations error get the data
	ConfigsErrorUpdate            Key = "ConfigsErrorUpdate"            // Configurations error updating data
	GeneralNotExist               Key = "GeneralNotExist"               // General settings file doesn't exist
	GeneralCreated                Key = "GeneralCreated"                // General settings file created
	GeneralWritten                Key = "GeneralWritten"                // General settings saved
	GeneralErrorRead              Key = "GeneralErrorRead"              // General error get the configuration
	GeneralErrorUpdate            Key = "GeneralErrorUpdate"            // General error updating the configuration
	LocalizationNotSet            Key = "LocalizationNotSet"            // Localization not set when retreiving text key
	RouteNotFound                 Key = "RouteNotFound"                 // API message returned when the route requested is not mapped
	MethodNotAllowed              Key = "MethodNotAllowed"              // API message returned when the method requested doesn't match the route
	GetSimsConfsFiles             Key = "GetSimsConfsFiles"             // Get the simulator configuration files
	ReadSimsConfsFiles            Key = "ReadSimsConfsFiles"            // Read the simulator configuration files
	GetModelsConfsFiles           Key = "GetModelsConfsFiles"           // Get the models configuration files
	ReadModelsConfsFiles          Key = "ReadModelsConfsFiles"          // Read the models configuration files
	CreateModelConfFile           Key = "CreateModelConfFile"           // Create a new model configuration file
	CreateModelConfFileError      Key = "CreateModelConfFileError"      // Unable to create the model configuration file
	CreateModelConfFileNameExists Key = "CreateModelConfFileNameExists" // The name of the new model already exists
	UpdateModelConfFileNotFound   Key = "UpdateModelConfFileNotFound"   // Not found the configuration file
	UpdateModelConfFile           Key = "UpdateModelConfFile"           // Update model configuration file
	UpdateModelConfFileError      Key = "UpdateModelConfFileError"      // Update model configuration file error
	DeleteModelConfFileNotFound   Key = "DeleteModelConfFileNotFound"   // Configuration ID not found
	DeleteModelConfFileError      Key = "DeleteModelConfFileError"      // Not able to delete the configuration file
	DeleteModelConfFile           Key = "DeleteModelConfFile"           // Delete model configuration file
	CreateSimConfFile             Key = "CreateSimConfFile"             // Create a new simulator configuration file
	CreateSimConfFileError        Key = "CreateSimConfFileError"        // Unable to create the simulator configuration file
	CreateSimConfFileNameExists   Key = "CreateSimConfFileNameExist"    // The name of the new asset already exists
	UpdateSimConfFileNotFound     Key = "UpdateSimConfFileNotFound"     // Not found the configuration file
	UpdateSimConfFile             Key = "UpdateSimConfFile"             // Update simulator configuration file
	DeleteSimConfFileNotFound     Key = "DeleteSimConfFileNotFound"     // Configuration ID not found
	DeleteSimConfFileError        Key = "DeleteSimConfFileError"        // Not able to delete the configuration file
	DeleteSimConfFile             Key = "DeleteSimConfFile"             // Delete simulator configuration file
	Ocpp16SimConfFileMoreEvse     Key = "Ocpp16SimConfFileMoreEvse"     // Ocpp 1.6 only can handle 1 evse per charge point error
	Ocpp201ServerStarted          Key = "Ocpp201ServerStarted"          // Simulator ocpp 2.0.1 server started
	Ocpp201ServerStopped          Key = "Ocpp201ServerStopped"          // Simulator ocpp 2.0.1 server stopped
	EvcMissingPhasesError         Key = "EvcMissingPhasesError"         // The phases property is required
	EvcMissingCurrentTypeError    Key = "EvcMissingCurrentTypeError"    // The current type property is required
	EvcMissingEvsesError          Key = "EvcMissingEvsesError"          // The evses property is required
	OcppMissingCPIdError          Key = "OcppMissingCPIdError"          // The CP Id property is required
	MissingModelError             Key = "MissingModelError"             // The model property is required
	OpenSimConfFileError          Key = "OpenSimConfFileError"          // Error opening the simulator configuration file
	WriteSimConfFileError         Key = "WriteSimConfFileError"         // Error writing the simulator configuration file
	InternalServerError           Key = "InternalServerError"           // API message returned when something breaks
	RequestBodyDoesntMatch        Key = "RequestBodyDoesntMatch"        // API message returned when the request body doesn't bind to the model
	UUIDParsingError              Key = "UUIDParsingError"              // API message returned when couldn't parse the id to uuid
	NoAssetsToRunError            Key = "NoAssetsToRunError"            // API message returned when there are no assets to run
	ModbusServerStarted           Key = "ModbusServerStarted"           // Simulator modbus server started
	ModbusServerStopped           Key = "ModbusServerStopped"           // Simulator modbus server stopped
	CaCertNotFound                Key = "CaCertNotFound"                // No CA certificate found
	StartUpConfigurations         Key = "StartUpConfigurations"         // Set the startup configurations for an asset
)

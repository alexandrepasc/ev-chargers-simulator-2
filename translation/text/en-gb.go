package text

/*
EnGB maps the keys defined with the message for the language/location en-GB
*/
var EnGB = map[Key]string{
	FolderNotExist:                "Folder doesn't exist ",
	FolderCreated:                 "Folder created ",
	ConfigsNotExist:               "Configurations file doesn't exist.",
	ConfigsCreated:                "Created configuration file",
	ConfigsWritten:                "Configurations written to file.",
	GeneralNotExist:               "General configuration file doesn't exist.",
	GeneralCreated:                "Created general configuration file.",
	GeneralWritten:                "General configurations written to file.",
	LocalizationNotSet:            "Language localization not set.",
	RouteNotFound:                 "Route not found.",
	MethodNotAllowed:              "Method not allowed.",
	GetSimsConfsFiles:             "Get simulators configuration files.",
	ReadSimsConfsFiles:            "Read simulators configuration files.",
	GetModelsConfsFiles:           "Get models configuration files.",
	CreateModelConfFile:           "Created model configuration file.",
	CreateModelConfFileError:      "Unexpected error creating the model configuration file.",
	CreateModelConfFileNameExists: "The model name already exists.",
	ReadModelsConfsFiles:          "Read models configuration files.",
	CreateSimConfFile:             "Created simulator configuration file.",
	CreateSimConfFileError:        "Unexpected error creating the simulator configuration file.",
	CreateSimConfFileNameExists:   "The simulator name already exists.",
	UpdateSimConfFileNotFound:     "Simulator not found.",
	UpdateSimConfFile:             "Update simulator configuration file.",
	DeleteSimConfFileNotFound:     "Simulator not found.",
	DeleteSimConfFileError:        "Unexpected error deleting the simulator configuration file.",
	DeleteSimConfFile:             "Delete simulator configuration file.",
	OpenSimConfFileError:          "Error opening the simulator configuration file.",
	WriteSimConfFileError:         "Error writing the simulator configuration file.",
	InternalServerError:           "Something went very wrong.",
	RequestBodyDoesntMatch:        "The request body is malformed.",
	UUIDParsingError:              "The id could not be parsed to UUID.",
}

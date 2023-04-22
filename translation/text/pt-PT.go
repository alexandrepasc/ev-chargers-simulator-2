package text

/*
PtPT maps the keys defined with the message for the language/location en-GB
*/
var PtPT = map[Key]string{
	RouteNotFound:                 "Rota não encontrada.",
	MethodNotAllowed:              "Método não permitido.",
	GetSimsConfsFiles:             "Obter os arquivos de configuração dos simuladores.",
	ReadSimsConfsFiles:            "Ler os arquivos de configuração dos simuladores.",
	GetModelsConfsFiles:           "Obter os arquivos de configuração dos modelos.",
	CreateModelConfFile:           "Ficheiro de configuração do modelo criado.",
	CreateModelConfFileError:      "Erro inesperado ao criar o ficheiro de configuração do modelo.",
	CreateModelConfFileNameExists: "O nome do modelo já existe.",
	ReadModelsConfsFiles:          "Ler os arquivos de configuração dos modelos.",
	CreateSimConfFile:             "Fichheiro de configuração do simulador criado.",
	CreateSimConfFileError:        "Erro inesperado ao criar o ficheiro de configuração do simulador.",
	CreateSimConfFileNameExists:   "O nome do simulador já existe.",
	UpdateSimConfFileNotFound:     "Simulador não encontrado.",
	UpdateSimConfFile:             "Actualiza o ficheiro de configuração do simulador.",
	DeleteSimConfFileNotFound:     "Simulador não encontrado.",
	DeleteSimConfFileError:        "Erro inesperado ao apagar o ficheiro de configuração do simulador.",
	DeleteSimConfFile:             "Apagar o ficheiro de configuração do simulador.",
	OpenSimConfFileError:          "Erro a abrir o ficheiro de configuração do simulador.",
	WriteSimConfFileError:         "Erro a escrever no ficheiro de configuração do simulador.",
	InternalServerError:           "Algo correu muito mal.",
	RequestBodyDoesntMatch:        "O corpo do pedido está mal formado.",
	UUIDParsingError:              "O id não pôde ser convertido para UUID.",
}

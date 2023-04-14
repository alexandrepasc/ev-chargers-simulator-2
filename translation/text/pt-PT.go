package text

/*
PtPT maps the keys defined with the message for the language/location en-GB
*/
var PtPT = map[Key]string{
	RouteNotFound:               "Rota não encontrada.",
	MethodNotAllowed:            "Método não permitido.",
	GetSimsConfsFiles:           "Obter os arquivos de configuração dos simuladores.",
	ReadSimsConfsFiles:          "Ler os arquivos de configuração dos simuladores.",
	GetModelsConfsFiles:         "Obter os arquivos de configuração dos modelos.",
	ReadModelsConfsFiles:        "Ler os arquivos de configuração dos modelos.",
	CreateSimConfFile:           "Fichheiro de configuração do simulador criado.",
	CreateSimConfFileError:      "Erro inesperado ao criar o ficheiro de configuração do simulador.",
	CreateSimConfFileNameExists: "O nome do simulador já existe.",
	InternalServerError:         "Algo correu muito mal.",
	RequestBodyDoesntMatch:      "O corpo do pedido está mal formado.",
}

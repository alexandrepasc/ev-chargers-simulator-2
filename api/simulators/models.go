package simulators

type GetSimulatorsTemp struct {
	Total  int64       `json:"total"`  // Total number of assets
	Assets interface{} `json:"assets"` // List of assets
}

type GetModelsTemp struct {
	Total  int64       `json:"total"`  // Total number of models
	Models interface{} `json:"models"` // List of models
}

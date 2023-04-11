package simulators

type getSimulatorsTemp struct {
	Total  int64       `json:"total"`  // Total number of assets
	Assets interface{} `json:"assets"` // List of assets
}

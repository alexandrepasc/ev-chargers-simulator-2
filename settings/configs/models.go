package configs

type Model struct {
	GeneralConfigFolder    string `json:"generalFolder"` // General configuration folder path
	SimulatorsConfigFolder string `json:"simFolder"`     // Simulators configuration folder path
}

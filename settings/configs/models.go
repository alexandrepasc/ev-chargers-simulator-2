package configs

type Model struct {
	GeneralConfigFolder    string `json:"generalFolder"` // General configuration folder path
	SimulatorsConfigFolder string `json:"simFolder"`     // Simulators configuration folder path
}

/*
This model exists to be used by the api logic, since the api only has the ability to change
the simulators configuration path.
*/
type SimPathModel struct {
	SimulatorsConfigFolder string `json:"simFolder" validate:"required"` // Simulators configuration folder path
}

package simulator

type Asset struct {
	Type          string `json:"type"`          // Type of the asset that the configuration will be used to (evc, pm)
	Protocol      string `json:"protocol"`      // Protocol used by the asset
	Model         string `json:"model"`         // The name of the configuration file set in the models folder
	Port          string `json:"port"`          // Communication ip port
	CPId          string `json:"cPId"`          // Charge point id to identify the unit (used in the ocpp protocol)
	StartCharging bool   `json:"startCharging"` // Set the asset to start charging behaviour by it self
	Phases        int64  `json:"phases"`        // Phases number
	CurrentType   string `json:"curentType"`    // Type of current of the asset (AC or DC)
	Evses         []Evse `json:"evses"`         // List of evses that the asset has
}

type Evse struct {
	Id         int64       `json:"id"`         // Evse identifier number
	Connectors []Connector `json:"connectors"` // The list of connectors of the evse
}

type Connector struct {
	Id      int64  `json:"id"` // Connector identifier number
	Enabled bool   // Since only one connector can be chargin this is just to control that state
	Data    []Data `json:"data"` // The loop of data
}

// TODO: evaluate if the charging state should be an number or the name of the state
type Data struct {
	Duration      int64 `json:"duration"`      // The duration in seconds that the current data will be in place
	ChargingState int64 `json:"chargingState"` // The state of charging for the current data
	ErrorCode     int64 `json:"errorCode"`     // Error code
	PowerFactor   int64 `json:"powerFactor"`   // The power factor
	Power         int64 `json:"power"`         // Power in W
	VoltageP1     int64 `json:"voltageP1"`     // Voltage of the phase 1 in V (in case the asset is single phase only need the voltageP1)
	VoltageP2     int64 `json:"voltageP2"`     // Voltage of the phase 2 in V (in case the asset is single phase only need the voltageP1)
	VoltageP3     int64 `json:"voltageP3"`     // Voltage of the phase 3 in V (in case the asset is single phase only need the voltageP1)
}

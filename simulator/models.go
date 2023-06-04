package simulator

import "github.com/google/uuid"

type Asset struct {
	SimID         uuid.UUID   `json:"simId" validate:"required"`                                // Simulator identifier
	Name          string      `json:"name" validate:"required"`                                 // Simulator name
	Type          AssetType   `json:"type" validate:"required,oneof=evc pm"`                    // Type of the asset that the configuration will be used to (evc, pm)
	Protocol      Protocol    `json:"protocol" validate:"required,oneof=ocpp201 ocpp16 modbus"` // Protocol used by the asset
	Model         uuid.UUID   `json:"model"`                                                    // Id of the configuration file set in the model's folder
	Port          string      `json:"port,omitempty"`                                           // Communication ip port
	CPId          string      `json:"cPId,omitempty"`                                           // Charge point id to identify the unit (used in the ocpp protocol)
	StartCharging bool        `json:"startCharging" validate:"boolean"`                         // Set the asset to start charging behaviour by itself
	Phases        Phases      `json:"phases" validate:"required,oneof=1 3"`                     // Phases number
	CurrentType   CurrentType `json:"curentType" validate:"required,oneof=ac dc"`               // Type of current of the asset (AC or DC)
	Evses         []Evse      `json:"evses" validate:"required"`                                // List of evses that the asset has
}

type Evse struct {
	ID         int64       `json:"id"`         // Evse identifier number
	Connectors []Connector `json:"connectors"` // The list of connectors of the evse
}

type Connector struct {
	ID      int64        `json:"id"`               // Connector identifier number
	Enabled bool         `json:"omitempty"`        // Since only one connector can be charging this is just to control that state
	Data    []Data       `json:"data"`             // The loop of data
	DP      DataPosition `json:"dp,omitempty"`     // Data position used to control the data
	TPower  float64      `json:"tPower,omitempty"` // Store the total power for the connector, this will be used by the application only
}

// TODO: evaluate if the charging state should be a number or the name of the state
type Data struct {
	Duration      int64   `json:"duration"`      // The duration in seconds that the current data will be in place
	ChargingState int64   `json:"chargingState"` // The state of charging for the current data
	ErrorCode     int64   `json:"errorCode"`     // Error code
	PowerFactor   int64   `json:"powerFactor"`   // The power factor
	Power         int64   `json:"power"`         // Power in W
	Voltage       []int64 `json:"voltage"`       // Array of voltages each entry for each phase in V
}

type DataPosition struct {
	Position int64
	Ticker   int64
}

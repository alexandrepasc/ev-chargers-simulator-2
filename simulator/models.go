package simulator

import "github.com/google/uuid"

type Asset struct {
	SimID           uuid.UUID     `json:"simId" validate:"required"`                                                                        // Simulator identifier
	Name            string        `json:"name" validate:"required"`                                                                         // Simulator name
	Type            AssetType     `json:"type" validate:"required,oneof=evc pm"`                                                            // Type of the asset that the configuration will be used to (evc, pm)
	Protocol        Protocol      `json:"protocol" validate:"required,oneof=ocpp201 ocpp16 modbus"`                                         // Protocol used by the asset
	TLS             bool          `json:"tls"`                                                                                              // Set if the asset will connect using TLS or not, the default is false (if yes model is required)
	Model           uuid.UUID     `json:"model"`                                                                                            // Id of the configuration file set in the model's folder
	Port            string        `json:"port,omitempty"`                                                                                   // Communication ip port
	CPId            string        `json:"cPId,omitempty"`                                                                                   // Charge point id to identify the unit (used in the ocpp protocol)
	StartCharging   bool          `json:"startCharging" validate:"boolean"`                                                                 // Set the asset to start charging behaviour by itself
	Phases          Phases        `json:"phases,omitempty" validate:"omitempty,oneof=1 3"`                                                  // Phases number
	PhaseRotation   PhaseRotation `json:"phaseRotation,omitempty" validate:"omitempty,oneof=NotApplicable Unknown RST RTS SRT STR TRS TSR"` // The asset phase rotation, if the asset is DC the value should be NotApplicable
	CurrentType     CurrentType   `json:"currentType,omitempty" validate:"omitempty,oneof=ac dc"`                                           // Type of current of the asset (AC or DC)
	AuthorizeRemote bool          `json:"authorizeRemote" validate:"boolean"`                                                               // Configurataion AuthorizeRemoteTxRequests
	AuthList        bool          `json:"authList" validate:"boolean"`                                                                      // Enable or disable authorization local list
	Evses           []Evse        `json:"evses,omitempty" validate:"omitempty,required"`                                                    // List of evses that the asset has
}

type Evse struct {
	ID         int64       `json:"id"`               // Evse identifier number
	Connectors []Connector `json:"connectors"`       // The list of connectors of the evse
	CIDTag     string      `json:"cIdTag,omitempty"` // current id tag being used in the charge session
}

type Connector struct {
	ID           int64        `json:"id"`                     // Connector identifier number
	Enabled      bool         `json:"omitempty"`              // Since only one connector can be charging this is just to control that state
	Data         []Data       `json:"data"`                   // The loop of data
	DP           DataPosition `json:"dp,omitempty"`           // Data position used to control the data
	TPower       float64      `json:"tPower,omitempty"`       // Store the total power for the connector, this will be used by the application only
	TPowerExport float64      `json:"tPowerExport,omitempty"` // Store the total power export for the connector, this will be used by the application only
	Energy       float64      `json:"energy,omitempty"`       // store the energy of the connector, this will be used by the application only
	EnergyExport float64      `json:"energyExport,omitempty"` // store the energy export of the connector, this will be used by the application only
	Availability string       `json:"availability,omitempty"` // Store the availability state of the connector, this will only be used by the application
	CurrentSoC   float64      `json:"currentSoC,omitempty"`   // Store the current ev state of charge in the current position, this will only be used by the applicaion
}

// TODO: evaluate if the charging state should be a number or the name of the state
type Data struct {
	Duration      int64   `json:"duration"`      // The duration in seconds that the current data will be in place
	ChargingState int64   `json:"chargingState"` // The state of charging for the current data
	ErrorCode     int64   `json:"errorCode"`     // Error code
	PowerFactor   int64   `json:"powerFactor"`   // The power factor
	Power         int64   `json:"power"`         // Power in W
	PowerExport   int64   `json:"powerExport"`   // Power exported by the ev to the cp in W, vehicle to grid
	Voltage       []int64 `json:"voltage"`       // Array of voltages each entry for each phase in V
	StartSoC      float64 `json:"startSoC"`      // Start ev charge state in percentage
	EndSoC        float64 `json:"endSoC"`        // End ev charge state in percentage
}

type DataPosition struct {
	Position int64 // Current data array position
	Ticker   int64 // Time in seconds that is in the position
}

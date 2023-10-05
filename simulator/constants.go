package simulator

const (
	Evc           AssetType     = "evc"           // Type electric vehicle charger
	Pm            AssetType     = "pm"            // Type power meter
	Ac            CurrentType   = "ac"            // Alternating current
	Dc            CurrentType   = "dc"            // Direct current
	Ocpp201       Protocol      = "ocpp201"       // Ocpp version 2.0.1
	Ocpp16        Protocol      = "ocpp16"        // Ocpp version 1.6
	Modbus        Protocol      = "modbus"        // Modbus protocol
	One           Phases        = 1               // One phase
	Three         Phases        = 3               // Three phases
	NotApplicable PhaseRotation = "NotApplicable" // Not applicable for dc chargers
	Unknown       PhaseRotation = "Unknown"       // Not able to retrieve the rotation
	RST           PhaseRotation = "RST"           // L1 L2 L3
	RTS           PhaseRotation = "RTS"           // L1 L3 L2
	SRT           PhaseRotation = "SRT"           // L2 L1 L3
	STR           PhaseRotation = "STR"           // L2 L3 L1
	TRS           PhaseRotation = "TRS"           // L3 L1 L2
	TSR           PhaseRotation = "TSR"           // L3 L2 L1
)

var (
	PhasesList      = []int64{1, 3}
	CurrentTypeList = []string{"ac", "dc"}
)

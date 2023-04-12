package simulator

const (
	Evc     AssetType   = "evc"     // Type electric vehicle charger
	Pm      AssetType   = "pm"      // Type power meter
	Ac      CurrentType = "ac"      // Alternating current
	Dc      CurrentType = "dc"      // Direct current
	Ocpp201 Protocol    = "ocpp201" // Ocpp version 2.0.1
	Ocpp16  Protocol    = "ocpp16"  // Ocpp version 1.6
	Modbus  Protocol    = "modbus"  // Modbus protocol
	One     Phases      = 1         // One phase
	Three   Phases      = 3         // Three phases
)

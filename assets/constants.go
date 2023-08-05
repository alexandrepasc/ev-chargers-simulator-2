package assets

const (
	Info  Severity = "info"
	Error Severity = "error"
	Warn  Severity = "warn"
	Panic Severity = "panic"
	Fatal Severity = "fatal"

	Request  string = "request"
	Response string = "response"
	CP       string = "charge point"
	CS       string = "central system"

	Available     ConnectorStatus = 1
	Preparing     ConnectorStatus = 2
	Charging      ConnectorStatus = 3
	SuspendedEV   ConnectorStatus = 4
	SuspendedEVSE ConnectorStatus = 5
	Finishing     ConnectorStatus = 6
	Reserved      ConnectorStatus = 7
	Unavailable   ConnectorStatus = 8
	Faulted       ConnectorStatus = 9

	Soft ResetType = "Soft"
	Hard ResetType = "Hard"

	Inoperative Availability = "Inoperative"
	Operative   Availability = "Operative"

	DefMeterValueSampleInterval          int64  = 20
	DefMeterValuesSampledData            string = "Energy.Active.Import.Register,Power.Active.Import"
	DefStopTxnSampledData                string = "Current.Import"
	DefClockAlignedDataInterval          int64  = 900
	DefMeterValuesAlignedData            string = "Energy.Active.Import.Register,Power.Active.Import"
	DefStopTxnAlignedData                string = "Power.Active.Import"
	DefHeartbeatInterval                 int64  = 30
	DefResetRetries                      int64  = 0
	DefStopTransactionOnEvSideDisconnect bool   = false
	DefStopTransactionOnInvalidID        bool   = true
)

var (
	Current                      = "Current"
	Power                        = "Power"
	EnergyActiveExportRegister   = "Energy.Active.Export.Register"
	EnergyActiveImportRegister   = "Energy.Active.Import.Register"
	EnergyReactiveImportRegister = "Energy.Reactive.Import.Register"
	Voltage                      = "Voltage"
	CurrentImport                = "Current.Import"
	PowerActiveImport            = "Power.Active.Import"
	Phases                       = []string{"L1", "L2", "L3"}
	Status                       = []string{
		"",
		"Available",
		"Preparing",
		"Charging",
		"SuspendedEV",
		"SuspendedEVSE",
		"Finishing",
		"Reserved",
		"Unavailable",
		"Faulted",
	}
	ErrorCode = []string{
		"NoError",
		"ConnectorLockFailure",
		"EVCommunicationError",
		"GroundFailure",
		"HighTemperature",
		"InternalError",
		"LocalListConflict",
		"OtherError",
		"OverCurrentFailure",
		"OverVoltage",
		"PowerMeterFailure",
		"PowerSwitchFailure",
		"ReaderFailure",
		"ResetFailure",
		"UnderVoltage",
		"WeakSignal",
	}
)

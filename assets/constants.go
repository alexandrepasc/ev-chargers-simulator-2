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
	DefSupportedFeatureProfiles          string = "Core,LocalAuthListManagement,RemoteTrigger"
	DefTransactionMessageAttempts        int64  = 1
	DefTransactionMessageRetryInterval   int64  = 10
	DefUnlockConnectorOnEVSideDisconnect bool   = true
	DefGetConfigurationMaxKeys           int64  = 10
	DefLocalAuthorizeOffline             bool   = true
	DefLocalPreAuthorize                 bool   = false
	DefHeartbeatInterval                 int64  = 30
	DefResetRetries                      int64  = 0
	DefStopTransactionOnEvSideDisconnect bool   = false
	DefStopTransactionOnInvalidID        bool   = true
)

var (
	Current                      = "Current"
	Power                        = "Power"
	CurrentExport                = "Current.Export"                  // Instantaneous current flow from EV
	CurrentImport                = "Current.Import"                  // Instantaneous current flow to EV
	CurrentOffered               = "Current.Offered"                 // Maximum current offered to EV
	EnergyActiveExportRegister   = "Energy.Active.Export.Register"   // Numerical value read from the "active electrical energy" (Wh or kWh) register of the (most authoritative) electrical meter measuring energy exported (to the grid).
	EnergyActiveImportRegister   = "Energy.Active.Import.Register"   // Numerical value read from the "active electrical energy" (Wh or kWh) register of the (most authoritative) electrical meter measuring energy imported (from the grid supply).
	EnergyReactiveExportRegister = "Energy.Reactive.Export.Register" // Numerical value read from the "reactive electrical energy" (VARh or kVARh) register of the (most authoritative) electrical meter measuring energy exported (to the grid).
	EnergyReactiveImportRegister = "Energy.Reactive.Import.Register" // Numerical value read from the "reactive electrical energy" (VARh or kVARh) register of the (most authoritative) electrical meter measuring energy imported (from the grid supply).
	EnergyActiveExportInterval   = "Energy.Active.Export.Interval"   // Absolute amount of "active electrical energy" (Wh or kWh) exported (to the grid) during an associated time "interval", specified by a Metervalues ReadingContext, and applicable interval duration configuration values (in seconds) for "ClockAlignedDataInterval" and "MeterValueSampleInterval".
	EnergyActiveImportInterval   = "Energy.Active.Import.Interval"   // Absolute amount of "active electrical energy" (Wh or kWh) imported (from the grid supply) during an associated time "interval", specified by a Metervalues ReadingContext, and applicable interval duration configuration values (in seconds) for "ClockAlignedDataInterval" and "MeterValueSampleInterval".
	Voltage                      = "Voltage"
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

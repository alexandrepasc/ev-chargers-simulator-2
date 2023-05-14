package assets

const (
	Info  Severity = "info"
	Error Severity = "error"
	Warn  Severity = "warn"
	Panic Severity = "panic"
	Fatal Severity = "fatal"
)

var (
	Current                    = "Current"
	Power                      = "Power"
	EnergyActiveExportRegister = "Energy.Active.Export.Register"
	EnergyActiveImportRegister = "Energy.Active.Import.Register"
)

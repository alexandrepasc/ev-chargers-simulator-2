package assets

const (
	Info  Severity = "info"
	Error Severity = "error"
	Warn  Severity = "warn"
	Panic Severity = "panic"
	Fatal Severity = "fatal"

	Request  string = "request"
	Response string = "response"
)

var (
	Current                    = "Current"
	Power                      = "Power"
	EnergyActiveExportRegister = "Energy.Active.Export.Register"
	EnergyActiveImportRegister = "Energy.Active.Import.Register"
)

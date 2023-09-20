package modbus

const (
	// bit16 int = 16
	addr2816 uint16 = 2816
	addr3590 uint16 = 3590
)

/*
Map the key word from the configuration with the function related in the service file.
*/
var functionMap = map[string]interface{}{
	"powerImport":   func(m *Modbus) int { return m.powerImportW() },
	"powerImportKw": func(m *Modbus) int { return m.powerImportKw() },
}

package modbus

const (
	// bit16 int = 16
	conversion float64 = 1000
	sqrt       float64 = 3
)

/*
Map the key word from the configuration with the function related in the service file.
*/
var functionMap = map[string]interface{}{
	"powerImportW":      func(m *Modbus) int { return m.powerImportW() },
	"powerImportKw":     func(m *Modbus) int { return m.powerImportKw() },
	"importVa":          func(m *Modbus) int { return m.importVa() },
	"importKvA":         func(m *Modbus) int { return m.importKvA() },
	"totalCurrentA":     func(m *Modbus) int { return m.totalCurrentA() },
	"totalCurrentKa":    func(m *Modbus) int { return m.totalCurrentKa() },
	"currentPerPhaseA":  func(m *Modbus) int { return m.currentPerPhaseA() },
	"currentPerPhaseKa": func(m *Modbus) int { return m.currentPerPhaseKa() },
}

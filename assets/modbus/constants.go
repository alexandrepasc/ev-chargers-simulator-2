package modbus

const bit16 int = 16

var functionMap = map[string]interface{}{
	"powerImport": func(m *Modbus) int { return m.powerImport() },
}

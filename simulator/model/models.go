package model

import "github.com/google/uuid"

type Struct struct {
	ID        uuid.UUID   `json:"id" validate:"required"`   // Model identifier
	Name      string      `json:"name" validate:"required"` // Model name
	Type      Type        `json:"type" validate:"required"` // Type of asset that this config can be used (ocpp, modbus)
	CA        string      `json:"ca,omitempty"`             // CA certificate location and name, full path
	Cert      string      `json:"cert,omitempty"`           // Client certificate location and name, full path
	Key       string      `json:"key,omitempty"`            // Client certificate key location and name, full path
	BasicAuth BasicAuth   `json:"basicAuth,omitempty"`      // HTTP basic authentication credentials
	Ocpp      OcppModel   `json:"ocpp,omitempty"`           // Ocpp structure
	Modbus    ModbusModel `json:"modbus,omitempty"`         // Modbus structure
}

type BasicAuth struct {
	Username string `json:"username"` // HTTP basic authentication username
	Password string `json:"password"` // HTTP basic authentication password
}

type OcppModel struct {
	SerialNumb      string    `json:"serialNumb"`      // Equipment serial number
	Model           string    `json:"model"`           // Equipment model name
	Vendor          string    `json:"vendor"`          // Equipment vendor
	VendorID        string    `json:"vendorId"`        // Vendor identifier
	FwVersion       string    `json:"fwVersion"`       // Firmware version
	MeterSerialNumb string    `json:"meterSerialNumb"` // Power meter serial number
	Modem           OcppModem `json:"modem"`           // Modem information
}

type OcppModem struct {
	Iccid string `json:"iccid"` // SIM card identifier
	Imsi  string `json:"imsi"`  // International Mobile Subscriber Identity
}

// TODO: to be done
type ModbusModel struct {
	Coils            Coils            `json:"coils,omitempty"`            // Coils mapping
	Discrete         Discrete         `json:"discrete,omitempty"`         // Descrete inputs mapping
	HoldingRegisters HoldingRegisters `json:"holdingRegisters,omitempty"` // Holding registers mapping
	InputRegisters   InputRegisters   `json:"inputRegisters,omitempty"`   // Input registers mapping
}

type Coils struct{}

type Discrete struct{}

type HoldingRegisters struct {
	Addresses map[int]string `json:"addresses"` // List of modbus addresses mapping
}

type InputRegisters struct{}

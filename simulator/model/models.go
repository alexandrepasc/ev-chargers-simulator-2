package model

import "github.com/google/uuid"

type Struct struct {
	ID     uuid.UUID   `json:"id" validate:"required"`   // Model identifier
	Name   string      `json:"name" validate:"required"` // Model name
	Type   Type        `json:"type" validate:"required"` // Type of asset that this config can be used (ocpp, modbus)
	Ocpp   OcppModel   `json:"ocpp"`                     // Ocpp structure
	Modbus ModbusModel `json:"modbus"`                   // Modbus structure
}

type OcppModel struct {
	SerialNumb      string    `json:"serialNumb"`      // Equipment serial number
	Model           string    `json:"model"`           // Equipment model name
	Vendor          string    `json:"vendor"`          // Equipment vendor
	VendorID        string    `json:"vendorId"`        // Vendor identifier
	FwVersion       string    `json:"fwVersion"`       // Firmware version
	Modem           OcppModem `json:"modem"`           // Modem information
	MeterSerialNumb string    `json:"meterSerialNumb"` // Power meter serial number
}

type OcppModem struct {
	Iccid string `json:"iccid"` // SIM card identifier
	Imsi  string `json:"imsi"`  // International Mobile Subscriber Identity
}

// TODO: to be done
type ModbusModel struct{}

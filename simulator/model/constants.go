package model

const (
	Ocpp   Type   = "ocpp"
	Modbus Type   = "modbus"
	jsonEx string = ".json" // Json file extension
)

var DefOcppMod = Struct{
	Name: "default",
	Type: Ocpp,
	Ocpp: OcppModel{
		SerialNumb:      "default-serial",
		Model:           "defModel",
		Vendor:          "defVendor",
		VendorID:        "asd-asd-asd-asd",
		FwVersion:       "0.0.0.0",
		MeterSerialNumb: "0.0.0.1",
		Modem: OcppModem{
			Iccid: "99999999999",
			Imsi:  "88888888888",
		},
	},
}

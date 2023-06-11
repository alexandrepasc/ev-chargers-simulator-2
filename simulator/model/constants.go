package model

const (
	Ocpp   Type = "ocpp"
	Modbus Type = "modbus"
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
		AuthorizeRemote: false,
		Modem: OcppModem{
			Iccid: "99999999999",
			Imsi:  "88888888888",
		},
	},
}

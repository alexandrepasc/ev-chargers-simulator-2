package model

type OcppModel struct {
	SerialNumb string    `json:"serialNumb"` // Equipment serial number
	Model      string    `json:"model"`      // Equipment model name
	Vendor     string    `json:"vendor"`     // Equipment vendor
	FwVersion  string    `json:"fwVersion"`  // Firmware version
	Modem      OcppModem `json:"modem"`      // Modem information
}

type OcppModem struct {
	Iccid string `json:"iccid"` // SIM card identifier
	Imsi  string `json:"imsi"`  // International Mobile Subscriber Identity
}

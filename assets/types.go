package assets

type Severity string

type ConnectorStatus int64

type ResetType string

type Availability string

type AssetStatus string

type ConnectType string

type BootSeq struct {
	IsToTrigger  bool        // To trigger the boot up sequence
	BootStatus   interface{} // The booting status of the cp
	BootReason   interface{} // Boot reason
	BootInterval int         // Handle the boot interval when the boot fails
}

type ResetSeq struct {
	IsToTrigger bool // To trigger the reset logic sequence
	EvseIndex   *int // The evse index to be reset
}

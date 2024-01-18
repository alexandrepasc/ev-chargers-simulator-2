package assets

import "github.com/google/uuid"

type DataInfo struct {
	UUID        uuid.UUID   // Identifier of the asset
	Name        string      // Name of the asset
	Status      AssetStatus // Status of the asset
	Power       float64     // Charge point consumption power
	PowerFactor int64       // Power factor
	Energy      float64     // Charge point energy
	Current     float64     // Charge point current
}

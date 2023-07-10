package handler

import "github.com/google/uuid"

// type sim struct {
// 	asset *simulator.Asset // Assets list
// 	// model model.OcppModel
// }

// type Channel struct {
// 	name string    // Name of the asset
// 	uuid uuid.UUID // Identifier of the asset
// }

type info struct {
	UUID   uuid.UUID   // Identifier of the asset
	Name   string      // Name of the asset
	Status assetStatus // Status of the asset
	Power  float64     // Charge point consumption power
	Energy float64     // Charge point energy
}

type Status struct {
	Total  int64    `json:"total"`
	Assets []Assets `json:"assets"`
}

type Assets struct {
	ID     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	State  string    `json:"state"`
	Power  float64   `json:"power"`
	Energy float64   `json:"energy"`
}

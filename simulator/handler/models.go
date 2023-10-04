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

type Status struct {
	Total  int64    `json:"total"`  // Total number of assets
	Assets []Assets `json:"assets"` // List of assets
}

type Assets struct {
	ID     uuid.UUID `json:"id"`     // Asset identifier
	Name   string    `json:"name"`   // Asset name
	State  string    `json:"state"`  // State of the asset (acive, inactive)
	Power  float64   `json:"power"`  // Current asset power
	Energy float64   `json:"energy"` // Energy consumption of the execution
}

package handler

import (
	"github.com/alexandrepasc/ev-chargers-simulator-2/simulator"
	"github.com/google/uuid"
)

type sim struct {
	asset *simulator.Asset // Assets list
	// model model.OcppModel
}

type channel struct {
	name string    // Name of the asset
	uuid uuid.UUID // Identifier of the asset
}

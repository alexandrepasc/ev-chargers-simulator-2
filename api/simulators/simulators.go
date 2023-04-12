package simulators

import (
	"net/http"

	"github.com/alexandrepasc/ev-chargers-simulator-2/simulator"
	"github.com/gin-gonic/gin"
)

type Simulators struct {
	Al []simulator.Asset
}

/*
Creates the simulators endpoint controller handlers

r	-	api engine (*gin.Engine)
*/
func (s Simulators) Simulators(r *gin.Engine) {
	r.GET(simulatorsEp, s.getSimulators)
}

/*
Sets the get simulators endpoint controller

c	-	request  context (*gin.context)
*/
func (s Simulators) getSimulators(c *gin.Context) {
	resp := GetSimulatorsTemp{
		Total:  int64(len(s.Al)),
		Assets: s.Al,
	}

	c.IndentedJSON(http.StatusOK, resp)
}

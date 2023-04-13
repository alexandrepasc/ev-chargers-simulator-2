package simulators

import (
	"net/http"

	"github.com/alexandrepasc/ev-chargers-simulator-2/simulator"
	"github.com/gin-gonic/gin"
)

type Simulators struct {
	Sim simulator.Simulator // Simulator package
}

/*
Creates the simulators endpoint controller handlers

r	-	api engine (*gin.Engine)
*/
func (s *Simulators) Simulators(r *gin.Engine) {
	r.GET(simulatorsEp, s.getSimulators)
	r.GET(simModelsEp, s.getModels)
}

/*
Sets the get simulators endpoint controller

c	-	request  context (*gin.context)
*/
func (s *Simulators) getSimulators(c *gin.Context) {
	s.Sim.GetSimsConfs()

	resp := GetSimulatorsTemp{
		Total:  int64(len(s.Sim.Al)),
		Assets: s.Sim.Al,
	}

	c.IndentedJSON(http.StatusOK, resp)
}

/*
Sets the get models endpoint controller

c	-	request  context (*gin.context)
*/
func (s *Simulators) getModels(c *gin.Context) {
	resp := GetModelsTemp{
		Total:  int64(len(s.Sim.Oml)),
		Models: s.Sim.Oml,
	}

	c.IndentedJSON(http.StatusOK, resp)
}

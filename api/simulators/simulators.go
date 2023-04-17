package simulators

import (
	"net/http"

	"github.com/alexandrepasc/ev-chargers-simulator-2/api/errors"
	"github.com/alexandrepasc/ev-chargers-simulator-2/common"
	"github.com/alexandrepasc/ev-chargers-simulator-2/simulator"
	"github.com/alexandrepasc/ev-chargers-simulator-2/translation"
	"github.com/alexandrepasc/ev-chargers-simulator-2/translation/text"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Simulators struct {
	Sim  simulator.Simulator     // Simulator package
	Lang translation.Translation // Translation language setting
}

/*
Creates the simulators endpoint controller handlers

r	-	api engine (*gin.Engine)
*/
func (s *Simulators) Simulators(r *gin.Engine) {
	r.GET(simulatorsEp, s.getSimulators)
	r.POST(simulatorsEp, s.postSimulators)
	r.PUT(simulatorsIDEp, s.putSimulators)
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
Sets the post simulators endpoint controller, it will generate a new sim configuration file.

c	-	request  context (*gin.context)
*/
func (s *Simulators) postSimulators(c *gin.Context) {
	var b = simulator.Asset{}

	err := c.BindJSON(&b)

	if err != nil {
		common.Log("postSimulators").Error(err)

		var r = errors.ErroMsg{
			Message: s.Lang.Get(text.RequestBodyDoesntMatch),
		}

		ok := errors.StructValidate(r, c, s.Lang)

		if ok {
			c.IndentedJSON(http.StatusBadRequest, r)
		}

		return
	}

	ok, msg, r := s.Sim.CreateSimConf(&b)

	if ok {
		c.JSON(http.StatusCreated, r)
	} else {
		common.Log("postSimulators").Error(msg)

		var r = errors.ErroMsg{
			Message: msg,
		}
		v := errors.StructValidate(r, c, s.Lang)
		if v {
			c.IndentedJSON(http.StatusBadRequest, r)
		}
	}
}

/*
Sets the put simulators endpoint controller, it will replace the current configurations with what
was sent by the user.

c	-	request  context (*gin.context)
*/
func (s *Simulators) putSimulators(c *gin.Context) {
	var id, errP = uuid.Parse(c.Param("id"))

	if errP != nil {
		common.Log("putSimulators").Error(errP)

		var r = errors.ErroMsg{
			Message: s.Lang.Get(text.UUIDParsingError),
		}

		v := errors.StructValidate(r, c, s.Lang)

		if v {
			c.IndentedJSON(http.StatusBadRequest, r)
		}

		return
	}

	var req = simulator.Asset{}

	err := c.BindJSON(&req)

	if err != nil {
		common.Log("putSimulators").Error(err)

		var r = errors.ErroMsg{
			Message: s.Lang.Get(text.RequestBodyDoesntMatch),
		}

		v := errors.StructValidate(r, c, s.Lang)

		if v {
			c.IndentedJSON(http.StatusBadRequest, r)
		}

		return
	}

	ok, msg, code, r := s.Sim.UpdateSimConf(id, &req)

	if ok {
		c.IndentedJSON(http.StatusOK, r)
	} else {
		common.Log("putSimulators").Error(msg)

		var r = errors.ErroMsg{
			Message: msg,
		}

		v := errors.StructValidate(r, c, s.Lang)

		if v {
			c.IndentedJSON(code, r)
		}
	}
}

/*
Sets the get models endpoint controller

c	-	request  context (*gin.context)
*/
func (s *Simulators) getModels(c *gin.Context) {
	s.Sim.GetSimsConfs()

	resp := GetModelsTemp{
		Total:  int64(len(s.Sim.Oml)),
		Models: s.Sim.Oml,
	}

	c.IndentedJSON(http.StatusOK, resp)
}

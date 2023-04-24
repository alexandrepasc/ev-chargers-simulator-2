//nolint:dupl // because it needs to be reviewed
package simulators

import (
	"net/http"

	"github.com/alexandrepasc/ev-chargers-simulator-2/api/errors"
	"github.com/alexandrepasc/ev-chargers-simulator-2/common"
	"github.com/alexandrepasc/ev-chargers-simulator-2/simulator"
	"github.com/alexandrepasc/ev-chargers-simulator-2/simulator/handler"
	"github.com/alexandrepasc/ev-chargers-simulator-2/simulator/model"
	"github.com/alexandrepasc/ev-chargers-simulator-2/translation"
	"github.com/alexandrepasc/ev-chargers-simulator-2/translation/text"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Simulators struct {
	Sim  simulator.Simulator     // Simulator package
	Mod  model.Model             // Model package
	Lang translation.Translation // Translation language setting
	H    handler.Handler
}

/*
Creates the simulators endpoint controller handlers

r	-	api engine (*gin.Engine)
*/
func (s *Simulators) Simulators(r *gin.Engine) {
	s.H = handler.Handler{
		L: s.Lang,
	}

	r.GET(simulatorsEp, s.getSimulators)
	r.POST(simulatorsEp, s.postSimulators)
	r.PUT(simulatorsIDEp, s.putSimulators)
	r.DELETE(simulatorsIDEp, s.deleteSimulators)

	r.GET(simModelsEp, s.getModels)
	r.POST(simModelsEp, s.postModels)
	r.PUT(simModelsIDEp, s.putModels)

	r.POST(simulatorsEp+"/run", s.postRun)
	r.POST(simulatorsEp+"/stop", s.postStop)
}

/*
Sets the get simulators endpoint controller

c	-	request  context (*gin.context)
*/
func (s *Simulators) getSimulators(c *gin.Context) {
	s.Sim.GetSimulators()

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
Sets the delete simulators endpoint controller, it will enable the user to delete one simulator
configuration file.

c	-	request  context (*gin.context)
*/
func (s *Simulators) deleteSimulators(c *gin.Context) {
	var id, errP = uuid.Parse(c.Param("id"))

	if errP != nil {
		common.Log("deleteSimulators").Error(errP)

		var r = errors.ErroMsg{
			Message: s.Lang.Get(text.UUIDParsingError),
		}

		v := errors.StructValidate(r, c, s.Lang)

		if v {
			c.IndentedJSON(http.StatusBadRequest, r)
		}

		return
	}

	ok, msg, cod := s.Sim.DeleteSimConf(id)

	if ok {
		c.IndentedJSON(cod, http.NoBody)
	} else {
		common.Log("deleteSimulators").Error(msg)

		var r = errors.ErroMsg{
			Message: msg,
		}

		v := errors.StructValidate(r, c, s.Lang)

		if v {
			c.IndentedJSON(cod, r)
		}
	}
}

/*
Sets the get models endpoint controller

c	-	request  context (*gin.context)
*/
func (s *Simulators) getModels(c *gin.Context) {
	s.Mod.GetModels()

	resp := GetModelsTemp{
		Total:  int64(len(s.Mod.Ml)),
		Models: s.Mod.Ml,
	}

	c.IndentedJSON(http.StatusOK, resp)
}

/*
Sets the post models endpoint controller, it will generate a new model configuration file.

c	-	request  context (*gin.context)
*/
func (s *Simulators) postModels(c *gin.Context) {
	var b = model.Struct{}

	err := c.BindJSON(&b)

	if err != nil {
		common.Log("postModels").Error(err)

		var r = errors.ErroMsg{
			Message: s.Lang.Get(text.RequestBodyDoesntMatch),
		}

		ok := errors.StructValidate(r, c, s.Lang)

		if ok {
			c.IndentedJSON(http.StatusBadRequest, r)
		}

		return
	}

	ok, msg, cod, r := s.Mod.CreateModel(&b)

	if ok {
		c.IndentedJSON(cod, r)
	} else {
		common.Log("postModels").Error(msg)

		var r = errors.ErroMsg{
			Message: msg,
		}
		v := errors.StructValidate(r, c, s.Lang)
		if v {
			c.IndentedJSON(cod, r)
		}
	}
}

/*
Sets the put models endpoint controller, it will replace the current configurations with what
was sent by the user.

c	-	request  context (*gin.context)
*/
func (s *Simulators) putModels(c *gin.Context) {
	var id, errP = uuid.Parse(c.Param("id"))

	if errP != nil {
		common.Log("putModels").Error(errP)

		var r = errors.ErroMsg{
			Message: s.Lang.Get(text.UUIDParsingError),
		}

		v := errors.StructValidate(r, c, s.Lang)

		if v {
			c.IndentedJSON(http.StatusBadRequest, r)
		}

		return
	}

	var req = model.Struct{}

	err := c.BindJSON(&req)

	if err != nil {
		common.Log("putModels").Error(err)

		var r = errors.ErroMsg{
			Message: s.Lang.Get(text.RequestBodyDoesntMatch),
		}

		v := errors.StructValidate(r, c, s.Lang)

		if v {
			c.IndentedJSON(http.StatusBadRequest, r)
		}

		return
	}

	ok, msg, cod, r := s.Mod.UpdateModel(id, &req)

	if ok {
		c.IndentedJSON(http.StatusOK, r)
	} else {
		common.Log("putModels").Error(msg)

		var r = errors.ErroMsg{
			Message: msg,
		}

		v := errors.StructValidate(r, c, s.Lang)

		if v {
			c.IndentedJSON(cod, r)
		}
	}
}

/*
Sets the post run endpoint controller

c	-	request  context (*gin.context)
*/
func (s *Simulators) postRun(c *gin.Context) {
	s.Sim.GetSimulators()

	s.H.Al = s.Sim.Al

	s.H.Quit = s.H.Start()

	c.IndentedJSON(http.StatusNoContent, http.NoBody)
}

/*
Sets the post stop endpoint controller

c	-	request  context (*gin.context)
*/
func (s *Simulators) postStop(c *gin.Context) {
	s.H.Stop()

	c.IndentedJSON(http.StatusNoContent, http.NoBody)
}

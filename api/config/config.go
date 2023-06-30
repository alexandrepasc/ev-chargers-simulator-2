package config

import (
	"net/http"

	"github.com/alexandrepasc/ev-chargers-simulator-2/api/errors"
	"github.com/alexandrepasc/ev-chargers-simulator-2/common"
	"github.com/alexandrepasc/ev-chargers-simulator-2/settings/general"
	"github.com/alexandrepasc/ev-chargers-simulator-2/translation"
	"github.com/alexandrepasc/ev-chargers-simulator-2/translation/text"
	"github.com/gin-gonic/gin"
)

type Configs struct {
	Lang translation.Translation // Translation language setting
	Path string                  // Path to the general configuration folder
}

/*
Creates the configs endpoint controler

r	-	api engine (*gin.Engine)
*/
func (cfg *Configs) Configs(r *gin.Engine) {
	r.GET(generalEp, cfg.getGeneral)
	r.PUT(generalEp, cfg.putGeneral)

	r.GET(simulatorsEp, getSimulators)
}

/*
Sets the get general endpoint controller

c	-	request  context (*gin.context)
*/
func (cfg *Configs) getGeneral(c *gin.Context) {
	var ok, msg, code, r = general.GetGeneralConf(cfg.Path, cfg.Lang)

	if !ok {
		common.Log("getGeneral").Error(msg)

		var e = errors.ErroMsg{
			Message: msg,
		}

		if errors.StructValidate(e, c, cfg.Lang) {
			c.IndentedJSON(code, e)
		}

		return
	}

	c.IndentedJSON(code, r)
}

/*
Sets the put general endpoint controller, it will replace the current configurations with what
was sent by the user.

c	-	request  context (*gin.context)
*/
func (cfg *Configs) putGeneral(c *gin.Context) {
	var req = general.Model{}

	var err = c.BindJSON(&req)
	if err != nil {
		common.Log("putGeneral").Error(err)

		var e = errors.ErroMsg{
			Message: cfg.Lang.Get(text.RequestBodyDoesntMatch),
		}

		if errors.StructValidate(e, c, cfg.Lang) {
			c.IndentedJSON(http.StatusBadRequest, e)
		}

		return
	}

	var ok, msg, code, r = general.UpdateGeneralConf(&req, cfg.Path, cfg.Lang)

	if !ok {
		common.Log("putGeneral").Error(msg)

		var e = errors.ErroMsg{
			Message: msg,
		}

		if errors.StructValidate(e, c, cfg.Lang) {
			c.IndentedJSON(code, e)
		}

		return
	}

	c.IndentedJSON(code, r)
}

/*
Sets the get simulators endpoint controller

c	-	request  context (*gin.context)
*/
func getSimulators(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}

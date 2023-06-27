package config

import (
	"net/http"

	"github.com/alexandrepasc/ev-chargers-simulator-2/translation"
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

	r.GET(simulatorsEp, getSimulators)
}

/*
Sets the get general endpoint controller

c	-	request  context (*gin.context)
*/
func (cfg *Configs) getGeneral(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}

/*
Sets the get simulators endpoint controller

c	-	request  context (*gin.context)
*/
func getSimulators(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}

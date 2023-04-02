package configs

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
Creates the configs endpoint controler

r	-	api engine (*gin.Engine)
*/
func Configs(r *gin.Engine) {
	r.GET(generalEp, getGeneral)

	r.GET(simulatorsEp, getSimulators)
}

/*
Sets the get general endpoint controller

c	-	request  context (*gin.context)
*/
func getGeneral(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}

/*
Sets the get simulators endpoint controller

c	-	request  context (*gin.context)
*/
func getSimulators(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}

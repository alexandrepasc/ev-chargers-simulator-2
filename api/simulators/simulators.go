package simulators

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
Creates the simulators endpoint controller handlers

r	-	api engine (*gin.Engine)
*/
func Simulators(r *gin.Engine) {
	r.GET(simulatorsEp, getSimulators)
}

/*
Sets the get simulators endpoint controller

c	-	request  context (*gin.context)
*/
func getSimulators(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}

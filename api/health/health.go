package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
Creates the health endpoint controler

router	-	api engine (*gin.Engine)
*/
func Health(router *gin.Engine) {
	router.GET("/health", getHealth)
}

/*
Sets the get health method handler

context	-	request  context (*gin.context)
*/
func getHealth(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{"message": "OK"})
}

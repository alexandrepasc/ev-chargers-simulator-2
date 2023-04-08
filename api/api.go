package api

import (
	"net/http"

	"github.com/alexandrepasc/ev-chargers-simulator-2/api/config"
	"github.com/alexandrepasc/ev-chargers-simulator-2/api/health"
	"github.com/alexandrepasc/ev-chargers-simulator-2/api/simulators"
	"github.com/alexandrepasc/ev-chargers-simulator-2/common"
	"github.com/gin-gonic/gin"
)

/*
Build the api instance configurations, sets the endpoints, and returns the router instance

Returns the engine instance (*gin.Engine)
*/
func New() (router *gin.Engine) {
	router = buildRouter()

	health.Health(router)
	config.Configs(router)
	simulators.Simulators(router)

	return
}

/*
Runs the http server for the api instace

router	-	framework engine instance (*gin.Engine)
*/
func Serve(router *gin.Engine) {
	var err = router.Run("localhost:8000")

	if err != nil {
		common.APILog("Serve").Error(err)
	}
}

/*
Creates the engine instance

Returns the instance (*gin.Engine)
*/
func buildRouter() *gin.Engine {
	eng := gin.Default()

	eng.HandleMethodNotAllowed = true

	eng.NoMethod(func(ctx *gin.Context) { methodNotAllowed(ctx) })
	eng.NoRoute(func(ctx *gin.Context) { routeNotFound(ctx) })

	return eng
}

/*
Handles the method not allowed response

c	-	Gin context (*gin.Context)
*/
func methodNotAllowed(c *gin.Context) {
	c.JSON(http.StatusMethodNotAllowed, gin.H{"message": "Method not allowed."})
}

/*
Handles the route not found response

c	-	Gin context (*gin.Context)
*/
func routeNotFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{"message": "Route not found."})
}

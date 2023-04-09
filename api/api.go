package api

import (
	"net/http"

	"github.com/alexandrepasc/ev-chargers-simulator-2/api/config"
	"github.com/alexandrepasc/ev-chargers-simulator-2/api/health"
	"github.com/alexandrepasc/ev-chargers-simulator-2/api/simulators"
	"github.com/alexandrepasc/ev-chargers-simulator-2/common"
	"github.com/alexandrepasc/ev-chargers-simulator-2/translation"
	"github.com/alexandrepasc/ev-chargers-simulator-2/translation/text"
	"github.com/gin-gonic/gin"
)

/*
API structure to store the settings to the package.
*/
type API struct {
	Lang translation.Translation // Translation language setting
}

/*
Build the api instance configurations, sets the endpoints, and returns the router instance

Returns the engine instance (*gin.Engine)
*/
func (a API) New() (router *gin.Engine) {
	router = a.buildRouter()

	health.Health(router)
	config.Configs(router)
	simulators.Simulators(router)

	return
}

/*
Runs the http server for the api instace

router	-	framework engine instance (*gin.Engine)
*/
func (a API) Serve(router *gin.Engine) {
	var err = router.Run("localhost:8000")

	if err != nil {
		common.APILog("Serve").Error(err)
	}
}

/*
Creates the engine instance

Returns the instance (*gin.Engine)
*/
func (a API) buildRouter() *gin.Engine {
	eng := gin.Default()

	eng.HandleMethodNotAllowed = true

	eng.NoMethod(func(ctx *gin.Context) { a.methodNotAllowed(ctx) })
	eng.NoRoute(func(ctx *gin.Context) { a.routeNotFound(ctx) })

	return eng
}

/*
Handles the method not allowed response

c	-	Gin context (*gin.Context)
*/
func (a API) methodNotAllowed(c *gin.Context) {
	c.JSON(http.StatusMethodNotAllowed, gin.H{"message": a.Lang.Get(text.MethodNotAllowed)})
}

/*
Handles the route not found response

c	-	Gin context (*gin.Context)
*/
func (a API) routeNotFound(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{"message": a.Lang.Get(text.RouteNotFound)})
}

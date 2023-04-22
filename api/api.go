package api

import (
	"net/http"

	"github.com/alexandrepasc/ev-chargers-simulator-2/api/config"
	"github.com/alexandrepasc/ev-chargers-simulator-2/api/errors"
	"github.com/alexandrepasc/ev-chargers-simulator-2/api/health"
	"github.com/alexandrepasc/ev-chargers-simulator-2/api/simulators"
	"github.com/alexandrepasc/ev-chargers-simulator-2/common"
	"github.com/alexandrepasc/ev-chargers-simulator-2/settings/general"
	"github.com/alexandrepasc/ev-chargers-simulator-2/simulator"
	"github.com/alexandrepasc/ev-chargers-simulator-2/simulator/model"
	"github.com/alexandrepasc/ev-chargers-simulator-2/translation"
	"github.com/alexandrepasc/ev-chargers-simulator-2/translation/text"
	"github.com/gin-gonic/gin"
)

/*
API structure to store the settings to the package.
*/
type API struct {
	Lang    translation.Translation // Translation language setting
	General general.Model           // General configurations model
	Sim     simulator.Simulator     // Simulator package
	Mod     model.Model             // Model package
}

/*
Build the api instance configurations, sets the endpoints, and returns the router instance

Returns the engine instance (*gin.Engine)
*/
func (a *API) New() (r *gin.Engine) {
	r = a.buildRouter()

	health.Health(r)
	config.Configs(r)

	s := simulators.Simulators{
		Sim:  a.Sim,  // Simulator package
		Mod:  a.Mod,  // Model package
		Lang: a.Lang, // Translation language setting
	}
	s.Simulators(r)

	return
}

/*
Runs the http server for the api instace

router	-	framework engine instance (*gin.Engine)
*/
func (a *API) Serve(router *gin.Engine) {
	var err = router.Run("localhost:8000")

	if err != nil {
		common.APILog("Serve").Error(err)
	}
}

/*
Creates the engine instance

Returns the instance (*gin.Engine)
*/
func (a *API) buildRouter() *gin.Engine {
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
func (a *API) methodNotAllowed(c *gin.Context) {
	var r = errors.ErroMsg{
		Message: a.Lang.Get(text.MethodNotAllowed),
	}

	ok := errors.StructValidate(r, c, a.Lang)

	if ok {
		c.IndentedJSON(http.StatusMethodNotAllowed, r)
	}
}

/*
Handles the route not found response

c	-	Gin context (*gin.Context)
*/
func (a *API) routeNotFound(c *gin.Context) {
	var r = errors.ErroMsg{
		Message: a.Lang.Get(text.RouteNotFound),
	}

	ok := errors.StructValidate(r, c, a.Lang)

	if ok {
		c.IndentedJSON(http.StatusNotFound, r)
	}
}

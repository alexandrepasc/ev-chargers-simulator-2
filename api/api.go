package api

import (
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
	router := gin.Default()

	return router
}

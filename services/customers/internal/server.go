package internal

import (
	"fmt"
	"sync"

	"github.com/gin-gonic/gin"
)


type Server struct {
	EnvVariables EnvironmentVariables
}

func getRouter(routes []Route, globMiddlewares []gin.HandlerFunc) *gin.Engine {
	var router = gin.Default()
	var v1 = router.Group("/api/v1")

	for _, globMiddleware := range globMiddlewares {
		v1.Use(globMiddleware)
	}

	for _, route := range routes {
		v1.Handle(route.HttpMethod, route.Path, route.HandlerFunc)
	}

	return router
}

func LaunchServer(routes []Route, globMiddlewares []gin.HandlerFunc, env EnvironmentVariables, wg *sync.WaitGroup){
	defer wg.Done()

	var router = getRouter(routes, globMiddlewares)

	router.Run(fmt.Sprintf(":%s", env.Port))
}
package internal

import (
	"fmt"
	"log"
	"sync"

	"github.com/gin-gonic/gin"
)

type Route struct {
	Method string
	Path string
	HandlerFunc gin.HandlerFunc
}

type Server struct {
	Router *gin.Engine
}

func getRouter(routes []Route, globalMiddlewares []gin.HandlerFunc, env EnvVariables) *gin.Engine {
	var router = gin.Default()
	var v1 = router.Group("/api/v1")

	for _, globalMiddleware := range globalMiddlewares {
		v1.Use(globalMiddleware)
	}

	for _, route := range routes {
		v1.Handle(route.Method, route.Path, route.HandlerFunc)
	}

	return router
}

func LaunchServer(routes []Route, globalMiddlewares []gin.HandlerFunc, env EnvVariables, wg *sync.WaitGroup){
	defer wg.Done()

	var router = getRouter(routes, globalMiddlewares, env)
	log.Fatalln(router.Run(fmt.Sprintf(":%s", env.Port)))
}

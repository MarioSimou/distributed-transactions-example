package internal

import (
	"github.com/gin-gonic/gin"
)

type Route struct {
	Method string
	Path string
	HandlerFunc gin.HandlerFunc
}

func GetRouter(routes []Route) *gin.Engine {
	var router = gin.Default()
	var v1 = router.Group("/api/v1")
	
	for _, route := range routes {
		v1.Handle(route.Method, route.Path, route.HandlerFunc)
	}

	return router
}


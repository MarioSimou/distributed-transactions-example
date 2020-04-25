package internal

import "github.com/gin-gonic/gin"


type Server struct {
	EnvVariables EnvironmentVariables
}

func (s *Server) Setup(routes []Route) *gin.Engine {
	var router = gin.Default()
	var v1 = router.Group("/api/v1")

	for _, route := range routes {
		v1.Handle(route.HttpMethod, route.Path, route.HandlerFunc)
	}
	return router
}
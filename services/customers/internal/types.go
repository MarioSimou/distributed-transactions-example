package internal

import (
	"github.com/gin-gonic/gin"
)

type EnvironmentVariables struct {
	Port string
}

type Route struct {
	HttpMethod string
	Path string
	HandlerFunc gin.HandlerFunc
}
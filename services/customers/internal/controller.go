package internal

import (
	"github.com/gin-gonic/gin"
)

type Controller struct {
	EnvVariables EnvironmentVariables
}

func (contr *Controller) Ping(c *gin.Context){
	c.JSON(200, gin.H{"message": "pong"})
}
package internal

import (
	"time"

	"github.com/gin-gonic/gin"
)

type EnvironmentVariables struct {
	Port string
	UIDomain string
}

type Route struct {
	HttpMethod string
	Path string
	HandlerFunc gin.HandlerFunc
}

type userId struct {
	Id int64 `uri:"id" binding:"required"`
} 

type signInBody struct {
	Email string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type postUserBody struct {
	ID        int32 `sql:"primary_key"`
	Username  string `json:"username" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	Balance   float64 `json:"balance" binding:"gte=0"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
type updateUserBody struct {
	ID        int32 `sql:"primary_key"`
	Username  string `json:"username" binding:"required_without_all=Email Password Balance"`
	Email     string `json:"email" binding:"omitempty,email"`
	Password  string `json:"password" binding:"omitempty,min=8"`
	Balance   *float64 `json:"balance" binding:"omitempty,gte=0"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`
}

type Response struct {
	Status int `json:"status"`
	Success bool `json:"success"`
	Message string `json:"message"`
	Data interface{} `json:"data"`
}

package internal

import (
	"fmt"
	"log"
	"os"
	"products/internal/models/products/public/model"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
)

type postProductBody struct {
	ID        int32  `json:"id" sql:"primary_key"` 
	Name      string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	Price     float64 `json:"price" binding:"required,gt=0"`
	Quantity  *int32 `json:"quantity" binding:"required,gte=0"`
	Currency  string `json:"currency" binding:"required,oneof=GBP EURO USD"`
	Image string `json:"image" binding:"required"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type updateProductBody struct {
	ID        int32 `json:"id" sql:"primary_key"`
	Name      string `json:"name" binding:"required_without_all=Description Price Quantity Currency"`
	Description string `json:"description" binding:"omitempty,max=500"`
	Price     float64 `json:"price" binding:"omitempty,gt=0"`
	Quantity  *int32 `json:"quantity" binding:"omitempty,gte=0"`
	Currency  model.Currency `json:"currency" binding:"omitempty,oneof=GBP USD EURO"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type productUri struct {
	Id int64 `json:"id" uri:"id" binding:"required"`
}

type postOrderBody struct {
	UID string `json:"uid" binding:"required"`
	ProductID int64 `json: "productId" binding:"required,gt=0"`
	Quantity  int64 `json:"quantity" binding:"required,gt=0"`
	UserID    int64 `json:"userId" binding:"required,gt=0"`
}

type orderUri struct {
	Id int64 `json:"id" uri:"id" binding:"required"`
}


type response struct {
	Status int `json:"status"`
	Success bool `json:"success"`
	Message string `json:"message"`
	Data interface{} `json:"data"`
}	

func copy(src interface{}, tgt interface{}) error {
	var tv = reflect.ValueOf(tgt)
	if tv.Kind() != reflect.Ptr || tv.Elem().Kind() != reflect.Struct {
		return  fmt.Errorf("Error: 'Target is not a pointer of struct'\n")
	}

	var sv = reflect.ValueOf(src)
	var te = tv.Elem()
	for i := 0; i < sv.NumField(); i++ {
		var sf = sv.Field(i)
		var tf = te.Field(i)

		if reflect.Zero(tf.Type()).Interface() == tf.Interface() {
			tf.Set(sf)
		}
	}

	return nil
}

func HandleError(e error){
	if e != nil {
		log.Fatalf("Error: %v\n", e)
	}
}

func HandleCORS(c *gin.Context){
	c.Header("Access-Control-Allow-Origin", os.Getenv("ALLOW_ORIGIN_DOMAIN"))
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, HEAD")
	c.Header("Access-Control-Allow-Headers", "Content-Type")
	c.Header("Access-Control-Allow-Credentials", "true")
	c.Next()
}

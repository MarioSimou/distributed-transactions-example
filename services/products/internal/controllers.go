package internal

import (
	"fmt"
	"net/http"
	"products/internal/models/products/public/model"
	"reflect"
	"time"

	. "products/internal/models/products/public/table"

	. "github.com/go-jet/jet/postgres"

	"github.com/gin-gonic/gin"
	"github.com/go-jet/jet/qrm"
)

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

type EnvVariables struct {
	DBUri string
	Port string
}

type Controllers struct {
	Env EnvVariables
	DB qrm.DB
}

type postProductBody struct {
	ID        int32  `json:"id" sql:"primary_key"` 
	Name      string `json:"name" binding:"required"`
	Price     float64 `json:"price" binding:"required,gt=0"`
	Quantity  *int32 `json:"quantity" binding:"required,gte=0"`
	Currency  string `json:"currency" binding:"required,oneof=GBP EURO USD"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type updateProductBody struct {
	ID        int32 `json:"id" sql:"primary_key"`
	Name      string `json:"name" binding:"omitempty,required_without_all=Price Quantity Currency"`
	Price     float64 `json:"price" binding:"omitempty,gt=0"`
	Quantity  *int32 `json:"quantity" binding:"omitempty,gte=0"`
	Currency  model.Currency `json:"currency" binding="omitempty,oneof=GBP USD EURO"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type productUri struct {
	Id int64 `json:"id" uri:"id" binding:"required"`
}

type Response struct {
	Status int `json:"status"`
	Success bool `json:"success"`
	Message string `json:"message"`
	Data interface{} `json:"data"`
}	

func (contr *Controllers) Ping(c *gin.Context){
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}

func (contr *Controllers) GetProducts(c *gin.Context){
	var statement = Products.SELECT(Products.AllColumns).FROM(Products)
	var dest []model.Products
	if e := statement.Query(contr.DB, &dest); e != nil {
		if e == qrm.ErrNoRows {
			c.JSON(http.StatusNotFound, Response{Status: http.StatusNotFound, Success: false, Message: e.Error()})	
		} else {
			c.JSON(http.StatusBadRequest, Response{Status: http.StatusBadRequest, Success: false, Message: e.Error()})
		}
		return
	}
	
	c.JSON(http.StatusOK, Response{Status: http.StatusOK, Success: true, Data: dest})
}

func (contr *Controllers) GetProduct(c *gin.Context){
	var productUri productUri
	var dest model.Products

	if e := c.ShouldBindUri(&productUri); e != nil {
		c.JSON(http.StatusBadRequest, Response{Status: http.StatusBadRequest, Success: false, Message: e.Error()})	
	}

	var statement = Products.SELECT(Products.AllColumns).FROM(Products).WHERE(Products.ID.EQ(Int(productUri.Id)))
	if e := statement.Query(contr.DB, &dest); e != nil {
		if e == qrm.ErrNoRows {
			c.JSON(http.StatusNotFound, Response{Status: http.StatusNotFound, Success: false, Message: e.Error()})	
		} else {
			c.JSON(http.StatusBadRequest, Response{Status: http.StatusBadRequest, Success: false, Message: e.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, Response{Status: http.StatusOK, Success:true, Data: dest})
}

func (contr *Controllers) CreateProduct(c *gin.Context){
	var body postProductBody
	var dest model.Products
	if e := c.ShouldBindJSON(&body); e != nil {
		c.JSON(http.StatusBadRequest, Response{Status: http.StatusBadRequest, Success: false, Message: e.Error()})
		return
	}

	var statement = Products.INSERT(
		Products.Name,
		Products.Price,
		Products.Quantity,
		Products.Currency,
		Products.CreatedAt,
		Products.UpdatedAt,
	).VALUES(
		body.Name,
		body.Price,
		body.Quantity,
		body.Currency,
		time.Now(),
		time.Now(),
	).RETURNING(Products.AllColumns)
	fmt.Println(statement.Sql())

	if e := statement.Query(contr.DB, &dest); e != nil {
		c.JSON(http.StatusInternalServerError, Response{Status: http.StatusInternalServerError, Success: false, Message: e.Error()})
		return
	}

	c.Writer.Header().Add("Location", fmt.Sprintf("/api/v1/products/%d", dest.ID))
	c.JSON(http.StatusOK, Response{Status: http.StatusOK, Success:true, Data: dest})
}

func (contr *Controllers) UpdateProduct(c *gin.Context){
	var productUri productUri
	var existingProduct model.Products
	var dest model.Products
	var body updateProductBody

	if e := c.ShouldBindUri(&productUri); e != nil {
		c.JSON(http.StatusBadRequest, Response{Status: http.StatusBadRequest, Success: false, Message: e.Error()})	
		return
	}
	if e := c.ShouldBindJSON(&body); e != nil {
		c.JSON(http.StatusBadRequest, Response{Status: http.StatusBadRequest, Success: false, Message: e.Error()})
		return	
	}
	if e := Products.SELECT(Products.AllColumns).FROM(Products).WHERE(Products.ID.EQ(Int(productUri.Id))).Query(contr.DB, &existingProduct); e != nil {
		c.JSON(http.StatusNotFound, Response{Status: http.StatusNotFound, Success: false, Message: e.Error()})	
		return
	}
	if e := copy(existingProduct, &body); e != nil {
		c.JSON(http.StatusInternalServerError, Response{Status: http.StatusInternalServerError, Success: false, Message: e.Error()})	
		return
	}
	var statement = Products.UPDATE(Products.MutableColumns).MODEL(body).WHERE(Products.ID.EQ(Int(productUri.Id))).RETURNING(Products.AllColumns)
	if e := statement.Query(contr.DB, &dest); e != nil {
		c.JSON(http.StatusInternalServerError, Response{Status: http.StatusInternalServerError, Success: false, Message: e.Error()})	
		return
	} 

	c.JSON(http.StatusOK, Response{Status: http.StatusOK, Success:true, Data: dest})
}

func (contr *Controllers) DeleteProduct(c *gin.Context){
	var productUri productUri
	var product model.Products

	if e := c.ShouldBindUri(&productUri); e != nil {
		c.JSON(http.StatusBadRequest, Response{Status: http.StatusBadRequest, Success: false, Message: e.Error()})	
	}

	var statement = Products.DELETE().WHERE(Products.ID.EQ(Int(productUri.Id))).RETURNING(Products.ID)
	if e := statement.Query(contr.DB, &product); e != nil {
		if e == qrm.ErrNoRows {
			c.JSON(http.StatusNotFound, Response{Status: http.StatusNotFound, Success: false, Message: e.Error()})	
		} else {
			c.JSON(http.StatusBadRequest, Response{Status: http.StatusBadRequest, Success: false, Message: e.Error()})
		}
		return
	}

	c.JSON(http.StatusNoContent, Response{Status: http.StatusNoContent, Success:true})
}


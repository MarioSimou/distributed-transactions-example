package internal

import (
	"database/sql"
	"fmt"
	"net/http"
	"products/internal/models/products/public/model"
	"time"

	. "products/internal/models/products/public/table"
	"products/internal/rabbitmq"

	. "github.com/go-jet/jet/postgres"

	"github.com/gin-gonic/gin"
	"github.com/go-jet/jet/qrm"
)

type EnvVariables struct {
	DBUri string
	Port string
	QueueUri string
	QueuesNames []string
}

type Controller struct {
	Env EnvVariables
	DB *sql.DB
	Publisher rabbitmq.PublisherInterface
}

func (contr *Controller) Ping(c *gin.Context){
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}

func (contr *Controller) GetProducts(c *gin.Context){
	var statement = Products.SELECT(Products.AllColumns).FROM(Products)
	var dest []model.Products
	if e := statement.Query(contr.DB, &dest); e != nil {
		c.JSON(http.StatusInternalServerError, response{Status: http.StatusInternalServerError, Success: false, Message: e.Error()})
		return
	}
	if len(dest) == 0 {
		c.JSON(http.StatusNotFound, response{Status: http.StatusNotFound, Success: false, Message: qrm.ErrNoRows.Error()})	
		return
	}
	
	c.JSON(http.StatusOK, response{Status: http.StatusOK, Success: true, Data: dest})
}

func (contr *Controller) GetProduct(c *gin.Context){
	var productUri productUri
	var dest model.Products

	if e := c.ShouldBindUri(&productUri); e != nil {
		c.JSON(http.StatusBadRequest, response{Status: http.StatusBadRequest, Success: false, Message: e.Error()})	
		return
	}
	var statement = Products.SELECT(Products.AllColumns).FROM(Products).WHERE(Products.ID.EQ(Int(productUri.Id)))
	if e := statement.Query(contr.DB, &dest); e != nil {
		if e == qrm.ErrNoRows {
			c.JSON(http.StatusNotFound, response{Status: http.StatusNotFound, Success: false, Message: e.Error()})	
		} else {
			c.JSON(http.StatusInternalServerError, response{Status: http.StatusInternalServerError, Success: false, Message: e.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, response{Status: http.StatusOK, Success:true, Data: dest})
}

func (contr *Controller) CreateProduct(c *gin.Context){
	var body postProductBody
	var dest model.Products
	if e := c.ShouldBindJSON(&body); e != nil {
		c.JSON(http.StatusBadRequest, response{Status: http.StatusBadRequest, Success: false, Message: e.Error()})
		return
	}

	var statement = Products.INSERT(
		Products.Name,
		Products.Description,
		Products.Price,
		Products.Quantity,
		Products.Currency,
		Products.Image,
		Products.CreatedAt,
		Products.UpdatedAt,
	).VALUES(
		body.Name,
		body.Description,
		body.Price,
		body.Quantity,
		body.Currency,
		body.Image,
		time.Now(),
		time.Now(),
	).RETURNING(Products.AllColumns)

	if e := statement.Query(contr.DB, &dest); e != nil {
		c.JSON(http.StatusInternalServerError, response{Status: http.StatusInternalServerError, Success: false, Message: e.Error()})
		return
	}

	c.Writer.Header().Add("Location", fmt.Sprintf("/api/v1/products/%d", dest.ID))
	c.JSON(http.StatusOK, response{Status: http.StatusOK, Success:true, Data: dest})
}

func (contr *Controller) UpdateProduct(c *gin.Context){
	var productUri productUri
	var existingProduct model.Products
	var dest model.Products
	var body updateProductBody

	if e := c.ShouldBindUri(&productUri); e != nil {
		c.JSON(http.StatusBadRequest, response{Status: http.StatusBadRequest, Success: false, Message: e.Error()})	
		return
	}
	if e := c.ShouldBindJSON(&body); e != nil {
		c.JSON(http.StatusBadRequest, response{Status: http.StatusBadRequest, Success: false, Message: e.Error()})
		return	
	}
	if e := Products.SELECT(Products.AllColumns).FROM(Products).WHERE(Products.ID.EQ(Int(productUri.Id))).Query(contr.DB, &existingProduct); e != nil {
		c.JSON(http.StatusNotFound, response{Status: http.StatusNotFound, Success: false, Message: qrm.ErrNoRows.Error()})	
		return
	}
	body.UpdatedAt = time.Now()
	if e := copy(existingProduct, &body); e != nil {
		c.JSON(http.StatusInternalServerError, response{Status: http.StatusInternalServerError, Success: false, Message: e.Error()})	
		return
	}
	var statement = Products.UPDATE(Products.MutableColumns).MODEL(body).WHERE(Products.ID.EQ(Int(productUri.Id))).RETURNING(Products.AllColumns)
	if e := statement.Query(contr.DB, &dest); e != nil {
		c.JSON(http.StatusInternalServerError, response{Status: http.StatusInternalServerError, Success: false, Message: e.Error()})	
		return
	} 

	c.JSON(http.StatusOK, response{Status: http.StatusOK, Success:true, Data: dest})
}

func (contr *Controller) DeleteProduct(c *gin.Context){
	var productUri productUri
	var result sql.Result
	var e error

	if e := c.ShouldBindUri(&productUri); e != nil {
		c.JSON(http.StatusBadRequest, response{Status: http.StatusBadRequest, Success: false, Message: e.Error()})	
		return
	}

	var statement = Products.DELETE().WHERE(Products.ID.EQ(Int(productUri.Id)))
	if result, e = statement.Exec(contr.DB); e != nil {
		c.JSON(http.StatusInternalServerError, response{Status: http.StatusInternalServerError, Success: false, Message: e.Error()})
		return
	}
	if n, _ := result.RowsAffected(); n == 0 {
		c.JSON(http.StatusNotFound, response{Status: http.StatusNotFound, Success: false, Message: qrm.ErrNoRows.Error()})	
		return
	}

	c.JSON(http.StatusOK, response{Status: http.StatusOK, Success:true})
}

func (contr *Controller) GetOrders(c *gin.Context){
	var statement = Orders.SELECT(Orders.AllColumns).FROM(Orders)
	var dest []model.Orders

	if e := statement.Query(contr.DB, &dest); e != nil {
		c.JSON(http.StatusInternalServerError, response{Status: http.StatusInternalServerError,Success: false, Message: e.Error()})
		return
	}
	if len(dest) == 0 {
		c.JSON(http.StatusNotFound, response{Status: http.StatusNotFound,Success: false, Message: qrm.ErrNoRows.Error()})
		return
	}

	c.JSON(http.StatusOK, response{Status: http.StatusOK, Success: true, Data: dest})
}
func (contr *Controller) GetOrder(c *gin.Context){
	var orderUri orderUri
	var dest model.Orders
	if e := c.ShouldBindUri(&orderUri); e != nil {
		c.JSON(http.StatusBadRequest, response{Status: http.StatusBadRequest,Success: false, Message: e.Error()})
		return
	}

	var statement = Orders.SELECT(Orders.AllColumns).FROM(Orders).WHERE(Orders.ID.EQ(Int(orderUri.Id)))
	if e := statement.Query(contr.DB, &dest); e != nil {
		if e == qrm.ErrNoRows {
			c.JSON(http.StatusNotFound, response{Status: http.StatusNotFound,Success: false, Message: e.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, response{Status: http.StatusInternalServerError,Success: false, Message: e.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, response{Status: http.StatusOK, Success: true, Data: dest})
}
func (contr *Controller) CreateOrder(c *gin.Context){
	var body postOrderBody
	var dest model.Orders
	var trans *sql.Tx
	var e error
	if e := c.ShouldBindJSON(&body); e != nil {
		c.JSON(http.StatusBadRequest, response{Status: http.StatusBadRequest,Success: false, Message: e.Error()})
		return
	}	
	if trans, e = contr.DB.Begin(); e != nil {
		c.JSON(http.StatusInternalServerError, response{Status: http.StatusInternalServerError,Success: false, Message: e.Error()})
		return
	}
	// checks if the product is available
	var matchedProduct model.Products
	var matchedProductStatement = Products.SELECT(
		Products.ID, 
		Products.Quantity, 
		Products.Price,
		).FROM(Products).WHERE(Products.ID.EQ(Int(body.ProductID)).AND(Products.Quantity.GT_EQ(Int(body.Quantity))))

	if e := matchedProductStatement.Query(trans, &matchedProduct); e != nil {
		if e == qrm.ErrNoRows {
			c.JSON(http.StatusNotFound, response{Status: http.StatusNotFound,Success: false, Message: fmt.Sprintf("Not enough resources for product with id %d", body.ProductID)})
		} else {
			c.JSON(http.StatusInternalServerError, response{Status: http.StatusInternalServerError,Success: false, Message: e.Error()})
		}
		trans.Rollback()
		return	
	}
	// creates the order
	var total = matchedProduct.Price * float64(body.Quantity)
	var statement = Orders.INSERT(
		Orders.UID,
		Orders.ProductID,
		Orders.Quantity,
		Orders.Total,
		Orders.UserID,
		Orders.CreatedAt,
		Orders.UpdatedAt,
	).VALUES(
		body.UID,
		body.ProductID,
		body.Quantity,
		total,
		body.UserID,
		time.Now(),
		time.Now(),
	).RETURNING(Orders.AllColumns)
	if e := statement.Query(trans, &dest); e != nil {
		c.JSON(http.StatusInternalServerError, response{Status: http.StatusInternalServerError,Success: false, Message: e.Error()})
		trans.Rollback()
		return
	}

	// decreases the quantity from products table
	var newQuantity = int64(*matchedProduct.Quantity) - int64(body.Quantity)
	var updateProductQuantityStatement = Products.UPDATE(Products.Quantity).SET(Int(newQuantity)).WHERE(Products.ID.EQ(Int(body.ProductID)))
	if _, e := updateProductQuantityStatement.Exec(trans); e != nil { 	
		c.JSON(http.StatusInternalServerError, response{Status: http.StatusInternalServerError,Success: false, Message: e.Error()})
		trans.Rollback()
		return
	}

	trans.Commit()

	if e := contr.Publisher.Pub("products_create_order_success", dest, "someUniqueID"); e != nil {
		c.JSON(http.StatusInternalServerError, response{Status: http.StatusInternalServerError,Success: false, Message: e.Error()})
		return
	}
	c.JSON(http.StatusOK, response{Status: http.StatusOK,Success: true, Data: dest})
}

func (contr *Controller) DeleteOrder(c *gin.Context){
	var orderUri orderUri
	var result sql.Result
	var e error

	if e := c.ShouldBindUri(&orderUri); e != nil {
		c.JSON(http.StatusBadRequest, response{Status: http.StatusBadRequest,Success: false, Message: e.Error()})
		return
	}
	if result, e = Orders.DELETE().WHERE(Orders.ID.EQ(Int(orderUri.Id))).Exec(contr.DB); e != nil {
		c.JSON(http.StatusInternalServerError, response{Status: http.StatusInternalServerError,Success: false, Message: e.Error()})
		return
	}
	if n, _ := result.RowsAffected(); n == 0 {
		c.JSON(http.StatusNotFound, response{Status: http.StatusNotFound,Success: false, Message: qrm.ErrNoRows.Error()})
		return
	}

	c.JSON(http.StatusNoContent, response{Status: http.StatusNoContent, Success: true})
}
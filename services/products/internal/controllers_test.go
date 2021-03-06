package internal 

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/go-jet/jet/qrm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)
var (
	env = EnvVariables{}
)

func toResponse(body io.Reader) response{
	var bf, e = ioutil.ReadAll(body)
	var res response
	if e != nil {
		log.Fatalln(e)
	}
	if e := json.Unmarshal(bf, &res); e != nil {
		log.Fatalln(e)
	}
	return res
}

type ControllersSuite struct {
	suite.Suite
	DB *sql.DB
	DBMock sqlmock.Sqlmock
}

func (cs *ControllersSuite) SetupSuite(){
	gin.SetMode("test")

	var db, dbMock, e = sqlmock.New()
	if e != nil {
		log.Fatalln(e)
	}
	cs.DB = db
	cs.DBMock =  dbMock
}

func (cs *ControllersSuite) TestGetProducts(){
	var table = []struct{
		setExpectations func(sqlmock.Sqlmock)
		expectedRes response
	}{
		{
			setExpectations: func(dbMock sqlmock.Sqlmock){
				var e = fmt.Errorf("Internal Error")
				dbMock.ExpectQuery("SELECT (.+) FROM public.products").WillReturnError(e)
			},
			expectedRes: response{
				Status: 500,
				Success: false,
				Message:  "Internal Error",
			},
		},
		{
			setExpectations: func(dbMock sqlmock.Sqlmock){
				var columns = []string{"products.name"}
				var rows = sqlmock.NewRows(columns)
				dbMock.ExpectQuery("SELECT (.+) FROM public.products").WillReturnRows(rows)
			},
			expectedRes: response{
				Status: 404,
				Success: false,
				Message: qrm.ErrNoRows.Error(),
			},
		},
		{
			setExpectations: func(dbMock sqlmock.Sqlmock){
				var columns = []string{"products.name"}
				var rows = sqlmock.NewRows(columns).AddRow("product")
				dbMock.ExpectQuery("SELECT (.+) FROM public.products").WillReturnRows(rows)
			},
			expectedRes: response{
				Status: 200,
				Success: true,
				Data: []interface {}{map[string]interface {}{"createdAt":"0001-01-01T00:00:00Z", "currency":interface {}(nil), "description":"", "id":0.0, "image":"", "name":"product", "price":0.0, "quantity":interface {}(nil), "updatedAt":"0001-01-01T00:00:00Z"}},
			},
		},
	}

	var assert = assert.New(cs.T())
	for _, row := range table {
		var w = httptest.NewRecorder()
		var c, _ = gin.CreateTestContext(w)
		var contr = Controller{Env: env, DB: cs.DB}

		row.setExpectations(cs.DBMock)
		contr.GetProducts(c)

		var result = w.Result()
		var res = toResponse(result.Body)

		assert.Equal(res, row.expectedRes)
		cs.DBMock.ExpectationsWereMet()
	}
}

func (cs *ControllersSuite) TestGetProduct(){
	var table = []struct{
		id string
		setExpectations func(sqlmock.Sqlmock)
		expectedRes response
	}{
		{
			id: "",
			setExpectations: func(dbMock sqlmock.Sqlmock){},
			expectedRes: response{
				Status: 400,
				Success: false,
				Message:  "Key: 'productUri.Id' Error:Field validation for 'Id' failed on the 'required' tag",
			},
		},
		{
			id: "1",
			setExpectations: func(dbMock sqlmock.Sqlmock){
				dbMock.ExpectQuery("^SELECT (.+) FROM public.products WHERE products.id = \\$1;$").
				WithArgs(1).
				WillReturnError(qrm.ErrNoRows)				
			},
			expectedRes: response{
				Status: 404,
				Success: false,
				Message: qrm.ErrNoRows.Error(),
			},
		},
		{
			id: "1",
			setExpectations: func(dbMock sqlmock.Sqlmock){
				var e = fmt.Errorf("Internal Error")
				dbMock.ExpectQuery("^SELECT (.+) FROM public.products WHERE products.id = \\$1;$").
				WithArgs(1).
				WillReturnError(e)				
			},
			expectedRes: response{
				Status: 500,
				Success: false,
				Message: "Internal Error",
			},
		},
		{
			id: "1",
			setExpectations: func(dbMock sqlmock.Sqlmock){
				var rows = sqlmock.NewRows([]string{"products.name"}).AddRow("product1")
				dbMock.ExpectQuery("^SELECT (.+) FROM public.products WHERE products.id = \\$1;$").
				WithArgs(1).
				WillReturnRows(rows)	
			},
			expectedRes: response{
				Status: 200,
				Success: true,
				Data: map[string]interface {}{"createdAt":"0001-01-01T00:00:00Z", "currency":interface {}(nil), "description":"", "id":0.0, "image":"", "name":"product1", "price":0.0, "quantity":interface {}(nil), "updatedAt":"0001-01-01T00:00:00Z"},
			},
		},
	}


	var assert = assert.New(cs.T())
	for _, row := range table {
		var w = httptest.NewRecorder()
		var c, _ = gin.CreateTestContext(w)
		var contr = Controller{Env: env, DB: cs.DB}

		c.Params = append(c.Params, gin.Param{Key:"id", Value:row.id})
		row.setExpectations(cs.DBMock)
		contr.GetProduct(c)

		var result = w.Result()
		var res = toResponse(result.Body)

		assert.Equal(res, row.expectedRes)
		cs.DBMock.ExpectationsWereMet()
	}
}

func (cs *ControllersSuite) TestCreateProduct(){
	var table = []struct{
		body []byte
		setExpectations func(sqlmock.Sqlmock)
		expectedRes response
	}{
		{
			body: []byte(`{}`),
			setExpectations: func(dbMock sqlmock.Sqlmock){},
			expectedRes: response{
				Status: 400,
				Success: false,
				Message:  "Key: 'postProductBody.Name' Error:Field validation for 'Name' failed on the 'required' tag\nKey: 'postProductBody.Description' Error:Field validation for 'Description' failed on the 'required' tag\nKey: 'postProductBody.Price' Error:Field validation for 'Price' failed on the 'required' tag\nKey: 'postProductBody.Quantity' Error:Field validation for 'Quantity' failed on the 'required' tag\nKey: 'postProductBody.Currency' Error:Field validation for 'Currency' failed on the 'required' tag\nKey: 'postProductBody.Image' Error:Field validation for 'Image' failed on the 'required' tag",
			},
		},
		{
			body: []byte(`{
				"name": "product",
				"price": 10.0,
				"quantity": 1,
				"description": "description",
				"image": "image",
				"currency": "GBP"
			}`),
			setExpectations: func(dbMock sqlmock.Sqlmock){
				dbMock.ExpectQuery("^INSERT INTO public.products(.+) VALUES (.+)").
				WithArgs("product","description", 10.0,1,"GBP","image", sqlmock.AnyArg(), sqlmock.AnyArg()).
				WillReturnError(fmt.Errorf("Internal Error"))
			},
			expectedRes: response{
				Status: 500,
				Success: false,
				Message:  "Internal Error",
			},
		},
		{
			body: []byte(`{
				"name": "product",
				"price": 10.0,
				"quantity": 1,
				"description": "description",
				"image": "image",
				"currency": "GBP"
			}`),
			setExpectations: func(dbMock sqlmock.Sqlmock){
				var rows = sqlmock.NewRows([]string{"name"}).AddRow("name")
				dbMock.ExpectQuery("^INSERT INTO public.products(.+) VALUES (.+)").
				WithArgs("product","description", 10.0,1,"GBP","image", sqlmock.AnyArg(), sqlmock.AnyArg()).
				WillReturnRows(rows)
			},
			expectedRes: response{
				Status: 200,
				Success: true,
				Data: map[string]interface {}{"createdAt":"0001-01-01T00:00:00Z", "currency":interface {}(nil), "description":"", "id":0.0, "image":"", "name":"", "price":0.0, "quantity":interface {}(nil), "updatedAt":"0001-01-01T00:00:00Z"},
			},
		},
	}


	var assert = assert.New(cs.T())
	for _, row := range table {
		var w = httptest.NewRecorder()
		var c, _ = gin.CreateTestContext(w)
		var req = httptest.NewRequest("GET", "/api/v1/products", bytes.NewReader(row.body))
		c.Request = req
		var contr = Controller{Env: env, DB: cs.DB}

		row.setExpectations(cs.DBMock)
		contr.CreateProduct(c)

		var result = w.Result()
		var res = toResponse(result.Body)

		assert.Equal(res, row.expectedRes)
		cs.DBMock.ExpectationsWereMet()
	}
}


func (cs *ControllersSuite) TestUpdateProduct(){
	var table = []struct{
		id string
		body []byte
		setExpectations func(sqlmock.Sqlmock)
		expectedRes response
	}{
		{
			id: "",
			body: []byte(`{}`),
			setExpectations: func(dbMock sqlmock.Sqlmock){},
			expectedRes: response{
				Status: 400,
				Success: false,
				Message: "Key: 'productUri.Id' Error:Field validation for 'Id' failed on the 'required' tag",
			},
		},
		{
			id: "1",
			body: []byte(`{}`),
			setExpectations: func(dbMock sqlmock.Sqlmock){},
			expectedRes: response{
				Status: 400,
				Success: false,
				Message: "Key: 'updateProductBody.Name' Error:Field validation for 'Name' failed on the 'required_without_all' tag",
			},
		},
		{
			id: "1",
			body: []byte(`{
				"currency": "EURO"
			}`),
			setExpectations: func(dbMock sqlmock.Sqlmock){
				dbMock.ExpectQuery("SELECT (.+) FROM public.products WHERE products.id = \\$1;$").
				WithArgs(1).
				WillReturnError(qrm.ErrNoRows)
			},
			expectedRes: response{
				Status: 404,
				Success: false,
				Message: "qrm: no rows in result set",
			},
		},
		{
			id: "1",
			body: []byte(`{
				"currency": "EURO"
			}`),
			setExpectations: func(dbMock sqlmock.Sqlmock){
				dbMock.ExpectQuery("SELECT (.+) FROM public.products WHERE products.id = \\$1;$").
				WithArgs(1).
				WillReturnRows(sqlmock.NewRows([]string{"products.currency"}).AddRow("USD"))

				dbMock.ExpectQuery("UPDATE public.products SET (.+) WHERE products.id = \\$9").
				WillReturnError(fmt.Errorf("Internal Error"))
			},
			expectedRes: response{
				Status: 500,
				Success: false,
				Message: "Internal Error",
			},
		},
		{
			id: "1",
			body: []byte(`{
				"currency": "EURO"
			}`),
			setExpectations: func(dbMock sqlmock.Sqlmock){
				dbMock.ExpectQuery("SELECT (.+) FROM public.products WHERE products.id = \\$1;$").
				WithArgs(1).
				WillReturnRows(sqlmock.NewRows([]string{"products.currency"}).AddRow("USD"))

				dbMock.ExpectQuery("UPDATE public.products SET (.+) WHERE products.id = \\$9").
				WillReturnRows(sqlmock.NewRows([]string{"products.name"}).AddRow("product"))
			},
			expectedRes: response{
				Status: 200,
				Success: true,
				Data: map[string]interface {}{"createdAt":"0001-01-01T00:00:00Z", "currency":interface {}(nil), "description":"", "id":0.0, "image":"", "name":"product", "price":0.0, "quantity":interface {}(nil), "updatedAt":"0001-01-01T00:00:00Z"},
			},
		},
	}


	var assert = assert.New(cs.T())
	for _, row := range table {
		var w = httptest.NewRecorder()
		var c, _ = gin.CreateTestContext(w)
		var req = httptest.NewRequest("PUT", "/api/v1/products/:id", bytes.NewReader(row.body))
		c.Request = req
		var contr = Controller{Env: env, DB: cs.DB}

		c.Params = append(c.Params, gin.Param{Key:"id", Value: row.id})
		row.setExpectations(cs.DBMock)
		contr.UpdateProduct(c)

		var result = w.Result()
		var res = toResponse(result.Body)

		assert.Equal(res, row.expectedRes)
		cs.DBMock.ExpectationsWereMet()
	}
}



func (cs *ControllersSuite) TestDeleteProduct(){
	var table = []struct{
		id string
		code int
		setExpectations func(sqlmock.Sqlmock)
		expectedRes response
	}{
		{
			id: "",
			code: 400,
			setExpectations: func(dbMock sqlmock.Sqlmock){},
			expectedRes: response{
				Status: 400,
				Success: false,
				Message:  "Key: 'productUri.Id' Error:Field validation for 'Id' failed on the 'required' tag",
			},
		},
		{
			id: "1",
			code: 500,
			setExpectations: func(dbMock sqlmock.Sqlmock){
				var e = fmt.Errorf("Internal Error")
				dbMock.ExpectExec("^DELETE FROM public.products WHERE products.id = \\$1;$").
				WithArgs(1).
				WillReturnError(e)				
			},
			expectedRes: response{
				Status: 500,
				Success: false,
				Message: "Internal Error",
			},
		},
		{
			id: "1",
			code: 404,
			setExpectations: func(dbMock sqlmock.Sqlmock){
				var result = sqlmock.NewResult(0,0)
				dbMock.ExpectExec("^DELETE FROM public.products WHERE products.id = \\$1;$").
				WithArgs(1).
				WillReturnResult(result)				
			},
			expectedRes: response{
				Status: 404,
				Success: false,
				Message: qrm.ErrNoRows.Error(),
			},
		},
		{
			id: "1",
			code: 200,
			setExpectations: func(dbMock sqlmock.Sqlmock){
				var result = sqlmock.NewResult(1,1)
				dbMock.ExpectExec("^DELETE FROM public.products WHERE products.id = \\$1;$").
				WithArgs(1).
				WillReturnResult(result)				
			},
			expectedRes: response{
				Status: 200,
				Success: true,
			},
		},
	}


	var assert = assert.New(cs.T())
	for _, row := range table {
		var w = httptest.NewRecorder()
		var c, _ = gin.CreateTestContext(w)
		var contr = Controller{Env: env, DB: cs.DB}

		c.Params = append(c.Params, gin.Param{Key:"id", Value:row.id})
		row.setExpectations(cs.DBMock)
		contr.DeleteProduct(c)

		var result = w.Result()
		assert.Equal(result.StatusCode,row.code)

		if result.StatusCode != http.StatusNoContent {
			var res = toResponse(result.Body)
			assert.Equal(res, row.expectedRes)
		}			
		cs.DBMock.ExpectationsWereMet()
	}
}


func (cs *ControllersSuite) TestGetOrders(){
	var table = []struct{
		setExpectations func(sqlmock.Sqlmock)
		expectedRes response
	}{
		{
			setExpectations: func(dbMock sqlmock.Sqlmock){
				var e = fmt.Errorf("Internal Error")
				dbMock.ExpectQuery("SELECT (.+) FROM public.orders").WillReturnError(e)
			},
			expectedRes: response{
				Status: 500,
				Success: false,
				Message:  "Internal Error",
			},
		},
		{
			setExpectations: func(dbMock sqlmock.Sqlmock){
				var rows = sqlmock.NewRows([]string{"orders.uid"})
				dbMock.ExpectQuery("SELECT (.+) FROM public.orders").WillReturnRows(rows)
			},
			expectedRes: response{
				Status: 404,
				Success: false,
				Message:  qrm.ErrNoRows.Error(),
			},
		},
		{
			setExpectations: func(dbMock sqlmock.Sqlmock){
				var rows = sqlmock.NewRows([]string{"orders.uid"}).AddRow("uid")
				dbMock.ExpectQuery("SELECT (.+) FROM public.orders").WillReturnRows(rows)
			},
			expectedRes: response{
				Status: 200,
				Success: true,
				Data: []interface {}{map[string]interface {}{"createdAt":"0001-01-01T00:00:00Z", "id":0.0, "productId":0.0, "quantity":0.0, "status":interface {}(nil), "total":0.0, "uid":"uid", "updatedAt":"0001-01-01T00:00:00Z", "userId":0.0}},
			},
		},
	}

	var assert = assert.New(cs.T())
	for _, row := range table {
		var w = httptest.NewRecorder()
		var c, _ = gin.CreateTestContext(w)
		var contr = Controller{Env: env, DB: cs.DB}

		row.setExpectations(cs.DBMock)
		contr.GetOrders(c)

		var result = w.Result()
		var res = toResponse(result.Body)

		assert.Equal(res, row.expectedRes)
		cs.DBMock.ExpectationsWereMet()
	}
}


func (cs *ControllersSuite) TestGetOrder(){
	var table = []struct{
		id string
		setExpectations func(sqlmock.Sqlmock)
		expectedRes response
	}{
		{
			id: "",
			setExpectations: func(dbMock sqlmock.Sqlmock){},
			expectedRes: response{
				Status: 400,
				Success: false,
				Message:  "Key: 'orderUri.Id' Error:Field validation for 'Id' failed on the 'required' tag",
			},
		},
		{
			id: "1",
			setExpectations: func(dbMock sqlmock.Sqlmock){
				dbMock.ExpectQuery("SELECT (.+) FROM public.orders WHERE orders.id = \\$1").
				WithArgs(1).
				WillReturnError(qrm.ErrNoRows)
			},
			expectedRes: response{
				Status: 404,
				Success: false,
				Message:  qrm.ErrNoRows.Error(),
			},
		},
		{
			id: "1",
			setExpectations: func(dbMock sqlmock.Sqlmock){
				dbMock.ExpectQuery("SELECT (.+) FROM public.orders WHERE orders.id = \\$1").
				WithArgs(1).
				WillReturnError(fmt.Errorf("Internal Error"))
			},
			expectedRes: response{
				Status: 500,
				Success: false,
				Message:  "Internal Error",
			},
		},
		{
			id: "1",
			setExpectations: func(dbMock sqlmock.Sqlmock){
				var rows = sqlmock.NewRows([]string{"orders.uid"}).AddRow("uid")
				dbMock.ExpectQuery("SELECT (.+) FROM public.orders WHERE orders.id = \\$1").
				WithArgs(1).
				WillReturnRows(rows)
			},
			expectedRes: response{
				Status: 200,
				Success: true,
				Data: map[string]interface {}{"createdAt":"0001-01-01T00:00:00Z", "id":0.0, "productId":0.0, "quantity":0.0, "status":interface {}(nil), "total":0.0, "uid":"uid", "updatedAt":"0001-01-01T00:00:00Z", "userId":0.0},
			},
		},
	}


	var assert = assert.New(cs.T())
	for _, row := range table {
		var w = httptest.NewRecorder()
		var c, _ = gin.CreateTestContext(w)
		var contr = Controller{Env: env, DB: cs.DB}

		c.Params = append(c.Params, gin.Param{Key:"id", Value:row.id})
		row.setExpectations(cs.DBMock)
		contr.GetOrder(c)

		var result = w.Result()
		var res = toResponse(result.Body)

		assert.Equal(res, row.expectedRes)
		cs.DBMock.ExpectationsWereMet()
	}
}


func (cs *ControllersSuite) TestDeleteOrder(){
	var table = []struct{
		id string
		code int
		setExpectations func(sqlmock.Sqlmock)
		expectedRes response
	}{
		{
			id: "",
			code: 400,
			setExpectations: func(dbMock sqlmock.Sqlmock){},
			expectedRes: response{
				Status: 400,
				Success: false,
				Message:  "Key: 'orderUri.Id' Error:Field validation for 'Id' failed on the 'required' tag",
			},
		},
		{
			id: "1",
			code: 500,
			setExpectations: func(dbMock sqlmock.Sqlmock){
				var e = fmt.Errorf("Internal Error")
				dbMock.ExpectExec("^DELETE FROM public.orders WHERE orders.id = \\$1;$").
				WithArgs(1).
				WillReturnError(e)				
			},
			expectedRes: response{
				Status: 500,
				Success: false,
				Message: "Internal Error",
			},
		},
		{
			id: "1",
			code: 404,
			setExpectations: func(dbMock sqlmock.Sqlmock){
				var result = sqlmock.NewResult(0,0)
				dbMock.ExpectExec("^DELETE FROM public.orders WHERE orders.id = \\$1;$").
				WithArgs(1).
				WillReturnResult(result)				
			},
			expectedRes: response{
				Status: 404,
				Success: false,
				Message: qrm.ErrNoRows.Error(),
			},
		},
		{
			id: "1",
			code: 204,
			setExpectations: func(dbMock sqlmock.Sqlmock){
				var result = sqlmock.NewResult(1,1)
				dbMock.ExpectExec("^DELETE FROM public.orders WHERE orders.id = \\$1;$").
				WithArgs(1).
				WillReturnResult(result)				
			},
			expectedRes: response{

			},
		},
	}


	var assert = assert.New(cs.T())
	for _, row := range table {
		var w = httptest.NewRecorder()
		var c, _ = gin.CreateTestContext(w)
		var contr = Controller{Env: env, DB: cs.DB}

		c.Params = append(c.Params, gin.Param{Key:"id", Value:row.id})
		row.setExpectations(cs.DBMock)
		contr.DeleteOrder(c)

		var result = w.Result()
		assert.Equal(result.StatusCode,row.code)

		if result.StatusCode != http.StatusNoContent {
			var res = toResponse(result.Body)
			assert.Equal(res, row.expectedRes)
		}			
		cs.DBMock.ExpectationsWereMet()
	}
}

func (cs *ControllersSuite) TestCreateOrder(){
	var table = []struct{
		body []byte
		setExpectations func(sqlmock.Sqlmock)
		expectedRes response
	}{
		{
			body: []byte(`{}`),
			setExpectations: func(dbMock sqlmock.Sqlmock){},
			expectedRes: response{
				Status: 400,
				Success: false,
				Message:  "Key: 'postOrderBody.UID' Error:Field validation for 'UID' failed on the 'required' tag\nKey: 'postOrderBody.Quantity' Error:Field validation for 'Quantity' failed on the 'required' tag\nKey: 'postOrderBody.UserID' Error:Field validation for 'UserID' failed on the 'required' tag",
			},
		},
		{
			body: []byte(`{
				"uid": "uid",
				"productId": 1,
				"quantity": 2,
				"userId": 1
			}`),
			setExpectations: func(dbMock sqlmock.Sqlmock){
				var e = fmt.Errorf("No transaction")
				dbMock.ExpectBegin().WillReturnError(e)
			},
			expectedRes: response{
				Status: 500,
				Success: false,
				Message:  "No transaction",
			},
		},
		{
			body: []byte(`{
				"uid": "uid",
				"productId": 1,
				"quantity": 2,
				"userId": 1
			}`),
			setExpectations: func(dbMock sqlmock.Sqlmock){
				dbMock.ExpectBegin()
				dbMock.ExpectQuery("SELECT (.+) FROM public.products WHERE (.+)").WithArgs(1,2).WillReturnError(qrm.ErrNoRows)
				dbMock.ExpectRollback()
			},
			expectedRes: response{
				Status: 404,
				Success: false,
				Message:  "Not enough resources for product with id 1",
			},
		},
		{
			body: []byte(`{
				"uid": "uid",
				"productId": 1,
				"quantity": 2,
				"userId": 1
			}`),
			setExpectations: func(dbMock sqlmock.Sqlmock){
				dbMock.ExpectBegin()
				dbMock.ExpectQuery("SELECT (.+) FROM public.products WHERE (.+)").WithArgs(1,2).WillReturnError(fmt.Errorf("Internal Error"))
				dbMock.ExpectRollback()
			},
			expectedRes: response{
				Status: 500,
				Success: false,
				Message:  "Internal Error",
			},
		},
		{
			body: []byte(`{
				"uid": "uid",
				"productId": 1,
				"quantity": 2,
				"userId": 1
			}`),
			setExpectations: func(dbMock sqlmock.Sqlmock){
				var rows = sqlmock.NewRows([]string{}).AddRow()
				dbMock.ExpectBegin()
				dbMock.ExpectQuery("SELECT (.+) FROM public.products WHERE (.+)").WithArgs(1,2).WillReturnRows(rows)
				dbMock.ExpectRollback()
			},
			expectedRes: response{
				Status: 404,
				Success: false,
				Message:  "Not enough resources for product with id 1",
			},
		},
	}


	var assert = assert.New(cs.T())
	for _, row := range table {
		var w = httptest.NewRecorder()
		var c, _ = gin.CreateTestContext(w)
		var req = httptest.NewRequest("POST", "/api/v1/orders", bytes.NewReader(row.body))
		c.Request = req
		var contr = Controller{Env: env, DB: cs.DB}

		row.setExpectations(cs.DBMock)
		contr.CreateOrder(c)

		var result = w.Result()
		var res = toResponse(result.Body)

		assert.Equal(res, row.expectedRes)
		cs.DBMock.ExpectationsWereMet()
	}
}

func TestControllerSuite(t *testing.T){
	suite.Run(t, new(ControllersSuite))
}
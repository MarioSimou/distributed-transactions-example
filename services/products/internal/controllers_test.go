package internal 

import (
	"bytes"
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
	DB qrm.DB
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
				Data: []interface {}{map[string]interface {}{"CreatedAt":"0001-01-01T00:00:00Z", "Currency":interface {}(nil), "ID":float64(0), "Name":"product", "Price": float64(0), "Quantity":interface {}(nil), "UpdatedAt":"0001-01-01T00:00:00Z"}},
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
				Data: map[string]interface {}{"CreatedAt":"0001-01-01T00:00:00Z", "Currency":interface {}(nil), "ID":float64(0), "Name":"product1", "Price":float64(0), "Quantity":interface {}(nil), "UpdatedAt":"0001-01-01T00:00:00Z"},
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
				Message:  "Key: 'postProductBody.Name' Error:Field validation for 'Name' failed on the 'required' tag\nKey: 'postProductBody.Price' Error:Field validation for 'Price' failed on the 'required' tag\nKey: 'postProductBody.Quantity' Error:Field validation for 'Quantity' failed on the 'required' tag\nKey: 'postProductBody.Currency' Error:Field validation for 'Currency' failed on the 'required' tag",
			},
		},
		{
			body: []byte(`{
				"name": "product",
				"price": 10.0,
				"quantity": 1,
				"currency": "GBP"
			}`),
			setExpectations: func(dbMock sqlmock.Sqlmock){
				dbMock.ExpectQuery("^INSERT INTO public.products(.+) VALUES (.+)").
				WithArgs("product", 10.0,1,"GBP", sqlmock.AnyArg(), sqlmock.AnyArg()).
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
				"currency": "GBP"
			}`),
			setExpectations: func(dbMock sqlmock.Sqlmock){
				var rows = sqlmock.NewRows([]string{"name"}).AddRow("name")
				dbMock.ExpectQuery("^INSERT INTO public.products(.+) VALUES (.+)").
				WithArgs("product", 10.0,1,"GBP", sqlmock.AnyArg(), sqlmock.AnyArg()).
				WillReturnRows(rows)
			},
			expectedRes: response{
				Status: 200,
				Success: true,
				Data: map[string]interface {}{"CreatedAt":"0001-01-01T00:00:00Z", "Currency":interface {}(nil), "ID":float64(0), "Name":"", "Price":float64(0), "Quantity":interface {}(nil), "UpdatedAt":"0001-01-01T00:00:00Z"},
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

				dbMock.ExpectQuery("UPDATE public.products SET (.+) WHERE products.id = \\$7").
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

				dbMock.ExpectQuery("UPDATE public.products SET (.+) WHERE products.id = \\$7").
				WillReturnRows(sqlmock.NewRows([]string{"products.name"}).AddRow("product"))
			},
			expectedRes: response{
				Status: 200,
				Success: true,
				Data: map[string]interface {}{"CreatedAt":"0001-01-01T00:00:00Z", "Currency":interface {}(nil), "ID":float64(0), "Name":"product", "Price":float64(0), "Quantity":interface {}(nil), "UpdatedAt":"0001-01-01T00:00:00Z"},
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
			code: 204,
			setExpectations: func(dbMock sqlmock.Sqlmock){
				var result = sqlmock.NewResult(1,1)
				dbMock.ExpectExec("^DELETE FROM public.products WHERE products.id = \\$1;$").
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
				Data: []interface {}{map[string]interface {}{"CreatedAt":"0001-01-01T00:00:00Z", "ID":float64(0), "ProductID":float64(0), "Quantity":float64(0), "Status":interface {}(nil), "Total":float64(0), "UID":"uid", "UpdatedAt":"0001-01-01T00:00:00Z", "UserID":float64(0)}},
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
				Data: map[string]interface {}{"CreatedAt":"0001-01-01T00:00:00Z", "ID":0.0, "ProductID":0.0, "Quantity": 0.0, "Status":interface {}(nil), "Total":0.0, "UID":"uid", "UpdatedAt":"0001-01-01T00:00:00Z", "UserID":0.0},
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

func TestControllerSuite(t *testing.T){
	suite.Run(t, new(ControllersSuite))
}
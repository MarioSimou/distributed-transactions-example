package internal

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http/httptest"
	"testing"
	"net/http"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/go-jet/jet/qrm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var (
	env = EnvironmentVariables{}
)

type ControllerSuite struct {
	suite.Suite
	DB *sql.DB
	DBMock sqlmock.Sqlmock
}

func (cs *ControllerSuite) SetupTest() {
	gin.SetMode("test")
	var db, dbMock, e = sqlmock.New()
	if e != nil {
		log.Fatalln(e)
	}

	cs.DB = db
	cs.DBMock = dbMock
}


func toResponse(body io.ReadCloser) (Response, error){
	var bf, e = ioutil.ReadAll(body.(io.Reader))
	if e != nil {
		return Response{}, e
	}

	var res Response
	json.Unmarshal(bf, &res)
	return res, nil
}

func (cs *ControllerSuite) TestControllerCreateUser(){
	var table = []struct{
		setAssertions func(db sqlmock.Sqlmock)
		body []byte
		expectedRes Response
	}{
		{
			setAssertions: func(db sqlmock.Sqlmock){},
			body: []byte("{}"),
			expectedRes: Response{
				Status: 400,
				Success: false,
				Message: "Key: 'postUserBody.Username' Error:Field validation for 'Username' failed on the 'required' tag\nKey: 'postUserBody.Email' Error:Field validation for 'Email' failed on the 'required' tag\nKey: 'postUserBody.Password' Error:Field validation for 'Password' failed on the 'required' tag",
			},
		},
		{
			setAssertions: func(dbMock sqlmock.Sqlmock){
				var e = fmt.Errorf("DB error")
				var query = "\nINSERT INTO public.users (.+)"
				
				dbMock.ExpectQuery(query).WillReturnError(e)
			},
			body: []byte(`{
				"username": "test",
				"email": "test@gmail.com",
				"password": "123456678",
				"balance": 10.0
			}`),
			expectedRes: Response{
				Status: 500,
				Success: false,
				Message: "DB error",
			},
		},
		{
			setAssertions: func(dbMock sqlmock.Sqlmock){
				var query = "\nINSERT INTO public.users (.+)"
				var columns = []string{"users.username"}

				var rows = sqlmock.NewRows(columns).AddRow("test")
				dbMock.ExpectQuery(query).
					WithArgs("test","test@gmail.com","123456678", 10.0, sqlmock.AnyArg(),sqlmock.AnyArg()).
					WillReturnRows(rows)
			},
			body: []byte(`{
				"username": "test",
				"email": "test@gmail.com",
				"password": "123456678",
				"balance": 10.0
			}`),
			expectedRes: Response{
				Status: 200,
				Success: true,
				Data: map[string]interface {}{"balance":interface {}(nil), "createdAt":interface {}(nil), "email":"", "id":0.0, "password":"", "updatedAt":interface {}(nil), "username":"test"},
			},
		},
	}

	var assert = assert.New(cs.T())
	for _, row := range table {
		var w = httptest.NewRecorder()
		var req = httptest.NewRequest("GET","/api/v1/users", bytes.NewReader(row.body))
		var c, _ = gin.CreateTestContext(w)
		var contr = Controller{EnvVariables: env, DB: cs.DB}

		c.Request = req
		row.setAssertions(cs.DBMock)

		// performs the call
		contr.CreateUser(c)

		var result = w.Result()
		var res, _ = toResponse(result.Body)

		assert.EqualValues(res,row.expectedRes)
		if e := cs.DBMock.ExpectationsWereMet(); e != nil {
			log.Fatalln(e)
		}
	}

}

func (cs *ControllerSuite) TestGetUser(){
	var table = []struct {
		id string
		expectedRes Response
		setAssertions func(sqlmock.Sqlmock)
	}{
		{
			id: "",
			expectedRes: Response{
				Status: 400,
				Success: false,
				Message: "Key: 'userId.Id' Error:Field validation for 'Id' failed on the 'required' tag",
			},
			setAssertions: func(dbMock sqlmock.Sqlmock){},
		},
		{
			id: "1",
			expectedRes: Response{
				Status: 404,
				Success: false,
				Message: "qrm: no rows in result set",
			},
			setAssertions: func(dbMock sqlmock.Sqlmock){
				dbMock.ExpectQuery("SELECT (.+)").WithArgs(1).WillReturnError(qrm.ErrNoRows)
			},
		},
		{
			id: "1",
			expectedRes: Response{
				Status: 500,
				Success: false,
				Message: "Internal Error",
			},
			setAssertions: func(dbMock sqlmock.Sqlmock){
				var e = fmt.Errorf("Internal Error")
				dbMock.ExpectQuery("SELECT (.+)").WithArgs(1).WillReturnError(e)
			},
		},
		{
			id: "1",
			expectedRes: Response{
				Status: 200,
				Success: true,
				Data: map[string]interface {}{"balance":interface {}(nil), "createdAt":interface {}(nil), "email":"", "id":0.0, "password":"", "updatedAt":interface {}(nil), "username":"test"},
			},
			setAssertions: func(dbMock sqlmock.Sqlmock){
				var columns = []string{"users.username"}
				var rows = sqlmock.NewRows(columns).AddRow("test")
				dbMock.ExpectQuery("SELECT (.+)").WithArgs(1).WillReturnRows(rows)
			},
		},
	}

	var t = cs.T()
	var assert = assert.New(t)

	for _, row := range table {
		var w = httptest.NewRecorder()
		var req = httptest.NewRequest("GET", fmt.Sprintf("/api/v1/users/:id"), nil)
		var c, _ = gin.CreateTestContext(w)
		var contr = Controller{EnvVariables: env, DB: cs.DB}
		

		row.setAssertions(cs.DBMock)
		c.Request = req
		c.Params = gin.Params{gin.Param{Key: "id", Value: row.id}}

		contr.GetUser(c)		

		var result = w.Result()
		var res, _ = toResponse(result.Body)
		
		assert.Equal(res, row.expectedRes)
		if e := cs.DBMock.ExpectationsWereMet(); e != nil {
			log.Fatalln(e)
		}
	}
}


func (cs *ControllerSuite) TestGetUsers(){
	var table = []struct {
		expectedRes Response
		setAssertions func(dbMock sqlmock.Sqlmock)
	}{
		{
			expectedRes: Response{
				Status: 500,
				Success: false,
				Message: "Internal Error",
			},
			setAssertions: func(dbMock sqlmock.Sqlmock){
				var e = fmt.Errorf("Internal Error")
				dbMock.ExpectQuery("SELECT (.+) FROM public.users").WillReturnError(e)
			},
		},
		{
			expectedRes: Response{
				Status: 404,
				Success: false,
				Message: "qrm: no rows in result set",
			},
			setAssertions: func(dbMock sqlmock.Sqlmock){
				var columns = []string{"users.username"}
				var rows = sqlmock.NewRows(columns)
				dbMock.ExpectQuery("SELECT (.+) FROM public.users").WillReturnRows(rows)
			},
		},
		{
			expectedRes: Response{
				Status: 200,
				Success: true,
				Data: []interface {}{map[string]interface {}{"balance":interface {}(nil), "createdAt":interface {}(nil), "email":"", "id":0.0, "password":"", "updatedAt":interface {}(nil), "username":"test1"}, map[string]interface {}{"balance":interface {}(nil), "createdAt":interface {}(nil), "email":"", "id":0.0, "password":"", "updatedAt":interface {}(nil), "username":"test2"}},
			},
			setAssertions: func(dbMock sqlmock.Sqlmock){
				var columns = []string{"users.username"}
				var rows = sqlmock.NewRows(columns).AddRow("test1").AddRow("test2")
				dbMock.ExpectQuery("SELECT (.+) FROM public.users").WillReturnRows(rows)
			},
		},
	}

	var t = cs.T()
	var assert = assert.New(t)

	for _, row := range table {
		var w = httptest.NewRecorder()
		var req = httptest.NewRequest("GET", fmt.Sprintf("/api/v1/users"), nil)
		var c, _ = gin.CreateTestContext(w)
		var contr = Controller{EnvVariables: env, DB: cs.DB}
		

		row.setAssertions(cs.DBMock)
		c.Request = req
		contr.GetUsers(c)		

		var result = w.Result()
		var res, _ = toResponse(result.Body)
		
		assert.Equal(res, row.expectedRes)
		if e := cs.DBMock.ExpectationsWereMet(); e != nil {
			log.Fatalln(e)
		}
	}
}

func (cs *ControllerSuite) TestDeleteUser(){
	var table = []struct {
		id string
		code int
		expectedRes Response
		setAssertions func(sqlmock.Sqlmock)
	}{
		{
			id: "",
			code: 400,
			expectedRes: Response{
				Status: 400,
				Success: false,
				Message: "Key: 'userId.Id' Error:Field validation for 'Id' failed on the 'required' tag",
			},
			setAssertions: func(dbMock sqlmock.Sqlmock){},
		},
		{
			id: "1",
			code: 500,
			expectedRes: Response{
				Status: 500,
				Success: false,
				Message: "Internal Error",
			},
			setAssertions: func(dbMock sqlmock.Sqlmock){
				var e = fmt.Errorf("Internal Error")
				dbMock.ExpectExec("DELETE FROM public.users WHERE users.id = \\$1").WithArgs(1).WillReturnError(e)
			},
		},
		{
			id: "1",
			code: 404,
			expectedRes: Response{
				Status: 404,
				Success: false,
				Message: qrm.ErrNoRows.Error(),
			},
			setAssertions: func(dbMock sqlmock.Sqlmock){
				var result = sqlmock.NewResult(0,0)
				dbMock.ExpectExec("DELETE FROM public.users WHERE users.id = \\$1").WithArgs(1).WillReturnResult(result)
			},
		},
		{
			id: "1",
			code: 204,
			expectedRes: Response{},
			setAssertions: func(dbMock sqlmock.Sqlmock){
				var result = sqlmock.NewResult(1,1)
				dbMock.ExpectExec("DELETE FROM public.users WHERE users.id = \\$1").WithArgs(1).WillReturnResult(result)
			},
		},
	}

	var t = cs.T()
	var assert = assert.New(t)

	for _, row := range table {
		var w = httptest.NewRecorder()
		var req = httptest.NewRequest("GET", fmt.Sprintf("/api/v1/users/:id"), nil)
		var c, _ = gin.CreateTestContext(w)
		var contr = Controller{EnvVariables: env, DB:cs.DB}
		

		row.setAssertions(cs.DBMock)
		c.Request = req
		c.Params = gin.Params{gin.Param{Key: "id", Value: row.id}}

		contr.DeleteUser(c)		

		var result = w.Result()
		
		var res, _ = toResponse(result.Body)
		
		assert.Equal(result.StatusCode, row.code)
		assert.Equal(res, row.expectedRes)
		if e := cs.DBMock.ExpectationsWereMet(); e != nil {
			log.Fatalln(e)
		}
	}
}


func (cs *ControllerSuite) TestUpdateUser(){
	var table = []struct {
		id string
		expectedRes Response
		setAssertions func(sqlmock.Sqlmock)
		body []byte
	}{
		{
			id: "",
			expectedRes: Response{
				Status: 400,
				Success: false,
				Message: "Key: 'userId.Id' Error:Field validation for 'Id' failed on the 'required' tag",
			},
			body: []byte(``),
			setAssertions: func(db sqlmock.Sqlmock){},
		},
		{
			id: "1",
			expectedRes: Response{
				Status: 400,
				Success: false,
				Message: "Key: 'updateUserBody.Username' Error:Field validation for 'Username' failed on the 'required_without_all' tag",
			},
			body: []byte(`{}`),
			setAssertions: func(db sqlmock.Sqlmock){},
		},
		{
			id: "1",
			expectedRes: Response{
				Status: 404,
				Success: false,
				Message: qrm.ErrNoRows.Error(),
			},
			body: []byte(`{
				"username": "test",
				"email": "test@gmail.com",
				"password": "12345678",
				"balance": 10.0
			}`),
			setAssertions: func(dbMock sqlmock.Sqlmock){
				dbMock.ExpectQuery("SELECT (.+) FROM public.users WHERE users.id = \\$1").WithArgs(1).WillReturnError(qrm.ErrNoRows)
			},
		},
		{
			id: "1",
			expectedRes: Response{
				Status: 500,
				Success: false,
				Message: "Internal Error",
			},
			body: []byte(`{
				"username": "test",
				"email": "test@gmail.com",
				"password": "12345678",
				"balance": 10.0
			}`),
			setAssertions: func(dbMock sqlmock.Sqlmock){
				var columns = []string{"users.username"}
				dbMock.
					ExpectQuery("SELECT (.+) FROM public.users WHERE users.id = \\$1").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows(columns).AddRow("test"))

				
				var e = fmt.Errorf("Internal Error")
				dbMock.
					ExpectQuery("UPDATE (.+)").WillReturnError(e)
			},
		},
		{
			id: "1",
			expectedRes: Response{
				Status: 200,
				Success: true,
				Data: map[string]interface {}{"balance":interface {}(nil), "createdAt":interface {}(nil), "email":"", "id":0.0, "password":"", "updatedAt":interface {}(nil), "username":"test"},
			},
			body: []byte(`{
				"username": "test",
				"email": "test@gmail.com",
				"password": "12345678",
				"balance": 10.0
			}`),
			setAssertions: func(dbMock sqlmock.Sqlmock){
				var columns = []string{"users.username"}
				dbMock.
					ExpectQuery("SELECT (.+) FROM public.users WHERE users.id = \\$1").
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows(columns).AddRow("test"))

				dbMock.
					ExpectQuery("UPDATE public.users SET (.+) WHERE users.id = \\$7 RETURNING (.+)").
					WithArgs("test","test@gmail.com", "12345678", 10.0, sqlmock.AnyArg(), sqlmock.AnyArg(), 1).
					WillReturnRows(sqlmock.NewRows(columns).AddRow("test"))
			},
		},
	}

	var t = cs.T()
	var assert = assert.New(t)

	for _, row := range table {
		var w = httptest.NewRecorder()
		var req = httptest.NewRequest("PUT", fmt.Sprintf("/api/v1/users/:id"), bytes.NewReader(row.body))
		var c, _ = gin.CreateTestContext(w)
		var contr = Controller{EnvVariables: env, DB: cs.DB}
		

		row.setAssertions(cs.DBMock)
		c.Request = req
		c.Params = gin.Params{gin.Param{Key: "id", Value: row.id}}

		contr.UpdateUser(c)		

		var result = w.Result()
		var res, _ = toResponse(result.Body)
		
		assert.Equal(res, row.expectedRes)
		if e := cs.DBMock.ExpectationsWereMet(); e != nil {
			log.Fatalln(e)
		}
	}
}

func (cs *ControllerSuite ) TestSignInUser(){
	var table = []struct{
		setAssertions func(sqlmock.Sqlmock)
		body []byte
		expectedRes Response
		cookie bool
	}{
		{
			setAssertions: func(sqlmock.Sqlmock){},
			body: []byte(`{}`),
			expectedRes: Response{
				Status: 400,
				Success: false,
				Message: "Key: 'signInBody.Email' Error:Field validation for 'Email' failed on the 'required' tag\nKey: 'signInBody.Password' Error:Field validation for 'Password' failed on the 'required' tag",
			},
			cookie: false,
		},
		{
			setAssertions: func(dbMock sqlmock.Sqlmock){
				var e = fmt.Errorf("Internal Error")
				dbMock.ExpectQuery("SELECT (.+) FROM public.users WHERE \\(users.email = \\$1\\) AND \\(users.password = \\$2\\)").
				WithArgs("test@gmail.com", "12345678").
				WillReturnError(e)
			},
			body: []byte(`{
				"email": "test@gmail.com",
				"password": "12345678"
			}`),
			expectedRes: Response{
				Status: 500,
				Success: false,
				Message: "Internal Error",
			},
			cookie: false,
		},
		{
			setAssertions: func(dbMock sqlmock.Sqlmock){
				var rows = sqlmock.NewRows([]string{"email"}).AddRow("test@gmail.com")
				dbMock.ExpectQuery("SELECT (.+) FROM public.users WHERE \\(users.email = \\$1\\) AND \\(users.password = \\$2\\)").
				WithArgs("test@gmail.com", "12345678").
				WillReturnRows(rows)

				var e = fmt.Errorf("Internal Error")
				dbMock.ExpectExec("INSERT INTO (.+)").WillReturnError(e)
			},
			body: []byte(`{
				"email": "test@gmail.com",
				"password": "12345678"
			}`),
			expectedRes: Response{
				Status: 500,
				Success: false,	
				Message: "Internal Error",
			},
			cookie: false,
		},
		{
			setAssertions: func(dbMock sqlmock.Sqlmock){
				var rows = sqlmock.NewRows([]string{"email"}).AddRow("test@gmail.com")
				dbMock.ExpectQuery("SELECT (.+) FROM public.users WHERE \\(users.email = \\$1\\) AND \\(users.password = \\$2\\)").
				WithArgs("test@gmail.com", "12345678").
				WillReturnRows(rows)

				var result = sqlmock.NewResult(0,0)
				dbMock.ExpectExec("INSERT INTO (.+)").	WillReturnResult(result)
			},
			body: []byte(`{
				"email": "test@gmail.com",
				"password": "12345678"
			}`),
			expectedRes: Response{
				Status: 500,
				Success: false,	
				Message: "Session was not created",
			},
			cookie: false,
		},
		{
			setAssertions: func(dbMock sqlmock.Sqlmock){
				var rows = sqlmock.NewRows([]string{"email"}).AddRow("test@gmail.com")
				dbMock.ExpectQuery("SELECT (.+) FROM public.users WHERE \\(users.email = \\$1\\) AND \\(users.password = \\$2\\)").
				WithArgs("test@gmail.com", "12345678").
				WillReturnRows(rows)

				var result = sqlmock.NewResult(1,1)
				dbMock.ExpectExec("INSERT INTO (.+)").WillReturnResult(result)
			},
			body: []byte(`{
				"email": "test@gmail.com",
				"password": "12345678"
			}`),
			expectedRes: Response{
				Status: 200,
				Success: true,
				Data:map[string]interface {}{"balance":interface {}(nil), "createdAt":interface {}(nil), "email":"", "id":0.0, "password":"", "updatedAt":interface {}(nil), "username":""},
			},
			cookie: true,
		},
	}

	var t = cs.T()
	var assert = assert.New(t)
	for _, row := range table {
		var w = httptest.NewRecorder()
		var req = httptest.NewRequest("POST", "/signin", bytes.NewReader(row.body))
		var c, _ = gin.CreateTestContext(w)
		c.Request = req
		var contr = Controller{EnvVariables: env, DB: cs.DB}

		row.setAssertions(cs.DBMock)

		contr.SignInUser(c)

		if row.cookie {
			var cookieHeader = w.HeaderMap["Set-Cookie"]
			assert.True(len(cookieHeader) > 0, row.cookie)
		}

		var result = w.Result()
		var res, _ = toResponse(result.Body)
		assert.EqualValues(res, row.expectedRes)

		cs.DBMock.ExpectationsWereMet()
	}
}

func (cs *ControllerSuite ) TestSignInUserWithGUID(){
	var table = []struct{
		cookieValue string
		setAssertions func(sqlmock.Sqlmock)
		expectedRes Response
		returnsCookie bool
	}{
		{
			cookieValue: "",
			setAssertions: func(dbMock sqlmock.Sqlmock){},
			expectedRes: Response{
				Status: 404,
				Success: false,
				Message: "http: named cookie not present",
			},
			returnsCookie: false,
		},
		{
			cookieValue: "guid",
			setAssertions: func(dbMock sqlmock.Sqlmock){
				dbMock.ExpectQuery("SELECT (.+) FROM public.sessions WHERE sessions.guid = \\$1").WithArgs("guid").WillReturnError(qrm.ErrNoRows)

			},
			expectedRes: Response{
				Status: 404,
				Success: false,
				Message: qrm.ErrNoRows.Error(),
			},
			returnsCookie: false,
		},
		{
			cookieValue: "guid",
			setAssertions: func(dbMock sqlmock.Sqlmock){
				var rows = sqlmock.NewRows([]string{"sessions.guid"}).AddRow("guid")
				dbMock.ExpectQuery("SELECT (.+) FROM public.sessions WHERE sessions.guid = \\$1").WithArgs("guid").WillReturnRows(rows)

				dbMock.ExpectQuery("SELECT (.+) FROM public.users WHERE (.+)").WithArgs(0).WillReturnError(qrm.ErrNoRows)				
			},
			expectedRes: Response{
				Status: 404,
				Success: false,
				Message: qrm.ErrNoRows.Error(),
			},
			returnsCookie: false,
		},
		{
			cookieValue: "guid",
			setAssertions: func(dbMock sqlmock.Sqlmock){
				var rows1 = sqlmock.NewRows([]string{"sessions.guid"}).AddRow("guid")
				dbMock.ExpectQuery("SELECT (.+) FROM public.sessions WHERE sessions.guid = \\$1").WithArgs("guid").WillReturnRows(rows1)

				var rows2 = sqlmock.NewRows([]string{"users.email"}).AddRow("test@gmail.com")
				dbMock.ExpectQuery("SELECT (.+) FROM public.users WHERE (.+)").WithArgs(0).WillReturnRows(rows2)	

				var e = fmt.Errorf("Internal Error")
				dbMock.ExpectExec("UPDATE public.sessions (.+) WHERE sessions.guid = \\$2").WithArgs(sqlmock.AnyArg(), "guid").WillReturnError(e)
			},
			expectedRes: Response{
				Status: 500,
				Success: false,
				Message: "Internal Error",
			},
			returnsCookie: false,
		},
		{
			cookieValue: "guid",
			setAssertions: func(dbMock sqlmock.Sqlmock){
				var rows1 = sqlmock.NewRows([]string{"sessions.guid"}).AddRow("guid")
				dbMock.ExpectQuery("SELECT (.+) FROM public.sessions WHERE sessions.guid = \\$1").WithArgs("guid").WillReturnRows(rows1)

				var rows2 = sqlmock.NewRows([]string{"users.email"}).AddRow("test@gmail.com")
				dbMock.ExpectQuery("SELECT (.+) FROM public.users WHERE (.+)").WithArgs(0).WillReturnRows(rows2)	

				var result = sqlmock.NewResult(1,1)
				dbMock.ExpectExec("UPDATE public.sessions (.+) WHERE sessions.guid = \\$2").WithArgs(sqlmock.AnyArg(), "guid").WillReturnResult(result)
			},
			expectedRes: Response{
				Status: 200,
				Success: true,
				Data: map[string]interface {}{"balance":interface {}(nil), "createdAt":interface {}(nil), "email":"test@gmail.com", "id":0.0, "password":"", "updatedAt":interface {}(nil), "username":""},
			},
			returnsCookie: true,
		},
	}

	var t = cs.T()
	var assert = assert.New(t)
	for _, row := range table {
		var w = httptest.NewRecorder()
		var c, _ = gin.CreateTestContext(w)
		var req, _ = http.NewRequest("GET", "/signin", nil)

		if row.cookieValue != ""{
			var cookie = http.Cookie{
				Name: "sid",
				Value: row.cookieValue,
				Domain: "ecommerce.com",
				Path: "/",
				MaxAge: 3600,
			}
			req.AddCookie(&cookie)	
		}
		var contr = Controller{EnvVariables: env, DB: cs.DB}

		c.Request = req
		row.setAssertions(cs.DBMock)
		contr.SignInUserWithCookie(c)

		if row.returnsCookie {
			var cookieHeader = w.HeaderMap["Set-Cookie"]
			assert.True(len(cookieHeader) > 0, row.returnsCookie)
		}

		var result = w.Result()
		var res, _ = toResponse(result.Body)
		assert.EqualValues(res, row.expectedRes)

		cs.DBMock.ExpectationsWereMet()
	}
}

func (cs *ControllerSuite ) TestLogOut(){
	var table = []struct{
		cookieValue string
		expectedRes Response
		returnsCookie bool
	}{
		{
			cookieValue: "",
			returnsCookie: false,
			expectedRes: Response{
				Status: 404,
				Success:false,
				Message: "http: named cookie not present",
			},
		},
		{
			cookieValue: "value",
			returnsCookie: true,
			expectedRes: Response{},
		},
	}

	var t = cs.T()
	var assert = assert.New(t)
	for _, row := range table {
		var w = httptest.NewRecorder()
		var c, _ = gin.CreateTestContext(w)
		var req, _ = http.NewRequest("GET", "/signin", nil)
		var contr = Controller{EnvVariables: env, DB: cs.DB}

		if row.cookieValue != ""{
			var cookie = http.Cookie{
				Name: "sid",
				Value: row.cookieValue,
				Domain: "ecommerce.com",
				Path: "/",
				MaxAge: 3600,
			}
			req.AddCookie(&cookie)	
		}


		c.Request = req
		contr.LogOut(c)

		if row.returnsCookie {
			var cookieHeader = w.HeaderMap["Set-Cookie"]
			assert.True(len(cookieHeader) > 0, row.returnsCookie)
		}

		var result = w.Result()
		var res, _ = toResponse(result.Body)
		assert.EqualValues(res, row.expectedRes)
	}
}


func TestControllerSuite(t *testing.T){
	suite.Run(t, new(ControllerSuite))
}


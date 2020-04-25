package internal

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"customers/internal/mocks"

	"github.com/gin-gonic/gin"
	"github.com/go-jet/jet/qrm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

var (
	env = EnvironmentVariables{}
)

type ControllerSuite struct {
	suite.Suite
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
		setAssertions func(db *mocks.DB)
		body []byte
		expectedRes Response
	}{
		{
			setAssertions: func(db *mocks.DB){},
			body: []byte("{}"),
			expectedRes: Response{
				Status: 400,
				Success: false,
				Message: "Key: 'User.Username' Error:Field validation for 'Username' failed on the 'required' tag\nKey: 'User.Email' Error:Field validation for 'Email' failed on the 'required' tag\nKey: 'User.Password' Error:Field validation for 'Password' failed on the 'required' tag",
			},
		},
		{
			setAssertions: func(db *mocks.DB){
				var e = errors.New("DB error")
				var query = "\nINSERT INTO public.users (username, email, password, balance, created_at, updated_at) VALUES\n     ($1, $2, $3, $4, $5, $6)\nRETURNING users.id AS \"users.id\",\n          users.username AS \"users.username\",\n          users.email AS \"users.email\",\n          users.password AS \"users.password\",\n          users.balance AS \"users.balance\",\n          users.created_at AS \"users.created_at\",\n          users.updated_at AS \"users.updated_at\";\n"

				db.On("QueryContext", 
							 mock.Anything, 
							 query, 
							 "test", 
							 "test@gmail.com",
							 "123456678",
							 10.0,
							 mock.AnythingOfType("time.Time"),
							 mock.AnythingOfType("time.Time"),
							 ).Return(nil, e)
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
	}

	var assert = assert.New(cs.T())
	for _, row := range table {
		var w = httptest.NewRecorder()
		var req = httptest.NewRequest("GET","/api/v1/users", bytes.NewReader(row.body))
		var c, _ = gin.CreateTestContext(w)
		var db = &mocks.DB{}
		var contr = Controller{EnvVariables: env, DB: db}

		c.Request = req
		row.setAssertions(db)

		// performs the call
		contr.CreateUser(c)

		var result = w.Result()
		var res, _ = toResponse(result.Body)

		assert.EqualValues(res,row.expectedRes)
		db.AssertExpectations(cs.T())
	}

}

func (cs *ControllerSuite) TestGetUser(){
	var table = []struct {
		id string
		expectedRes Response
		setAssertions func(db *mocks.DB)
	}{
		{
			id: "",
			expectedRes: Response{
				Status: 400,
				Success: false,
				Message: "Key: 'userId.Id' Error:Field validation for 'Id' failed on the 'required' tag",
			},
			setAssertions: func(db *mocks.DB){},
		},
		{
			id: "1",
			expectedRes: Response{
				Status: 500,
				Success: false,
				Message: "Internal Error",
			},
			setAssertions: func(db *mocks.DB){
				var e = errors.New("Internal Error")

				db.On("QueryContext", 
			  				mock.Anything,
							 "\nSELECT users.id AS \"users.id\",\n     users.username AS \"users.username\",\n     users.email AS \"users.email\",\n     users.password AS \"users.password\",\n     users.balance AS \"users.balance\",\n     users.created_at AS \"users.created_at\",\n     users.updated_at AS \"users.updated_at\"\nFROM public.users\nWHERE users.id = $1;\n", 
							 int64(1)).Return(nil, e)
			},
		},
		{
			id: "1",
			expectedRes: Response{
				Status: 404,
				Success: false,
				Message: "qrm: no rows in result set",
			},
			setAssertions: func(db *mocks.DB){
				var e = qrm.ErrNoRows

				db.On("QueryContext", 
			  				mock.Anything,
							 "\nSELECT users.id AS \"users.id\",\n     users.username AS \"users.username\",\n     users.email AS \"users.email\",\n     users.password AS \"users.password\",\n     users.balance AS \"users.balance\",\n     users.created_at AS \"users.created_at\",\n     users.updated_at AS \"users.updated_at\"\nFROM public.users\nWHERE users.id = $1;\n", 
							 int64(1)).Return(nil, e)
			},
		},
	}

	var t = cs.T()
	var assert = assert.New(t)

	for _, row := range table {
		var w = httptest.NewRecorder()
		var req = httptest.NewRequest("GET", fmt.Sprintf("/api/v1/users/:id"), nil)
		var c, _ = gin.CreateTestContext(w)
		var db = &mocks.DB{}
		var contr = Controller{EnvVariables: env, DB: db}
		

		row.setAssertions(db)
		c.Request = req
		c.Params = gin.Params{gin.Param{Key: "id", Value: row.id}}

		contr.GetUser(c)		

		var result = w.Result()
		var res, _ = toResponse(result.Body)
		
		assert.Equal(res, row.expectedRes)
		db.AssertExpectations(t)
	}
}


func (cs *ControllerSuite) TestGetUsers(){
	var table = []struct {
		expectedRes Response
		setAssertions func(db *mocks.DB)
	}{
		{
			expectedRes: Response{
				Status: 500,
				Success: false,
				Message: "Internal error",
			},
			setAssertions: func(db *mocks.DB){
				var e = errors.New("Internal error")
				db.On("QueryContext", mock.Anything, mock.Anything).Return(nil, e)
			},
		},
		{
			expectedRes: Response{
				Status: 404,
				Success: false,
				Message: "qrm: no rows in result set",
			},
			setAssertions: func(db *mocks.DB){
				var e = qrm.ErrNoRows
				db.On("QueryContext", mock.Anything, mock.Anything).Return(nil, e)
			},
		},
	}

	var t = cs.T()
	var assert = assert.New(t)

	for _, row := range table {
		var w = httptest.NewRecorder()
		var req = httptest.NewRequest("GET", fmt.Sprintf("/api/v1/users"), nil)
		var c, _ = gin.CreateTestContext(w)
		var db = &mocks.DB{}
		var contr = Controller{EnvVariables: env, DB: db}
		

		row.setAssertions(db)
		c.Request = req
		contr.GetUsers(c)		

		var result = w.Result()
		var res, _ = toResponse(result.Body)
		
		assert.Equal(res, row.expectedRes)
		db.AssertExpectations(t)
	}
}

func (cs *ControllerSuite) TestDeleteUser(){
	var table = []struct {
		id string
		expectedRes Response
		setAssertions func(db *mocks.DB)
	}{
		{
			id: "",
			expectedRes: Response{
				Status: 400,
				Success: false,
				Message: "Key: 'userId.Id' Error:Field validation for 'Id' failed on the 'required' tag",
			},
			setAssertions: func(db *mocks.DB){},
		},
		{
			id: "1",
			expectedRes: Response{
				Status: 500,
				Success: false,
				Message: "Internal Error",
			},
			setAssertions: func(db *mocks.DB){
				var e = errors.New("Internal Error")

				db.On("QueryContext", 
			  				mock.Anything,
								"\nDELETE FROM public.users\nWHERE users.id = $1\nRETURNING users.id AS \"users.id\";\n", 
							 int64(1)).Return(nil, e)
			},
		},
		{
			id: "1",
			expectedRes: Response{
				Status: 404,
				Success: false,
				Message: "qrm: no rows in result set",
			},
			setAssertions: func(db *mocks.DB){
				var e = qrm.ErrNoRows

				db.On("QueryContext", 
			  				mock.Anything,
							 "\nDELETE FROM public.users\nWHERE users.id = $1\nRETURNING users.id AS \"users.id\";\n", 
							 int64(1)).Return(nil, e)
			},
		},
	}

	var t = cs.T()
	var assert = assert.New(t)

	for _, row := range table {
		var w = httptest.NewRecorder()
		var req = httptest.NewRequest("GET", fmt.Sprintf("/api/v1/users/:id"), nil)
		var c, _ = gin.CreateTestContext(w)
		var db = &mocks.DB{}
		var contr = Controller{EnvVariables: env, DB: db}
		

		row.setAssertions(db)
		c.Request = req
		c.Params = gin.Params{gin.Param{Key: "id", Value: row.id}}

		contr.DeleteUser(c)		

		var result = w.Result()
		var res, _ = toResponse(result.Body)
		
		assert.Equal(res, row.expectedRes)
		db.AssertExpectations(t)
	}
}


func (cs *ControllerSuite) TestUpdateUser(){
	var table = []struct {
		id string
		expectedRes Response
		setAssertions func(db *mocks.DB)
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
			setAssertions: func(db *mocks.DB){},
		},
		{
			id: "1",
			expectedRes: Response{
				Status: 400,
				Success: false,
				Message: "Key: 'UpdateUserBody.Username' Error:Field validation for 'Username' failed on the 'required_without_all' tag",
			},
			body: []byte(`{}`),
			setAssertions: func(db *mocks.DB){},
		},
		{
			id: "1",
			expectedRes: Response{
				Status: 404,
				Success: false,
				Message: "qrm: no rows in result set",
			},
			body: []byte(`{
				"username": "test"
			}`),
			setAssertions: func(db *mocks.DB){
				var e = qrm.ErrNoRows
				db.On("QueryContext", mock.Anything, "\nSELECT users.id AS \"users.id\",\n     users.username AS \"users.username\",\n     users.email AS \"users.email\",\n     users.password AS \"users.password\",\n     users.balance AS \"users.balance\",\n     users.created_at AS \"users.created_at\",\n     users.updated_at AS \"users.updated_at\"\nFROM public.users\nWHERE users.id = $1;\n", int64(1)).Return(nil, e)
			},
		},
	}

	var t = cs.T()
	var assert = assert.New(t)

	for _, row := range table {
		var w = httptest.NewRecorder()
		var req = httptest.NewRequest("PUT", fmt.Sprintf("/api/v1/users/:id"), bytes.NewReader(row.body))
		var c, _ = gin.CreateTestContext(w)
		var db = &mocks.DB{}
		var contr = Controller{EnvVariables: env, DB: db}
		

		row.setAssertions(db)
		c.Request = req
		c.Params = gin.Params{gin.Param{Key: "id", Value: row.id}}

		contr.UpdateUser(c)		

		var result = w.Result()
		var res, _ = toResponse(result.Body)
		
		assert.Equal(res, row.expectedRes)
		db.AssertExpectations(t)
	}
}


func TestControllerSuite(t *testing.T){
	suite.Run(t, new(ControllerSuite))
}


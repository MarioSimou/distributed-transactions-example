package internal

import (
	"fmt"
	"net/http"
	"time"

	"customers/internal/models/customers/public/model"
	. "customers/internal/models/customers/public/table"

	. "github.com/go-jet/jet/postgres"

	"github.com/gin-gonic/gin"
	"github.com/go-jet/jet/qrm"
)

type Controller struct {
	EnvVariables EnvironmentVariables
	DB qrm.DB
}

func (contr *Controller) Ping(c *gin.Context){
	c.JSON(200, gin.H{"message": "pong"})
}

func (contr *Controller) GetUser(c *gin.Context){
	var userId userId
	var dest model.Users

	if e := c.ShouldBindUri(&userId); e != nil {
		c.JSON(http.StatusBadRequest, Response{Status: http.StatusBadRequest, Success:false, Message: e.Error()})
		return
	}

	var statement = Users.SELECT(Users.AllColumns).FROM(Users).WHERE(Users.ID.EQ(Int(userId.Id)))

	if e := statement.Query(contr.DB, &dest); e != nil {
		if e == qrm.ErrNoRows {
			c.JSON(http.StatusNotFound, Response{Status: http.StatusNotFound, Success:false, Message: e.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, Response{Status: http.StatusInternalServerError, Success:false, Message: e.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, Response{Status: http.StatusOK, Success: true, Data: dest})
}

func (contr *Controller) GetUsers(c *gin.Context){
	var dest []model.Users
	var statement = Users.SELECT(Users.AllColumns).FROM(Users)

	if e := statement.Query(contr.DB, &dest); e != nil {
		if e == qrm.ErrNoRows {
			c.JSON(http.StatusNotFound, Response{Status: http.StatusNotFound, Success:false, Message: e.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, Response{Status: http.StatusInternalServerError, Success:false, Message: e.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, Response{Status: http.StatusOK, Success: true, Data: dest})
}

func (contr *Controller) CreateUser(c *gin.Context){
	var user User
	var dest model.Users
	if e := c.ShouldBindJSON(&user); e != nil {
		c.JSON(http.StatusBadRequest, Response{Status: http.StatusBadRequest, Success:false, Message: e.Error()})
		return
	}
	
	var statement = Users.INSERT(
		Users.Username,
		Users.Email, 
		Users.Password, 
		Users.Balance, 
		Users.CreatedAt, 
		Users.UpdatedAt,
		).VALUES(
			user.Username, 
			user.Email, 
			user.Password, 
			user.Balance,
			time.Now(), 
			time.Now(),
		).RETURNING(Users.AllColumns)
		
	if e := statement.Query(contr.DB, &dest); e != nil {
		fmt.Println("Error: ", e)
		c.JSON(http.StatusInternalServerError, Response{Status: http.StatusInternalServerError, Success:false, Message: e.Error()})
		return
	}

	c.Writer.Header().Set("Location", fmt.Sprintf("/api/v1/users/%d", dest.ID))
	c.JSON(http.StatusOK, Response{Status: http.StatusOK, Success: true, Data: dest})
}

func (contr *Controller) DeleteUser(c *gin.Context){
	var userId userId
	var user model.Users

	if e := c.ShouldBindUri(&userId); e != nil {
		c.JSON(http.StatusBadRequest, Response{Status: http.StatusBadRequest, Success:false, Message: e.Error()})
		return
	}
	var statement = Users.DELETE().WHERE(Users.ID.EQ(Int(userId.Id))).RETURNING(Users.ID)
	if e := statement.Query(contr.DB, &user) ; e != nil {
		if e == qrm.ErrNoRows {
			c.JSON(http.StatusNotFound, Response{Status: http.StatusNotFound, Success:false, Message: e.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, Response{Status: http.StatusInternalServerError, Success:false, Message: e.Error()})
		}
		return
	}

	c.JSON(http.StatusNoContent, Response{Status: http.StatusNoContent, Success: true})
}

func (contr *Controller) UpdateUser(c *gin.Context){
	var body UpdateUserBody
	var dest model.Users
	var userId userId
	var existingUser model.Users


	if e := c.ShouldBindUri(&userId); e != nil {
		c.JSON(http.StatusBadRequest, Response{Status: http.StatusBadRequest, Success:false, Message: e.Error()})
		return
	}
	if e := c.ShouldBindJSON(&body); e != nil {
		c.JSON(http.StatusBadRequest, Response{Status: http.StatusBadRequest, Success:false, Message: e.Error()})
		return
	}
	if e := SELECT(Users.AllColumns).FROM(Users).WHERE(Users.ID.EQ(Int(userId.Id))).Query(contr.DB, &existingUser); e != nil {
		c.JSON(http.StatusNotFound, Response{Status: http.StatusNotFound, Success:false, Message: e.Error()})
		return
	}
	if e := copy(existingUser, &body); e != nil {
		c.JSON(http.StatusInternalServerError, Response{Status: http.StatusInternalServerError, Success:false, Message: e.Error()})
		return
	}

	var statement = Users.
									UPDATE(Users.MutableColumns).
									MODEL(body).
									WHERE(Users.ID.EQ(Int(userId.Id))).
									RETURNING(Users.AllColumns)

	if e := statement.Query(contr.DB,&dest) ; e != nil {
		c.JSON(http.StatusInternalServerError, Response{Status: http.StatusInternalServerError, Success:false, Message: e.Error()})
		return
	}

	c.JSON(http.StatusOK, Response{Status: http.StatusOK, Success: true, Data: dest})
}

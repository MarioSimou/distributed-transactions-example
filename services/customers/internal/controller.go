package internal

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"customers/internal/models/customers/public/model"
	. "customers/internal/models/customers/public/table"

	"github.com/google/uuid"

	. "github.com/go-jet/jet/postgres"

	"github.com/gin-gonic/gin"
	"github.com/go-jet/jet/qrm"
)

type Controller struct {
	EnvVariables EnvironmentVariables
	DB *sql.DB
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
		c.JSON(http.StatusInternalServerError, Response{Status: http.StatusInternalServerError, Success:false, Message: e.Error()})
		return
	}
	if len(dest) == 0 {
		c.JSON(http.StatusNotFound, Response{Status: http.StatusNotFound, Success:false, Message: "User not found"})
		return
	}

	c.JSON(http.StatusOK, Response{Status: http.StatusOK, Success: true, Data: dest})
}

func (contr *Controller) SignInUser(c *gin.Context){
	var body signInBody
	var dest model.Users

	if e := c.ShouldBindJSON(&body); e != nil {
		c.JSON(http.StatusBadRequest, Response{Status: http.StatusBadRequest, Success: false, Message: e.Error()})
		return
	}

	var statement = Users.SELECT(Users.AllColumns).FROM(Users).WHERE(Users.Email.EQ(String(body.Email)).AND(Users.Password.EQ(String(body.Password))))
	if e := statement.Query(contr.DB, &dest); e != nil {
		c.JSON(http.StatusInternalServerError, Response{Status: http.StatusInternalServerError, Success: false, Message: e.Error()})
		return
	}
	var uid = uuid.New()
	c.SetCookie("sid", uid.String(), int(3600), "/", "ecommerce.com", false, false)
	c.JSON(http.StatusOK, Response{Status: http.StatusOK, Success: true, Data: dest})
}

func (contr *Controller) CreateUser(c *gin.Context){
	var user postUserBody
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
		c.JSON(http.StatusInternalServerError, Response{Status: http.StatusInternalServerError, Success:false, Message: e.Error()})
		return
	}

	c.Writer.Header().Set("Location", fmt.Sprintf("/api/v1/users/%d", dest.ID))
	c.JSON(http.StatusOK, Response{Status: http.StatusOK, Success: true, Data: dest})
}

func (contr *Controller) DeleteUser(c *gin.Context){
	var userId userId
	var e error
	var result sql.Result

	if e := c.ShouldBindUri(&userId); e != nil {
		c.JSON(http.StatusBadRequest, Response{Status: http.StatusBadRequest, Success:false, Message: e.Error()})
		return
	}
	var statement = Users.
										DELETE().
										WHERE(Users.ID.EQ(Int(userId.Id)))

	if result, e = statement.Exec(contr.DB); e != nil {
		c.JSON(http.StatusInternalServerError, Response{Status: http.StatusInternalServerError, Success:false, Message: e.Error()})
		return
	}			
	if n, _ := result.RowsAffected(); n == 0 {
		c.JSON(http.StatusNotFound, Response{Status: http.StatusNotFound, Success:false, Message: qrm.ErrNoRows.Error()})
	}

	c.JSON(http.StatusNoContent, Response{Status: http.StatusNoContent, Success: true})
}

func (contr *Controller) UpdateUser(c *gin.Context){
	var body updateUserBody
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

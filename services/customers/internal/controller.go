package internal

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"customers/internal/models/customers/public/model"
	. "customers/internal/models/customers/public/table"

	. "github.com/go-jet/jet/postgres"

	"github.com/gin-gonic/gin"
	"github.com/go-jet/jet/qrm"
	"github.com/google/uuid"
)

const (
	SESSION_NOT_CREATE = "Session was not created"
	MAX_AGE = 3600
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
		c.JSON(http.StatusNotFound, Response{Status: http.StatusNotFound, Success:false, Message: qrm.ErrNoRows.Error()})
		return
	}

	c.JSON(http.StatusOK, Response{Status: http.StatusOK, Success: true, Data: dest})
}

func (contr *Controller) LogOut(c *gin.Context){
	var sessionGUID string
	var e error

	if sessionGUID, e = c.Cookie("sid"); e != nil {
		c.JSON(http.StatusNotFound, Response{Status: http.StatusNotFound, Success: false, Message: e.Error()})
		return
	}

	c.SetCookie("sid", sessionGUID, -1, "/", "ecommerce.com", false, true)
	c.JSON(http.StatusNoContent, Response{Status: http.StatusNoContent, Success: true})
}

func (contr *Controller) SignInUserWithCookie(c *gin.Context){
	var session model.Sessions
var user model.Users
	var sessionGUID string
	var e error

	if sessionGUID, e = c.Cookie("sid"); e != nil {
		c.JSON(http.StatusNotFound, Response{Status: http.StatusNotFound, Success: false, Message: e.Error()})
		return
	}
	var getSessionStmt = Sessions.SELECT(Sessions.AllColumns).FROM(Sessions).WHERE(Sessions.GUID.EQ(String(sessionGUID)))
	if e := getSessionStmt.Query(contr.DB, &session); e != nil {
		c.JSON(http.StatusNotFound, Response{Status: http.StatusNotFound, Success: false, Message: e.Error()})
		return
	}
	var getUserStmt = Users.SELECT(Users.AllColumns).FROM(Users).WHERE(Users.ID.EQ(Int(int64(session.UserID))));
	if e:= getUserStmt.Query(contr.DB, &user); e != nil {
		c.JSON(http.StatusNotFound, Response{Status: http.StatusNotFound, Success: false, Message: e.Error()})
		return
	}
	var updateSessionStmt = Sessions.UPDATE(Sessions.ExpiresAt).SET(time.Now().Add(time.Second * time.Duration(MAX_AGE))).WHERE(Sessions.GUID.EQ(String(session.GUID)))
	if _, e := updateSessionStmt.Exec(contr.DB); e != nil {
		c.JSON(http.StatusInternalServerError, Response{Status: http.StatusInternalServerError, Success: false, Message: e.Error()})
		return
	}

	c.SetCookie("sid", session.GUID, MAX_AGE, "/", "ecommerce.com", false, true)
	c.JSON(http.StatusOK, Response{Status: http.StatusOK, Success: true, Data: user})
}

func (contr *Controller) SignInUser(c *gin.Context){
	var body signInBody
	var dest model.Users
	var result sql.Result
	var e error

	if e := c.ShouldBindJSON(&body); e != nil {
		c.JSON(http.StatusBadRequest, Response{Status: http.StatusBadRequest, Success: false, Message: e.Error()})
		return
	}
	var statement = Users.SELECT(Users.AllColumns).FROM(Users).WHERE(Users.Email.EQ(String(body.Email)).AND(Users.Password.EQ(String(body.Password))))
	if e := statement.Query(contr.DB, &dest); e != nil {
		c.JSON(http.StatusInternalServerError, Response{Status: http.StatusInternalServerError, Success: false, Message: e.Error()})
		return
	}

	var sid = uuid.New()
	var createdAt = time.Now()
	var expiresAt = createdAt.Add(time.Second * time.Duration(MAX_AGE))
	var session = model.Sessions{UserID: dest.ID,GUID: sid.String(),CreatedAt: &createdAt, ExpiresAt: expiresAt}
	var sessionStatement = Sessions.INSERT(Sessions.MutableColumns).MODEL(session)

	if result, e = sessionStatement.Exec(contr.DB); e != nil {
		c.JSON(http.StatusInternalServerError, Response{Status: http.StatusInternalServerError, Success: false, Message: e.Error()})
		return
	}
	if n, _ := result.RowsAffected(); n == 0 {
		c.JSON(http.StatusInternalServerError, Response{Status: http.StatusInternalServerError, Success: false, Message: SESSION_NOT_CREATE})
		return
	}

	c.SetCookie("sid", sid.String(), MAX_AGE, "/", "ecommerce.com", false, true)
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

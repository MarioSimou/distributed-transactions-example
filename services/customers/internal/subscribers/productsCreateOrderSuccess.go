package subscribers

import (
	"context"
	r "customers/internal/rabbitmq"
	"database/sql"
	"encoding/json"
	"time"

	"customers/internal/models/customers/public/model"
	. "customers/internal/models/customers/public/table"

	. "github.com/go-jet/jet/postgres"
)

type order struct {
	ID        int32 `json:"id" sql:"primary_key"`
	UID       string `json:"uid"`
	ProductID int32 `json:"productId"`
	Quantity  int32 `json:"quantity"`
	Total     float64 `json:"total"`
	UserID    int32 `json:"userId"`
	Status    string `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type orderUser struct {
	Order order
	User model.Users
	Err string
}

func publishSuccess(publisher r.PublisherInterface, body interface{}, requestID string){
	publisher.Pub("customers_charge_customer_success", body, requestID)
}
func publishFailure(publisher r.PublisherInterface, body interface{}, requestID string){
	publisher.Pub("customers_charge_customer_failure", body, requestID)
}

func productsCreateOrderSuccess (msg r.Message, ctx context.Context) error {
	var db = ctx.Value("DB").(*sql.DB)
	var publisher = ctx.Value("Publisher").(r.PublisherInterface)
	var o order
	var user model.Users
	var requestID = msg.CorrelationId

	if e := json.Unmarshal(msg.Body, &o); e != nil {
		publishFailure(publisher, orderUser{Err: e.Error()}, requestID)
		return e
	}
	
	var filter = Users.ID.EQ(Int(int64(o.UserID)))
	var getUser = Users.SELECT(Users.AllColumns).FROM(Users).WHERE(filter)
	if e := getUser.Query(db, &user); e != nil {
		publishFailure(publisher, orderUser{Order: o, Err: e.Error()}, requestID)
		return e
	}
	var newBalance = *user.Balance - o.Total 
	var newUser model.Users
	var updateUser = Users.UPDATE(Users.Balance).SET(newBalance).WHERE(filter).RETURNING(Users.AllColumns)
	if e := updateUser.Query(db,&newUser); e != nil {
		publishFailure(publisher, orderUser{Order:o, Err: e.Error()}, requestID)
		return e
	}

	publishSuccess(publisher, orderUser{Order: o, User: newUser}, requestID)
	return nil
}


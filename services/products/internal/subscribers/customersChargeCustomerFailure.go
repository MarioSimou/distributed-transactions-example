package subscribers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"products/internal/models/products/public/model"
	. "products/internal/models/products/public/table"
	r "products/internal/rabbitmq"

	. "github.com/go-jet/jet/postgres"
)

type user struct {
	ID        int32 `json:"id" sql:"primary_key"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Balance   *float64 `json:"balance"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`
}

type orderUser struct {
	User user `json:"user"`
	Order model.Orders `json:"order"`
	Err string `json:"error"`
}


func customersChargeCustomerFailure (msg r.Message, ctx context.Context) error {
	var db = ctx.Value("DB").(*sql.DB)
	var ou orderUser

	if e := json.Unmarshal(msg.Body, &ou); e != nil {
		return e
	}

	fmt.Println(ou)
	var filter = Orders.ID.EQ(Int(int64(ou.Order.ID)))
	var updateOrder = Orders.UPDATE(Orders.Status).SET(model.OrderStatus_Declined).WHERE(filter)
	if _, e := updateOrder.Exec(db); e != nil {
		return e
	}
	
	if ou.Err != "" {
		return fmt.Errorf(ou.Err)
	}

	return nil
}

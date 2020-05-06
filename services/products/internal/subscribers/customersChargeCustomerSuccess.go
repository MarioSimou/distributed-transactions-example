package subscribers

import (
	"context"
	"database/sql"
	"encoding/json"

	"products/internal/models/products/public/model"
	. "products/internal/models/products/public/table"
	r "products/internal/rabbitmq"

	. "github.com/go-jet/jet/postgres"
)

func customersChargeCustomerSuccess (msg r.Message, ctx context.Context) error {
	var db = ctx.Value("DB").(*sql.DB)
	var ou orderUser

	if e := json.Unmarshal(msg.Body, &ou); e != nil {
		return e
	}

	var filter = Orders.ID.EQ(Int(int64(ou.Order.ID)))
	var updateOrder = Orders.UPDATE(Orders.Status).SET(model.OrderStatus_Approved).WHERE(filter)
	if _, e := updateOrder.Exec(db); e != nil {
		return e
	}

	return nil
}

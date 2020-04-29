//
// Code generated by go-jet DO NOT EDIT.
// Generated at Wednesday, 29-Apr-20 09:00:05 UTC
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package table

import (
	"github.com/go-jet/jet/postgres"
)

var Orders = newOrdersTable()

type OrdersTable struct {
	postgres.Table

	//Columns
	ID        postgres.ColumnInteger
	UID       postgres.ColumnString
	ProductID postgres.ColumnInteger
	Quantity  postgres.ColumnInteger
	Total     postgres.ColumnFloat
	UserID    postgres.ColumnInteger
	CreatedAt postgres.ColumnTimestampz
	UpdatedAt postgres.ColumnTimestampz
	Status    postgres.ColumnString

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

// creates new OrdersTable with assigned alias
func (a *OrdersTable) AS(alias string) *OrdersTable {
	aliasTable := newOrdersTable()

	aliasTable.Table.AS(alias)

	return aliasTable
}

func newOrdersTable() *OrdersTable {
	var (
		IDColumn        = postgres.IntegerColumn("id")
		UIDColumn       = postgres.StringColumn("uid")
		ProductIDColumn = postgres.IntegerColumn("product_id")
		QuantityColumn  = postgres.IntegerColumn("quantity")
		TotalColumn     = postgres.FloatColumn("total")
		UserIDColumn    = postgres.IntegerColumn("user_id")
		CreatedAtColumn = postgres.TimestampzColumn("created_at")
		UpdatedAtColumn = postgres.TimestampzColumn("updated_at")
		StatusColumn    = postgres.StringColumn("status")
	)

	return &OrdersTable{
		Table: postgres.NewTable("public", "orders", IDColumn, UIDColumn, ProductIDColumn, QuantityColumn, TotalColumn, UserIDColumn, CreatedAtColumn, UpdatedAtColumn, StatusColumn),

		//Columns
		ID:        IDColumn,
		UID:       UIDColumn,
		ProductID: ProductIDColumn,
		Quantity:  QuantityColumn,
		Total:     TotalColumn,
		UserID:    UserIDColumn,
		CreatedAt: CreatedAtColumn,
		UpdatedAt: UpdatedAtColumn,
		Status:    StatusColumn,

		AllColumns:     postgres.ColumnList{IDColumn, UIDColumn, ProductIDColumn, QuantityColumn, TotalColumn, UserIDColumn, CreatedAtColumn, UpdatedAtColumn, StatusColumn},
		MutableColumns: postgres.ColumnList{UIDColumn, ProductIDColumn, QuantityColumn, TotalColumn, UserIDColumn, CreatedAtColumn, UpdatedAtColumn, StatusColumn},
	}
}

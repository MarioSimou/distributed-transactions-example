//
// Code generated by go-jet DO NOT EDIT.
// Generated at Tuesday, 05-May-20 21:51:13 UTC
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package table

import (
	"github.com/go-jet/jet/postgres"
)

var Products = newProductsTable()

type ProductsTable struct {
	postgres.Table

	//Columns
	ID          postgres.ColumnInteger
	Name        postgres.ColumnString
	Description postgres.ColumnString
	Price       postgres.ColumnFloat
	Quantity    postgres.ColumnInteger
	Currency    postgres.ColumnString
	Image       postgres.ColumnString
	CreatedAt   postgres.ColumnTimestampz
	UpdatedAt   postgres.ColumnTimestampz

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

// creates new ProductsTable with assigned alias
func (a *ProductsTable) AS(alias string) *ProductsTable {
	aliasTable := newProductsTable()

	aliasTable.Table.AS(alias)

	return aliasTable
}

func newProductsTable() *ProductsTable {
	var (
		IDColumn          = postgres.IntegerColumn("id")
		NameColumn        = postgres.StringColumn("name")
		DescriptionColumn = postgres.StringColumn("description")
		PriceColumn       = postgres.FloatColumn("price")
		QuantityColumn    = postgres.IntegerColumn("quantity")
		CurrencyColumn    = postgres.StringColumn("currency")
		ImageColumn       = postgres.StringColumn("image")
		CreatedAtColumn   = postgres.TimestampzColumn("created_at")
		UpdatedAtColumn   = postgres.TimestampzColumn("updated_at")
	)

	return &ProductsTable{
		Table: postgres.NewTable("public", "products", IDColumn, NameColumn, DescriptionColumn, PriceColumn, QuantityColumn, CurrencyColumn, ImageColumn, CreatedAtColumn, UpdatedAtColumn),

		//Columns
		ID:          IDColumn,
		Name:        NameColumn,
		Description: DescriptionColumn,
		Price:       PriceColumn,
		Quantity:    QuantityColumn,
		Currency:    CurrencyColumn,
		Image:       ImageColumn,
		CreatedAt:   CreatedAtColumn,
		UpdatedAt:   UpdatedAtColumn,

		AllColumns:     postgres.ColumnList{IDColumn, NameColumn, DescriptionColumn, PriceColumn, QuantityColumn, CurrencyColumn, ImageColumn, CreatedAtColumn, UpdatedAtColumn},
		MutableColumns: postgres.ColumnList{NameColumn, DescriptionColumn, PriceColumn, QuantityColumn, CurrencyColumn, ImageColumn, CreatedAtColumn, UpdatedAtColumn},
	}
}

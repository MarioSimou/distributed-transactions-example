//
// Code generated by go-jet DO NOT EDIT.
// Generated at Sunday, 26-Apr-20 12:26:11 UTC
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package table

import (
	"github.com/go-jet/jet/postgres"
)

var Users = newUsersTable()

type UsersTable struct {
	postgres.Table

	//Columns
	ID        postgres.ColumnInteger
	Username  postgres.ColumnString
	Email     postgres.ColumnString
	Password  postgres.ColumnString
	Balance   postgres.ColumnFloat
	CreatedAt postgres.ColumnTimestampz
	UpdatedAt postgres.ColumnTimestampz

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

// creates new UsersTable with assigned alias
func (a *UsersTable) AS(alias string) *UsersTable {
	aliasTable := newUsersTable()

	aliasTable.Table.AS(alias)

	return aliasTable
}

func newUsersTable() *UsersTable {
	var (
		IDColumn        = postgres.IntegerColumn("id")
		UsernameColumn  = postgres.StringColumn("username")
		EmailColumn     = postgres.StringColumn("email")
		PasswordColumn  = postgres.StringColumn("password")
		BalanceColumn   = postgres.FloatColumn("balance")
		CreatedAtColumn = postgres.TimestampzColumn("created_at")
		UpdatedAtColumn = postgres.TimestampzColumn("updated_at")
	)

	return &UsersTable{
		Table: postgres.NewTable("public", "users", IDColumn, UsernameColumn, EmailColumn, PasswordColumn, BalanceColumn, CreatedAtColumn, UpdatedAtColumn),

		//Columns
		ID:        IDColumn,
		Username:  UsernameColumn,
		Email:     EmailColumn,
		Password:  PasswordColumn,
		Balance:   BalanceColumn,
		CreatedAt: CreatedAtColumn,
		UpdatedAt: UpdatedAtColumn,

		AllColumns:     postgres.ColumnList{IDColumn, UsernameColumn, EmailColumn, PasswordColumn, BalanceColumn, CreatedAtColumn, UpdatedAtColumn},
		MutableColumns: postgres.ColumnList{UsernameColumn, EmailColumn, PasswordColumn, BalanceColumn, CreatedAtColumn, UpdatedAtColumn},
	}
}

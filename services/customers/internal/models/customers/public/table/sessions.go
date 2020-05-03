//
// Code generated by go-jet DO NOT EDIT.
// Generated at Sunday, 03-May-20 15:51:16 UTC
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package table

import (
	"github.com/go-jet/jet/postgres"
)

var Sessions = newSessionsTable()

type SessionsTable struct {
	postgres.Table

	//Columns
	ID        postgres.ColumnInteger
	UserID    postgres.ColumnInteger
	GUID      postgres.ColumnString
	CreatedAt postgres.ColumnTimestampz
	ExpiresAt postgres.ColumnTimestampz

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

// creates new SessionsTable with assigned alias
func (a *SessionsTable) AS(alias string) *SessionsTable {
	aliasTable := newSessionsTable()

	aliasTable.Table.AS(alias)

	return aliasTable
}

func newSessionsTable() *SessionsTable {
	var (
		IDColumn        = postgres.IntegerColumn("id")
		UserIDColumn    = postgres.IntegerColumn("user_id")
		GUIDColumn      = postgres.StringColumn("guid")
		CreatedAtColumn = postgres.TimestampzColumn("created_at")
		ExpiresAtColumn = postgres.TimestampzColumn("expires_at")
	)

	return &SessionsTable{
		Table: postgres.NewTable("public", "sessions", IDColumn, UserIDColumn, GUIDColumn, CreatedAtColumn, ExpiresAtColumn),

		//Columns
		ID:        IDColumn,
		UserID:    UserIDColumn,
		GUID:      GUIDColumn,
		CreatedAt: CreatedAtColumn,
		ExpiresAt: ExpiresAtColumn,

		AllColumns:     postgres.ColumnList{IDColumn, UserIDColumn, GUIDColumn, CreatedAtColumn, ExpiresAtColumn},
		MutableColumns: postgres.ColumnList{UserIDColumn, GUIDColumn, CreatedAtColumn, ExpiresAtColumn},
	}
}

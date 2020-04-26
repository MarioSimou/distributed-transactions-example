//
// Code generated by go-jet DO NOT EDIT.
// Generated at Sunday, 26-Apr-20 19:29:36 UTC
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package table

import (
	"github.com/go-jet/jet/postgres"
)

var SchemaMigrations = newSchemaMigrationsTable()

type SchemaMigrationsTable struct {
	postgres.Table

	//Columns
	Version postgres.ColumnInteger
	Dirty   postgres.ColumnBool

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

// creates new SchemaMigrationsTable with assigned alias
func (a *SchemaMigrationsTable) AS(alias string) *SchemaMigrationsTable {
	aliasTable := newSchemaMigrationsTable()

	aliasTable.Table.AS(alias)

	return aliasTable
}

func newSchemaMigrationsTable() *SchemaMigrationsTable {
	var (
		VersionColumn = postgres.IntegerColumn("version")
		DirtyColumn   = postgres.BoolColumn("dirty")
	)

	return &SchemaMigrationsTable{
		Table: postgres.NewTable("public", "schema_migrations", VersionColumn, DirtyColumn),

		//Columns
		Version: VersionColumn,
		Dirty:   DirtyColumn,

		AllColumns:     postgres.ColumnList{VersionColumn, DirtyColumn},
		MutableColumns: postgres.ColumnList{DirtyColumn},
	}
}

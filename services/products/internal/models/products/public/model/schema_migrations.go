//
// Code generated by go-jet DO NOT EDIT.
// Generated at Saturday, 02-May-20 15:49:41 UTC
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package model

type SchemaMigrations struct {
	Version int64 `sql:"primary_key"`
	Dirty   bool
}

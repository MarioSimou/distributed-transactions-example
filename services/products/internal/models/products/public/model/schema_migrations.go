//
// Code generated by go-jet DO NOT EDIT.
// Generated at Tuesday, 05-May-20 21:51:13 UTC
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package model

type SchemaMigrations struct {
	Version int64 `sql:"primary_key"`
	Dirty   bool
}

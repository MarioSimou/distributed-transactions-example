//
// Code generated by go-jet DO NOT EDIT.
// Generated at Monday, 27-Apr-20 18:42:10 UTC
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package model

import (
	"time"
)

type Products struct {
	ID        int32 `sql:"primary_key"`
	Name      string
	Price     float64
	Quantity  *int32
	Currency  *Currency
	CreatedAt time.Time
	UpdatedAt time.Time
}

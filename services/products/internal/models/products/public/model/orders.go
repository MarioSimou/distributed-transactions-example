//
// Code generated by go-jet DO NOT EDIT.
// Generated at Saturday, 02-May-20 09:33:59 UTC
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package model

import (
	"time"
)

type Orders struct {
	ID        int32 `sql:"primary_key"`
	UID       string
	ProductID int32
	Quantity  int32
	Total     float64
	UserID    int32
	Status    *OrderStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}
//
// Code generated by go-jet DO NOT EDIT.
// Generated at Tuesday, 05-May-20 22:20:54 UTC
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package model

import (
	"time"
)

type Users struct {
	ID        int32 `json:"id" sql:"primary_key"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Balance   *float64 `json:"balance"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt"`
}

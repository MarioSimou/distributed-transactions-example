//
// Code generated by go-jet DO NOT EDIT.
// Generated at Sunday, 03-May-20 12:39:38 UTC
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package model

import (
	"time"
)

type Products struct {
	ID          int32 `json:"id" sql:"primary_key"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       float64 `json:"price"`
	Quantity    *int32 `json:"quantity"`
	Currency    *Currency `json:"currency"`
	Image       string `json:"image"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
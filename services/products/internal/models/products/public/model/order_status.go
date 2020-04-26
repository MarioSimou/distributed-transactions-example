//
// Code generated by go-jet DO NOT EDIT.
// Generated at Sunday, 26-Apr-20 10:14:25 UTC
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package model

import "errors"

type OrderStatus string

const (
	OrderStatus_Pending  OrderStatus = "pending"
	OrderStatus_Approved OrderStatus = "approved"
	OrderStatus_Declined OrderStatus = "declined"
)

func (e *OrderStatus) Scan(value interface{}) error {
	if v, ok := value.(string); !ok {
		return errors.New("jet: Invalid data for OrderStatus enum")
	} else {
		switch string(v) {
		case "pending":
			*e = OrderStatus_Pending
		case "approved":
			*e = OrderStatus_Approved
		case "declined":
			*e = OrderStatus_Declined
		default:
			return errors.New("jet: Inavlid data " + string(v) + "for OrderStatus enum")
		}

		return nil
	}
}

func (e OrderStatus) String() string {
	return string(e)
}

//
// Code generated by go-jet DO NOT EDIT.
// Generated at Monday, 27-Apr-20 18:42:10 UTC
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package model

import "errors"

type Currency string

const (
	Currency_Gbp  Currency = "GBP"
	Currency_Euro Currency = "EURO"
	Currency_Usd  Currency = "USD"
)

func (e *Currency) Scan(value interface{}) error {
	if v, ok := value.(string); !ok {
		return errors.New("jet: Invalid data for Currency enum")
	} else {
		switch string(v) {
		case "GBP":
			*e = Currency_Gbp
		case "EURO":
			*e = Currency_Euro
		case "USD":
			*e = Currency_Usd
		default:
			return errors.New("jet: Inavlid data " + string(v) + "for Currency enum")
		}

		return nil
	}
}

func (e Currency) String() string {
	return string(e)
}

//
// Code generated by go-jet DO NOT EDIT.
// Generated at Saturday, 02-May-20 09:33:59 UTC
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package enum

import "github.com/go-jet/jet/postgres"

var Currency = &struct {
	Gbp  postgres.StringExpression
	Euro postgres.StringExpression
	Usd  postgres.StringExpression
}{
	Gbp:  postgres.NewEnumValue("GBP"),
	Euro: postgres.NewEnumValue("EURO"),
	Usd:  postgres.NewEnumValue("USD"),
}
//
// Code generated by go-jet DO NOT EDIT.
// Generated at Monday, 27-Apr-20 21:13:20 UTC
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

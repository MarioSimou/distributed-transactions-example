//
// Code generated by go-jet DO NOT EDIT.
// Generated at Sunday, 03-May-20 17:43:30 UTC
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package enum

import "github.com/go-jet/jet/postgres"

var OrderStatus = &struct {
	Pending  postgres.StringExpression
	Approved postgres.StringExpression
	Declined postgres.StringExpression
}{
	Pending:  postgres.NewEnumValue("pending"),
	Approved: postgres.NewEnumValue("approved"),
	Declined: postgres.NewEnumValue("declined"),
}

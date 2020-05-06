package subscribers

import (
	"context"
	r "products/internal/rabbitmq"
)

func GetSubscribers(ctx context.Context) r.Subscribers {
	return r.NewSubscribers(
		r.NewSubscriber(
			"customers_charge_customer_failure",
			customersChargeCustomerFailure,
			ctx,
		),
		r.NewSubscriber(
			"customers_charge_customer_success",
			customersChargeCustomerSuccess,
			ctx,
		),
	)
}
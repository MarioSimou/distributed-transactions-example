package subscribers

import (
	"context"
	r "customers/internal/rabbitmq"
)

func GetSubscribers(ctx context.Context) r.Subscribers {
	return r.NewSubscribers(
		r.NewSubscriber(
			"products_create_order_success",
			productsCreateOrderSuccess,
			ctx,
		),
	)
}

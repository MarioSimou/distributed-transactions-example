package main

import (
	"context"
	"database/sql"
	"log"
	"os"

	i "products/internal"
	r "products/internal/rabbitmq"
	s "products/internal/subscribers"
	"sync"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main(){
	var e error
	var db *sql.DB
	var wg sync.WaitGroup
	var publisher r.PublisherInterface
	var queueConn = &r.ConnectionStruct{}
	var subResChan chan r.SubscriptionResponse
	var env = i.EnvVariables{
		DBUri: os.Getenv("DB_URI"),
		Port: os.Getenv("PORT"),
		QueueUri: os.Getenv("QUEUE_URI"),
		QueuesNames: []string{
			"products_create_order_success",
		},
	}
	var globalMiddlewares = []gin.HandlerFunc{
		i.HandleCORS,
	}
	
	if db , e = sql.Open("postgres", env.DBUri); e != nil {
		i.HandleError(e)
	}
	defer db.Close()
	
	if e := db.Ping(); e != nil {
		i.HandleError(e)
	}
	
	if e := queueConn.Start(env.QueueUri); e != nil {
		log.Fatalln(e)
	}
	defer queueConn.Close()

	if publisher, e = r.NewPublisher(env.QueuesNames, queueConn); e != nil {
		log.Fatalln(e)
	}
	
	var bg = context.Background()
	var withDB = context.WithValue(bg,"DB", db)
	var withPublisher = context.WithValue(withDB, "Publisher", publisher)
	if subResChan, e = r.NewSubscription(s.GetSubscribers(withPublisher), queueConn); e != nil {
		log.Fatalln(e)
	}
	
	var contr = i.Controller{
		Env: env,
		DB: db,
		Publisher: publisher,
	}

	var routes = []i.Route{
		i.Route{
			Method: "GET",
			Path: "/ping",
			HandlerFunc: contr.Ping,
		},
		i.Route{
			Method: "GET",
			Path: "/products",
			HandlerFunc: contr.GetProducts,
		},
		i.Route{
			Method: "GET",
			Path: "/products/:id",
			HandlerFunc: contr.GetProduct,
		},
		i.Route{
			Method: "POST",
			Path: "/products",
			HandlerFunc: contr.CreateProduct,
		},
		i.Route{
			Method: "PUT",
			Path: "/products/:id",
			HandlerFunc: contr.UpdateProduct,
		},
		i.Route{
			Method: "DELETE",
			Path: "/products/:id",
			HandlerFunc: contr.DeleteProduct,
		},
		i.Route{
			Method: "GET",
			Path: "/orders",
			HandlerFunc: contr.GetOrders,
		},
		i.Route{
			Method: "GET",
			Path: "/orders/:id",
			HandlerFunc: contr.GetOrder,
		},
		i.Route{
			Method: "POST",
			Path: "/orders",
			HandlerFunc: contr.CreateOrder,
		},
		i.Route{
			Method: "DELETE",
			Path: "/orders/:id",
			HandlerFunc: contr.DeleteOrder,
		},
	}	
	
	wg.Add(2)
	go i.LaunchServer(routes, globalMiddlewares, env, &wg)
	go i.HandleSubscribersResponses(subResChan, &wg)
	wg.Wait()
}
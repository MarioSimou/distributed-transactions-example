package main

import (
	"database/sql"
	"log"
	"os"
	"sync"

	"context"
	c "customers/internal"
	r "customers/internal/rabbitmq"
	s "customers/internal/subscribers"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main(){
	var envVariables = c.EnvironmentVariables{
		Port: "3000",
		UIDomain: os.Getenv("UI_DOMAIN"),
		QueueUri: os.Getenv("QUEUE_URI"),
		QueuesNames: []string{
			"customers_charge_customer_success",
			"customers_charge_customer_failure",
		},
	}
	var wg sync.WaitGroup
	var publisher r.PublisherInterface
	var subResChan chan r.SubscriptionResponse
	var conn = &r.ConnectionStruct{}

	var db, e = sql.Open("postgres", os.Getenv("DB_URI"))
	if e != nil {
		log.Fatalf("Error: %v\n", e)
	}
	if e := db.Ping(); e != nil {
		log.Fatalf("Error: %v\n", e)
	}
	if e := conn.Start(envVariables.QueueUri); e != nil {
		log.Fatalf("Error: %v\n", e)
	}
	if publisher, e = r.NewPublisher(envVariables.QueuesNames, conn); e != nil {
		log.Fatalf("Error: %v\n", e)
	}

	var parent = context.Background()
	var ctxDB = context.WithValue(parent, "DB", db)
	var ctx = context.WithValue(ctxDB, "Publisher", publisher)
	if subResChan, e = r.NewSubscription(s.GetSubscribers(ctx), conn); e != nil {
		log.Fatalf("Error: %v\n", e)
	}

	var controller = c.Controller{
		EnvVariables: envVariables,
		DB: db,
		Publisher: publisher,
	}
	var globMiddlewares = []gin.HandlerFunc{
		c.HandleCORS,
	}
	var routes = []c.Route{
		c.Route{
			HttpMethod: "GET",
			Path: "/ping",
			HandlerFunc: controller.Ping,
		},
		c.Route{
			HttpMethod: "GET",
			Path: "/users/:id",
			HandlerFunc: controller.GetUser,
		},
		c.Route{
			HttpMethod: "GET",
			Path: "/users",
			HandlerFunc: controller.GetUsers,
		},
		c.Route{
			HttpMethod: "POST",
			Path: "/users",
			HandlerFunc: controller.CreateUser,
		},
		c.Route{
			HttpMethod: "DELETE",
			Path: "/users/:id",
			HandlerFunc: controller.DeleteUser,
		},
		c.Route{
			HttpMethod: "PUT",
			Path: "/users/:id",
			HandlerFunc: controller.UpdateUser,
		},
		c.Route{
			HttpMethod: "POST",
			Path: "/signin",
			HandlerFunc: controller.SignInUser,
		},
		c.Route{
			HttpMethod: "GET",
			Path: "/signin",
			HandlerFunc: controller.SignInUserWithCookie,
		},
		c.Route{
			HttpMethod: "POST",
			Path:"/logout",
			HandlerFunc: controller.LogOut,
		},
	}

	wg.Add(2)
	go c.LaunchServer(routes, globMiddlewares, envVariables, &wg)
	go c.HandleSubscribersResponses(subResChan,&wg)
	wg.Wait()
}
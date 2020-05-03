package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	i "products/internal"
	"sync"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main(){
	var e error
	var db *sql.DB
	var wg sync.WaitGroup
	var publisher i.Publishing
	var subResChan chan i.SubResponse
	var env = i.EnvVariables{
		DBUri: os.Getenv("DB_URI"),
		Port: os.Getenv("PORT"),
		QueueUri: os.Getenv("QUEUE_URI"),
	}
	var globalMiddlewares = []gin.HandlerFunc{
		i.HandleCORS,
	}
	var queuesNames = []string{
		"products_created_order_success",
	}
	
	if db , e = sql.Open("postgres", env.DBUri); e != nil {
		i.HandleError(e)
	}
	defer db.Close()
	if e := db.Ping(); e != nil {
		i.HandleError(e)
	}
	if publisher, e = i.NewPublisher(env.QueueUri,queuesNames); e != nil {
		log.Fatalln(e)
	}
	var subscribers = []i.Subscriber{
		i.Subscriber{
			QueueName: "products_created_order_success",
			HandlerFunc: func(d i.Message) error {
				fmt.Println("HELLO WORLD")
				return nil
			},
		},
	}
	if subResChan, e = i.NewSubscription(env.QueueUri, subscribers); e != nil {
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
	go i.LaunchServer(routes, globalMiddlewares, env, wg)
	go i.HandleSubscribersResponses(subResChan, wg)
	wg.Wait()
}
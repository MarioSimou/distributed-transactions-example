package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	i "products/internal"
	// "time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main(){
	var e error
	var db *sql.DB
	var env = i.EnvVariables{
		DBUri: os.Getenv("DB_URI"),
		Port: os.Getenv("PORT"),
	}

	if db , e = sql.Open("postgres", env.DBUri); e != nil {
		i.HandleError(e)
	}
	defer db.Close()
	if e := db.Ping(); e != nil {
		i.HandleError(e)
	}
	var contr = i.Controller{
		Env: env,
		DB: db,
	}
	var handleCors = cors.New(cors.Config{
		AllowOrigins:     []string{os.Getenv("ALLOW_ORIGIN_DOMAIN")},
		AllowMethods:     []string{"OPTIONS","GET","POST","DELETE","PUT", "PATCH"},
		AllowHeaders:     []string{"Content-Type: application/json"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		// MaxAge: 12 * time.Hour,
	})

	var globalMiddlewares = []gin.HandlerFunc{
		handleCors,
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
	var router = i.GetRouter(routes, globalMiddlewares, env)

	log.Fatalln(router.Run(fmt.Sprintf(":%s", env.Port)))
}
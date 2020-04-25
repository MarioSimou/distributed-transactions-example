package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	i "products/internal"

	_ "github.com/lib/pq"
)

func handleError(e error){
	if e != nil {
		log.Fatalf("Error: %v\n", e)
	}
}


func main(){
	var e error
	var db *sql.DB
	var env = i.EnvVariables{
		DBUri: os.Getenv("DB_URI"),
		Port: os.Getenv("PORT"),
	}

	if db , e = sql.Open("postgres", env.DBUri); e != nil {
		handleError(e)
	}
	defer db.Close()
	if e := db.Ping(); e != nil {
		handleError(e)
	}
	var contr = i.Controllers{
		Env: env,
		DB: db,
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
	}	
	var router = i.GetRouter(routes)

	log.Fatalln(router.Run(fmt.Sprintf(":%s", env.Port)))
}
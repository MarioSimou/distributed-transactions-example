package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	c "customers/internal"

	_ "github.com/lib/pq"
)

func main(){
	var envVariables = c.EnvironmentVariables{
		Port: "3000",
	}

	var server = c.Server{
		EnvVariables: envVariables,
	}
	var db, e = sql.Open("postgres", os.Getenv("DB_URI"))
	if e != nil {
		log.Fatalf("Error: %v", e)
	}
	if e := db.Ping(); e != nil {
		log.Fatalf("Error: %v", e)
	}
	var controller = c.Controller{
		EnvVariables: envVariables,
		DB: db,
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
			Path: "users",
			HandlerFunc: controller.CreateUser,
		},
		c.Route{
			HttpMethod: "DELETE",
			Path: "users/:id",
			HandlerFunc: controller.DeleteUser,
		},
		c.Route{
			HttpMethod: "PUT",
			Path: "users/:id",
			HandlerFunc: controller.UpdateUser,
		},
	}
	var router = server.Setup(routes)

	router.Run(fmt.Sprintf(":%s", envVariables.Port))
}
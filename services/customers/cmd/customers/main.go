package main

import (
	"database/sql"
	"log"
	"os"

	c "customers/internal"
	_ "github.com/lib/pq"
)

func main(){
	var envVariables = c.EnvironmentVariables{
		Port: ":3000",
	}

	var server = c.Server{
		EnvVariables: envVariables,
	}

	var controller = c.Controller{
		EnvVariables: envVariables,
	}

	var db, e = sql.Open("postgres", os.Getenv("DB_URI"))
	if e != nil {
		log.Fatalf("Error: %v", e)
	}
	if e := db.Ping(); e != nil {
		log.Fatalf("Error: %v", e)
	}

	var routes = []c.Route{
		c.Route{
			HttpMethod: "GET",
			Path: "/ping",
			HandlerFunc: controller.Ping,
		},
	}
	var router = server.Setup(routes)

	router.Run(envVariables.Port)
}
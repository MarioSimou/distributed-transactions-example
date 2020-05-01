package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	c "customers/internal"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main(){
	var envVariables = c.EnvironmentVariables{
		Port: "3000",
		UIDomain: os.Getenv("UI_DOMAIN"),
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
	var router = server.Setup(routes, globMiddlewares)
	router.Run(fmt.Sprintf(":%s", envVariables.Port))
}
package main

import (
	"log"
	"net/http"

	"calendly_adventures/config"
	"calendly_adventures/db"
	"calendly_adventures/routes"
)

func main() {
	cfg := config.LoadConfig()
	db.ConnectDB(cfg.DSN)

	router := routes.SetupRoutes()
	log.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

/*
Assumptions:
1. Lets say user is active all days and 9am to 5pm. Trying to schedule in this time period
*/

package main

import (
	"context"
	"fmt"
	"healthcheck/config"
	"healthcheck/db"
	"healthcheck/route"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func Init() (*mux.Router, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Printf("Error loading environment file: %v\n", err)
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := db.ConnectMongo(ctx, cfg.MongoURI)
	if err != nil {
		log.Printf("Failed to connect to MongoDB: %v", err)
		return nil, err
	}
	r := mux.NewRouter()
	route.RegisterHealthCheckRoutes(r, client, cfg)

	return r, nil
}

func main() {
	r, err := Init()
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(http.ListenAndServe(":8080", r))
	fmt.Println("Server running on port 8080")

}

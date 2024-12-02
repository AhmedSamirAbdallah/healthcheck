package cmd

import (
	"fmt"
	"healthcheck/config"
	"healthcheck/db"
	"healthcheck/route"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func Init() (*mux.Router, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Printf("Error loading environment file: %v\n", err)
	}
	client, err := db.ConnectMongo(cfg.MongoURI)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
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

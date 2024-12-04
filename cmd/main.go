package main

import (
	"fmt"
	"healthcheck/config"
	"healthcheck/route"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func Init() (*mux.Router, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Printf("Error loading environment file: %v\n", err)
		return nil, err
	}

	r := mux.NewRouter()
	route.RegisterHealthCheckRoutes(r, cfg)

	return r, nil
}

func main() {
	r, err := Init()
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(http.ListenAndServe(":8020", r))
	fmt.Println("Server running on port 8020")

}

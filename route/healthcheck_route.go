package route

import (
	"healthcheck/config"
	"healthcheck/handler"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterHealthCheckRoutes(mux *mux.Router, client *mongo.Client, config *config.Config) {
	mux.HandleFunc("/api/health-check", handler.HealthCheckHandler(client, config)).Methods("GET")
}

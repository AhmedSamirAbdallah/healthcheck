package route

import (
	"healthcheck/config"
	"healthcheck/handler"

	"github.com/gorilla/mux"
)

func RegisterHealthCheckRoutes(mux *mux.Router, config *config.Config) {
	mux.HandleFunc("/api/health-check", handler.HealthCheckHandler(config)).Methods("GET")
}

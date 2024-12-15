package route

import (
	"healthcheck/handler"
	"healthcheck/internal/config"

	"github.com/gorilla/mux"
)

func RegisterHealthCheckRoutes(mux *mux.Router, config *config.Config) {
	mux.HandleFunc("/api/health-check", handler.HealthCheckHandler(config)).Methods("GET")
}

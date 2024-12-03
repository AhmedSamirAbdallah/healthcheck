package handler

import (
	"encoding/json"
	"healthcheck/config"
	"healthcheck/db"
	"healthcheck/kafka"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type HealthCheckResponse struct {
	Status       string                 `json:"status"`
	UpTime       string                 `json:"upTime"`
	Dependancies map[string]interface{} `json:"dependancies"`
}

func HealthCheckHandler(client *mongo.Client, config *config.Config) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		upTime := time.Now().String()

		dbStatus := map[string]interface{}{
			"connection": db.CheckDatabase(client),
			"read":       db.CheckReadOnDB(client, config.DatabaseName),
			"write":      db.CheckWriteOnDB(client, config.DatabaseName),
		}

		kafkaStatus := map[string]interface{}{
			"connection": kafka.CheckKafka(config.KafkaBroker),
			"produce":    kafka.CheckProduce(config.KafkaBroker, config.KafkaTopic),
			"consume":    kafka.CheckConsume(config.KafkaBroker, config.KafkaTopic, config.KafkaGroupID),
		}

		response := HealthCheckResponse{
			Status: "UP",
			UpTime: upTime,
			Dependancies: map[string]interface{}{
				"database": dbStatus,
				"kafka":    kafkaStatus,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}
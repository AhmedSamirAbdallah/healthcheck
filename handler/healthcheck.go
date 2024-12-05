package handler

import (
	"encoding/json"
	"healthcheck/config"
	"healthcheck/db"
	"healthcheck/kafka"
	"healthcheck/redis"
	"healthcheck/temporal"
	"net/http"
	"time"
)

type HealthCheckResponse struct {
	Status       string                 `json:"status"`
	UpTime       string                 `json:"upTime"`
	Dependancies map[string]interface{} `json:"dependancies"`
}

func HealthCheckHandler(cfg *config.Config) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		upTime := time.Now().String()

		response := HealthCheckResponse{
			Status:       "UP",
			UpTime:       upTime,
			Dependancies: map[string]interface{}{},
		}

		queryParam := r.URL.Query()

		checkDatabase := queryParam.Get("database") != ""
		if checkDatabase {
			db.InitDB(cfg.MongoURI)
			dbStatus := map[string]interface{}{
				"connection": db.CheckDatabase(),
				"read":       db.CheckReadOnDB(cfg.DatabaseName),
				"write":      db.CheckWriteOnDB(cfg.DatabaseName),
			}
			response.Dependancies["database"] = dbStatus
		}

		checkKafka := queryParam.Get("kafka") != ""
		if checkKafka {
			kafka.InitKafka(cfg.KafkaBroker, cfg.KafkaGroupID)
			kafkaStatus := map[string]interface{}{
				"connection": kafka.CheckKafka(),
				"produce":    kafka.CheckProduce(cfg.KafkaTopic),
				// "consume":    kafka.CheckConsume(config.KafkaTopic),
			}
			response.Dependancies["kafka"] = kafkaStatus
		}

		checkRedis := queryParam.Get("redis") != ""
		if checkRedis {
			redis.InitRedis(cfg.RedisHost+":"+cfg.RedisPort, cfg.RedisPassword, cfg.RedisDB)
			redisStatus := map[string]interface{}{
				"connection": redis.CheckRedisConnection(),
				"write":      redis.CheckWriteOnRedis("healthcheck", "healthy"),
				"read":       redis.CheckReadOnRedis("healthcheck"),
			}
			response.Dependancies["redis"] = redisStatus
		}

		checkTemporal := queryParam.Get("temporal") != ""
		if checkTemporal {
			temporalStatus := map[string]interface{}{
				"connection": temporal.CheckTemporalConnection(cfg.TemporalUrl, cfg.WithTLS),
			}
			response.Dependancies["temporal"] = temporalStatus
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

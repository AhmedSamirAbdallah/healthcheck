package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	ServiceName   string
	MongoURI      string
	DatabaseName  string
	KafkaBroker   string
	KafkaTopic    string
	KafkaGroupID  string
	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       int
	TemporalUrl   string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load("./.env")
	if err != nil {
		log.Printf("Error loading environment file: %v\n", err)
	}
	redisDB, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		log.Printf("Invalid value for REDIS_DB: %v. Using default value %d.\n", err, redisDB)
		redisDB = 0
	}

	return &Config{
		ServiceName:   os.Getenv("SERVICE_NAME"),
		MongoURI:      os.Getenv("MONGO_URI"),
		DatabaseName:  os.Getenv("DATABASE_NAME"),
		KafkaBroker:   os.Getenv("KAFKA_BROKER"),
		KafkaTopic:    os.Getenv("KAFKA_TOPIC"),
		KafkaGroupID:  os.Getenv("KAFKA_GROUP_ID"),
		RedisHost:     os.Getenv("REDIS_HOST"),
		RedisPort:     os.Getenv("REDIS_PORT"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
		RedisDB:       redisDB,
		TemporalUrl:   os.Getenv("TEMPORAL_URL"),
	}, nil
}

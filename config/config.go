package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServiceName  string
	MongoURI     string
	DatabaseName string
	KafkaBroker  string
	KafkaTopic   string
	KafkaGroupID string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load("./.env")
	if err != nil {
		log.Printf("Error loading environment file: %v\n", err)
	}
	return &Config{
		ServiceName:  os.Getenv("SERVICE_NAME"),
		MongoURI:     os.Getenv("MONGO_URI"),
		DatabaseName: os.Getenv("DATABASE_NAME"),
		KafkaBroker:  os.Getenv("KAFKA_BROKER"),
		KafkaTopic:   os.Getenv("KAFKA_TOPIC"),
		KafkaGroupID: os.Getenv("KAFKA_GROUP_ID"),
	}, nil
}
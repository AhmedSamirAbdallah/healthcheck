package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type HealthCheck struct {
	ID        string `bson:"_id,omitempty" json:"id"`
	Service   string `bson:"service" json:"service"`
	Timestamp string `bson:"timestamp" json:"timestamp"`
}

func CheckDatabase(client *mongo.Client) bool {
	err := client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Printf("Error pinging MongoDB: %v", err)
		return false
	}
	return true
}

func CheckReadOnDB(client *mongo.Client, databaseName string) bool {
	collection := client.Database(databaseName).Collection("healthcheck")
	var result HealthCheck
	err := collection.FindOne(context.Background(), map[string]interface{}{}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Println("Healthcheck collection is empty, but the database is accessible.")
			return true
		}
		log.Printf("Error reading from healthcheck collection: %v", err)
		return false
	}
	log.Printf("%v", result)
	return true
}

func CheckWriteOnDB(client *mongo.Client, databaseName string) bool {
	collection := client.Database(databaseName).Collection("healthcheck")
	healthCheck := HealthCheck{
		Service:   "healthcheck",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	_, err := collection.InsertOne(context.Background(), healthCheck)
	if err != nil {
		log.Printf("Error writing to healthcheck collection: %v", err)
		return false
	}
	return true
}

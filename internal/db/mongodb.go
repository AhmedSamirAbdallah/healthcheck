package db

import (
	"context"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	sharedClient *mongo.Client
	initOnce     sync.Once
)

type HealthCheck struct {
	ID        string `bson:"_id,omitempty" json:"id"`
	Service   string `bson:"service" json:"service"`
	Timestamp string `bson:"timestamp" json:"timestamp"`
}

func InitDB(MongoURI string) {
	initOnce.Do(func() {
		cliebtOptions := options.Client().ApplyURI(MongoURI)
		client, err := mongo.Connect(context.Background(), cliebtOptions)
		if err != nil {
			log.Printf("Error connecting to MongoDB: %v", err)
			return
		}
		err = client.Ping(context.Background(), readpref.Primary())
		if err != nil {
			log.Printf("Error pinging MongoDB: %v", err)
			return
		}
		sharedClient = client
		log.Println("DBConnection : MongoDB connection successful.")
	})
}

func CheckDatabase() bool {
	err := sharedClient.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Printf("Error pinging MongoDB: %v", err)
		return false
	}
	return true
}

func CheckReadOnDB(databaseName string) bool {
	collection := sharedClient.Database(databaseName).Collection("healthcheck")
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
	log.Printf("DBRead : %v", result)
	return true
}

func CheckWriteOnDB(databaseName string) bool {
	collection := sharedClient.Database(databaseName).Collection("healthcheck")

	var lastRecord HealthCheck

	filter := map[string]interface{}{}
	options := options.FindOne().SetSort(map[string]interface{}{"timestamp": -1})
	err := collection.FindOne(context.Background(), filter, options).Decode(&lastRecord)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return InsertFirstHealthCheckRecord(collection)
		}
		log.Printf("Error finding last record: %v", err)
		return false
	}

	update := map[string]interface{}{
		"$set": map[string]interface{}{
			"timestamp": time.Now().String(),
		},
	}

	objectID, err := primitive.ObjectIDFromHex(lastRecord.ID)

	if err != nil {
		log.Printf("Error converting ID to ObjectID: %v", err)
		return false
	}

	filter = map[string]interface{}{
		"_id": objectID,
	}

	_, err = collection.UpdateOne(context.Background(), filter, update)

	if err != nil {
		log.Printf("Error updating the last inserted record: %v", err)
		return false
	}

	log.Println("DBWrite : Updated the timestamp of the last inserted healthcheck record.")
	return true
}

func InsertFirstHealthCheckRecord(collection *mongo.Collection) bool {
	healthCheck := HealthCheck{
		Service:   "healthcheck",
		Timestamp: time.Now().String(),
	}
	_, err := collection.InsertOne(context.Background(), healthCheck)
	if err != nil {
		log.Printf("Error inserting first healthcheck record: %v", err)
		return false
	}
	log.Println("DBWrite : Inserted first healthcheck record.")
	return true

}

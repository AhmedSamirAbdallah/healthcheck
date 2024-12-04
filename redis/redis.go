package redis

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	redisClient *redis.Client
	initOnce    sync.Once
)

func InitRedis(serverAddress string, password string, redisDB int) {
	initOnce.Do(func() {
		redisOption := &redis.Options{
			Addr:     serverAddress,
			Password: password,
			DB:       redisDB,
		}
		client := redis.NewClient(redisOption)
		err := client.Ping(context.Background()).Err()
		if err != nil {
			log.Printf("Error connecting to Redis: %v", err)
			return
		}
		redisClient = client
		log.Println("RedisConnection : Redis connection successful.")
	})
}

func CheckRedisConnection() bool {
	err := redisClient.Ping(context.Background()).Err()
	if err != nil {
		log.Printf("Error pinging Redis: %v", err)
		return false
	}
	return true
}

func CheckWriteOnRedis(key string, value string) bool {
	err := redisClient.Set(context.Background(), key, value, 10*time.Second).Err()
	if err != nil {
		log.Printf("Error writing to Redis: %v", err)
		return false
	}
	log.Printf("RedisWrite : Successfully wrote to Redis: %s = %s", key, value)
	return true
}

func CheckReadOnRedis(key string) bool {
	val, err := redisClient.Get(context.Background(), key).Result()
	if err != nil {
		log.Printf("Error reading from Redis: %v", err)
		return false
	}
	log.Printf("RedisRead : Read value from Redis: %s = %s", key, val)
	return true
}

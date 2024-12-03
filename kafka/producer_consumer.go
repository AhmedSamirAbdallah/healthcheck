package kafka

import (
	"context"
	"log"
)

func CheckProduce(brokers string, topic string) bool {
	producer, err := CreateProducer(brokers)
	if err != nil {
		log.Printf("Error creating Kafka producer: %v", err)
		return false
	}
	defer producer.Close()
	err = SendMessage(producer, topic, "produce within the health check")
	if err != nil {
		log.Printf("Error sending message to Kafka topic %s: %v", topic, err)
		return false
	}
	return true
}

func CheckConsume(brokers string, topic string, groupID string) bool {
	consumerGroup, err := CreateConsumer(brokers, groupID)
	if err != nil {
		log.Printf("Error creating Kafka consumer group: %v", err)
		return false
	}
	defer consumerGroup.Close()

	consumer := &Consumer{
		ready:    make(chan bool),
		received: make(chan bool),
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Run consumer in a separate goroutine
	go func() {
		for {
			err := consumerGroup.Consume(ctx, []string{topic}, consumer)
			if err != nil {
				log.Printf("Error during message consumption: %v", err)
				return
			}
		}
	}()

	log.Println("Kafka consumer is ready for health check")

	// Wait briefly to allow a message to be consumed or timeout
	select {
	case <-consumer.received:
		log.Println("Consumer received at least one message")
		return true
	case <-ctx.Done():
		log.Println("Context done without consuming messages")
		return false
	}
}

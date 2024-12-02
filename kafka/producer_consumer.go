package kafka

import "log"

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

func CheckConsume(brokers string, topic string) bool {
	return true
}

package kafka

import (
	"log"
	"strings"
	"sync"

	"github.com/IBM/sarama"
)

var (
	sharedProducer   sarama.SyncProducer
	sharedConsumer   sarama.Consumer
	sharedKafkaAdmin sarama.ClusterAdmin
	initOnce         sync.Once
)

func InitKafka(brokers string, topic string) error {
	var err error
	initOnce.Do(func() {
		sharedProducer, err = CreateProducer(brokers)
		if err != nil {
			log.Printf("Failed to initialize Kafka producer: %v", err)
			return
		}

		sharedConsumer, err = CreateConsumer(brokers)
		if err != nil {
			log.Printf("Failed to initialize Kafka consumer: %v", err)
			return
		}

		brokerList := strings.Split(brokers, ",")
		config := sarama.NewConfig()
		config.Version = sarama.V2_0_0_0
		config.ClientID = "health-check-client"

		sharedKafkaAdmin, err = sarama.NewClusterAdmin(brokerList, config)
		if err != nil {
			log.Printf("Failed to initialize Kafka admin client: %v", err)
			return
		}
	})
	return err
}

func CreateProducer(brokers string) (sarama.SyncProducer, error) {
	if sharedProducer != nil {
		log.Println("Reusing shared Kafka producer")
		return sharedProducer, nil
	}

	brokerList := strings.Split(brokers, ",")

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokerList, config)
	if err != nil {
		return nil, err
	}

	return producer, nil
}

func SendMessage(producer sarama.SyncProducer, topic, message string) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}

	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		return err
	}

	log.Printf("KafkaProduce : Message sent to partition %d with offset %d\n", partition, offset)
	return nil
}

func CreateConsumer(brokers string) (sarama.Consumer, error) {
	brokerList := strings.Split(brokers, ",")

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer(brokerList, config)
	if err != nil {
		return nil, err
	}
	return consumer, nil
}

func CheckProduce(topic string) bool {
	if sharedProducer == nil {
		log.Println("Kafka producer is not initialized")
		return false
	}

	err := SendMessage(sharedProducer, topic, "produce within the health check")
	if err != nil {
		log.Printf("Error sending message to Kafka topic %s: %v", topic, err)
		return false
	}
	return true
}
func CheckConsumer(topic string) bool {
	if sharedConsumer == nil {
		log.Println("Kafka consumer is not initialized")
		return false
	}
	resultChan := make(chan bool)
	go func() {
		consumer, err := sharedConsumer.ConsumePartition(topic, 0, sarama.OffsetOldest)
		if err != nil {
			log.Printf("Error : %v", err)
			resultChan <- false
			return
		}
		defer consumer.Close()
		select {
		case err := <-consumer.Errors():
			log.Printf("Kafka consumer error: %v", err)
			resultChan <- false
		case message := <-consumer.Messages():
			log.Printf("KafkaConsumer: Successfully consumed message: topic=%s partition=%d offset=%d key=%s value=%s",
				message.Topic, message.Partition, message.Offset, string(message.Key), string(message.Value))
			resultChan <- true
		}
	}()
	return <-resultChan
}

func CheckKafka() bool {
	if sharedKafkaAdmin == nil {
		log.Println("Kafka admin client is not initialized")
		return false
	}
	brokersList, _, err := sharedKafkaAdmin.DescribeCluster()
	if err != nil {
		log.Printf("Failed to describe Kafka cluster: %v", err)
		return false
	}
	log.Printf("KafkaConnection : Kafka health check successful: found %d broker(s)\n", len(brokersList))
	return true
}

package kafka

import (
	"context"
	"log"
	"strings"
	"sync"

	"github.com/IBM/sarama"
)

var (
	sharedProducer   sarama.SyncProducer
	sharedConsumer   sarama.ConsumerGroup
	sharedKafkaAdmin sarama.ClusterAdmin
	initOnce         sync.Once
)

func InitKafka(brokers string, groupID string) error {
	var err error

	initOnce.Do(func() {
		sharedProducer, err = CreateProducer(brokers)
		if err != nil {
			log.Printf("Failed to initialize Kafka producer: %v", err)
			return
		}

		sharedConsumer, err = CreateConsumer(brokers, groupID)
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

func CreateConsumer(brokers string, groupID string) (sarama.ConsumerGroup, error) {
	brokerList := strings.Split(brokers, ",")

	config := sarama.NewConfig()
	config.Version = sarama.V2_1_0_0
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	consumerGroup, err := sarama.NewConsumerGroup(brokerList, groupID, config)
	if err != nil {
		return nil, err
	}

	return consumerGroup, nil
}

type Consumer struct {
	ready    chan bool
	received chan bool
}

func (c *Consumer) Setup(sarama.ConsumerGroupSession) error {
	close(c.ready)
	return nil
}

func (c *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		log.Printf("KafkaConsum : Message received: topic=%s partition=%d offset=%d key=%s value=%s",
			message.Topic, message.Partition, message.Offset, string(message.Key), string(message.Value))
		session.MarkMessage(message, "")
		c.received <- true
	}
	return nil
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

func CheckConsume(topic string) bool {
	if sharedConsumer == nil {
		log.Println("Kafka consumer is not initialized")
		return false
	}
	consumer := &Consumer{
		ready:    make(chan bool),
		received: make(chan bool),
	}
	ctx, _ := context.WithCancel(context.Background())
	// defer cancel()

	go func() {
		for {
			err := sharedConsumer.Consume(ctx, []string{topic}, consumer)
			if err != nil {
				log.Printf("Error during message consumption: %v", err)
				return
			}
		}
	}()

	select {
	case <-consumer.received:
		log.Println("Consumer received at least one message")
		return true
	case <-ctx.Done():
		log.Println("Context timeout without consuming messages")
		return false
	}
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

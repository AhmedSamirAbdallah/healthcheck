package kafka

import (
	"context"
	"log"
	"strings"

	"github.com/IBM/sarama"
)

func CreateProducer(brokers string) (sarama.SyncProducer, error) {
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

	log.Printf("Message sent to partition %d with offset %d\n", partition, offset)
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
		log.Printf("Message received: topic=%s partition=%d offset=%d key=%s value=%s",
			message.Topic, message.Partition, message.Offset, string(message.Key), string(message.Value))
		session.MarkMessage(message, "")
		c.received <- true
	}
	return nil
}

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

	select {
	case <-consumer.received:
		log.Println("Consumer received at least one message")
		return true
	case <-ctx.Done():
		log.Println("Context done without consuming messages")
		return false
	}
}

func CheckKafka(brokers string) bool {
	brokerList := strings.Split(brokers, ",")
	log.Printf("brokers : %v", brokerList)
	config := sarama.NewConfig()
	config.Version = sarama.V2_0_0_0
	config.ClientID = "health-check-client"

	admin, err := sarama.NewClusterAdmin(brokerList, config)
	if err != nil {
		log.Printf("failed to create Kafka admin client: %v", err)
		return false
	}
	defer admin.Close()

	brokersList, _, err := admin.DescribeCluster()
	if err != nil {
		log.Printf("failed to describe cluster: %v", err)
		return false
	}
	log.Printf("Kafka health check successful: found %d broker(s)\n", len(brokersList))
	return true
}

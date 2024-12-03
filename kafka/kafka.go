package kafka

import (
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

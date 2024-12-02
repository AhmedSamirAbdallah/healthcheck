package kafka

import (
	"log"
	"strings"

	"github.com/IBM/sarama"
)

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

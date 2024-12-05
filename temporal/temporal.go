package temporal

import (
	"log"
	"sync"

	"go.temporal.io/sdk/client"
)

var (
	initOnce       sync.Once
	temporalClient *client.Client
)

func CheckTemporalConnection(url string) bool {
	initOnce.Do(func() {

		options := client.Options{
			HostPort: url,
		}
		temporalClient, err := client.Dial(options)

		if err != nil {
			log.Printf("Error creating Temporal client: %v", err)
			return
		}
		log.Printf("temp : : %v", temporalClient)
		// _, err := temporalClient.
		// if err != nil {
		// 	log.Printf("Error connecting to Temporal namespace: %v", err)
		// 	temporalClient = nil
		// 	return
		// }
		log.Println("Temporal client initialized and connected successfully.")

	})
	return temporalClient != nil
}

package temporal

import (
	"log"
	"net/http"
	"sync"
)

var (
	initOnce       sync.Once
	temporalClient *http.Client
)

func CheckTemporalConnection(url string) bool {
	initOnce.Do(func() {
		client := &http.Client{}

		response, err := client.Get(url)
		if err != nil {
			log.Printf("Error connecting to Temporal service: %v", err)
			return
		}
		if response.StatusCode == http.StatusOK {
			log.Println("Temporal service is healthy.")
			temporalClient = client
			return
		}

	})
	return temporalClient != nil
}

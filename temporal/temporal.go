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

		// withTLS := os.Getenv("TEMPORAL_WITH_TLS") == "true"
		// var tlsConfig *tls.Config
		// if withTLS {
		// 	tlsConfig = &tls.Config{}
		// } else {
		// 	tlsConfig = nil
		// }
		c, err := client.Dial(client.Options{
			HostPort: url,
			// ConnectionOptions: client.ConnectionOptions{
			// 	TLS: tlsConfig,
			// },
		})
		if err != nil {
			log.Fatal(err)
		}
		//
		log.Println("Temporal client initialized and connected successfully.")
		temporalClient = &c

	})
	return temporalClient != nil
}

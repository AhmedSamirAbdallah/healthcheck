package temporal

import (
	"crypto/tls"
	"log"
	"sync"

	"go.temporal.io/sdk/client"
)

var (
	initOnce       sync.Once
	temporalClient *client.Client
)

func CheckTemporalConnection(url string, withTLS bool) bool {
	initOnce.Do(func() {
		var tlsConfig *tls.Config
		if withTLS {
			tlsConfig = &tls.Config{}
		} else {
			tlsConfig = nil
		}
		c, err := client.Dial(client.Options{
			HostPort: url,
			ConnectionOptions: client.ConnectionOptions{
				TLS: tlsConfig,
			},
		})
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Temporal client initialized and connected successfully.")
		temporalClient = &c
	})
	return temporalClient != nil
}

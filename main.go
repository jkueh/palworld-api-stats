package main

import (
	"fmt"
	"os"
	"time"

	palworld_api_client "github.com/jkueh/palworld-api-stats/client"
)

const RestAPIPasswordEnvKey string = "REST_API_PASSWORD"
const RestAPIHostnameEnvKey string = "REST_API_HOSTNAME"

var RestAPIPassword string
var RestAPIHostname string
var appVersion string = "unknown"

func main() {

	if Verbose {
		fmt.Println("Version:", appVersion)
	}

	RestAPIPassword = os.Getenv(RestAPIPasswordEnvKey)
	if RestAPIPassword == "" {
		fmt.Fprintf(os.Stderr, "Error: API Password environment variable (%s) not set\n", RestAPIPasswordEnvKey)
		os.Exit(1)
	}

	RestAPIHostname = os.Getenv(RestAPIHostnameEnvKey)

	if RestAPIHostname == "" {
		RestAPIHostname = "localhost"
		if Verbose {
			fmt.Fprintf(
				os.Stderr,
				"[WARNING]: API Hostname environment variable (%s) not set, defaulting to '%s'\n",
				RestAPIHostnameEnvKey,
				RestAPIHostname,
			)
		}
	}

	client := palworld_api_client.New(&palworld_api_client.RESTAPIClientConfig{
		// As far as I can tell, there's no way to change the username on the REST API side, so leaving it statically
		// defined... For now
		Username: "admin",
		Password: RestAPIPassword,
		Host:     RestAPIHostname,
	})

	// info := client.GetInfo()
	// if Verbose {
	// 	fmt.Println("Connected to server:", info.ServerName)
	// }
	interval := time.Duration(MetricsInterval) * time.Second
	if Verbose {
		fmt.Println("Starting ticker with interval of", interval.String())
	}
	for range time.Tick(interval) {
		go func() {
			PublishMetric(&client)
			if Verbose {
				fmt.Println("Published metric at:", time.Now().String())
			}
		}()
	}

	if Verbose {
		fmt.Println("Done!")
	}

}

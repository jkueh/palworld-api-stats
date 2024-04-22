package main

import (
	"fmt"
	"os"
	"time"

	"github.com/jkueh/palworld-api-stats/cloudwatch_api_client"
	palworld_api_client "github.com/jkueh/palworld-api-stats/palworld_api_client"
)

const RestAPIPasswordEnvKey string = "REST_API_PASSWORD"
const RestAPIHostnameEnvKey string = "REST_API_HOSTNAME"

var RestAPIPassword string
var appVersion string = "unknown"

func main() {

	if Verbose {
		fmt.Println("Starting: Version ", appVersion)
	}

	RestAPIPassword = os.Getenv(RestAPIPasswordEnvKey)
	if RestAPIPassword == "" {
		fmt.Fprintf(os.Stderr, "Error: API Password environment variable (%s) not set\n", RestAPIPasswordEnvKey)
		os.Exit(1)
	}

	palworld_client := palworld_api_client.New(&palworld_api_client.ClientConfig{
		// As far as I can tell, there's no way to change the username on the REST API side, so leaving it statically
		// defined... For now
		Username: "admin",
		Password: RestAPIPassword,
		Host:     RestAPIHostname,
		Port:     RestAPIPort,
		Verbose:  Verbose,
		Debug:    Debug,
	})

	if InfoRequested {
		info := palworld_client.GetInfo()
		fmt.Println("Connected to server:", info.ServerName)
		fmt.Println("Server version:", info.Version)
		os.Exit(0)
	}

	// Get the Cloudwatch client struct
	cloudwatch_client := cloudwatch_api_client.New(&cloudwatch_api_client.ClientConfig{
		MetricsNamespace: MetricsNamespace,
		Verbose:          Verbose,
		Debug:            Debug,
	})

	interval := time.Duration(MetricsInterval) * time.Second
	if Verbose {
		fmt.Println("Starting ticker with interval of", interval.String())
	}
	for range time.Tick(interval) {
		go func() {
			cloudwatch_client.PublishMetrics(palworld_client.GetMetrics())
			if Debug {
				fmt.Println("PublishMetrics() called at:", time.Now().String())
			}
		}()
	}

	if Verbose {
		fmt.Println("Done!")
	}

}

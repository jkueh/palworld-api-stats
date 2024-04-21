package main

import (
	"fmt"
	"os"
	"strings"

	palworld_api_client "github.com/jkueh/palworld-api-stats/client"
)

const VerboseEnvKey string = "VERBOSE"
const RestAPIPasswordEnvKey string = "REST_API_PASSWORD"
const RestAPIHostnameEnvkey string = "REST_API_HOSTNAME"

var RestAPIPassword string
var RestAPIHostname string
var Verbose bool
var appVersion string = "unknown"

func init() {
	// Set the verbose flag
	Verbose = strings.ToLower(os.Getenv(VerboseEnvKey)) == "true"
}

func main() {

	if Verbose {
		fmt.Println("Version:", appVersion)
	}

	RestAPIPassword = os.Getenv(RestAPIPasswordEnvKey)
	if RestAPIPassword == "" {
		fmt.Fprintf(os.Stderr, "Error: API Password environment variable (%s) not set\n", RestAPIPasswordEnvKey)
		os.Exit(1)
	}

	RestAPIHostname = os.Getenv(RestAPIHostnameEnvkey)

	if RestAPIHostname == "" {
		RestAPIHostname = "localhost"
		if Verbose {
			fmt.Fprintf(
				os.Stderr,
				"[WARNING]: API Hostname environment variable (%s) not set, defaulting to '%s'\n",
				RestAPIHostnameEnvkey,
				RestAPIHostname,
			)
		}
	}
	fmt.Println("Password:", RestAPIPassword)

	client := palworld_api_client.New(&palworld_api_client.RESTAPIClientConfig{
		// As far as I can tell, there's no way to change the username on the REST API side, so leaving it statically
		// defined... For now
		Username: "admin",
		Password: RestAPIPassword,
		Host:     "127.0.0.1",
	})

	serverInfo := client.GetInfo()
	if Verbose {
		fmt.Println("Connected to server:", serverInfo.ServerName)
	}
}

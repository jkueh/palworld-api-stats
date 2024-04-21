package main

import (
	"flag"
	"fmt"
	"os"
)

var Verbose bool
var InfoRequested bool
var MetricsInterval int
var MetricsNamespace string
var RestAPIHostname string
var RestAPIPort int

func init() {
	// Set the verbose flag
	var versionRequested bool
	flag.BoolVar(&versionRequested, "version", false, "Output version information and exit.")
	flag.BoolVar(&Verbose, "verbose", false, "Enables verbose logging.")
	flag.BoolVar(&InfoRequested, "info", false, "Outputs information about the server, then exits.")
	flag.IntVar(&MetricsInterval, "interval", 60, "Set how frequently (in seconds) metrics will be published.")
	flag.StringVar(
		&MetricsNamespace,
		"namespace",
		"palworld-api-stats",
		"The namespace that metrics will be published to",
	)
	flag.StringVar(
		&RestAPIHostname,
		"hostname",
		"localhost",
		"The host that is serving the REST API.",
	)
	flag.StringVar(
		&RestAPIHostname,
		"hostname",
		"localhost",
		"The hostname for the server running the REST API",
	)
	flag.IntVar(
		&RestAPIPort,
		"port",
		8212,
		"The port that the REST API service is listening to.",
	)

	flag.Parse()

	// If the version was requested, we can exit early here
	if versionRequested {
		fmt.Println("Version:", appVersion)
		os.Exit(0)
	}
}

package main

import (
	"flag"
	"fmt"
	"os"
)

var Verbose bool
var Debug bool
var InfoRequested bool
var MetricsInterval int
var MetricsNamespace string
var RestAPIHostname string
var RestAPIPort int
var CloudwatchRegion string
var CloudwatchStorageResolution int

func init() {
	// Set the verbose flag
	var versionRequested bool
	flag.BoolVar(&versionRequested, "version", false, "Output version information and exit.")
	flag.BoolVar(&Verbose, "verbose", false, "Enables verbose logging.")
	flag.BoolVar(&Debug, "debug", false, "Enables debug logging - Warning: Noisy!")
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
	flag.IntVar(
		&RestAPIPort,
		"port",
		8212,
		"The port that the REST API service is listening to.",
	)
	flag.StringVar(
		&CloudwatchRegion,
		"cloudwatch-region",
		"",
		"The AWS region to publish Cloudwatch metrics to",
	)
	flag.IntVar(
		&CloudwatchStorageResolution,
		"cloudwatch-storage-resolution",
		60,
		"The Cloudwatch storage resolution - Set to '1' for high-resolution metrics "+
			"(Note that this has a different cost implication to standard metrics!)",
	)

	flag.Parse()

	// If the version was requested, we can exit early here
	if versionRequested {
		fmt.Println("Version:", appVersion)
		os.Exit(0)
	}

	// Double check that we actually have a manually-specified value for the cloudwatch region
	if CloudwatchRegion == "" {
		fmt.Fprintln(os.Stderr, "Error: Cloudwatch Region (-cloudwatch-region) not specified.")
		os.Exit(1)
	}
}

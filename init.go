package main

import (
	"flag"
	"fmt"
	"os"
)

var Verbose bool
var MetricsInterval int

func init() {
	// Set the verbose flag
	var versionRequested bool
	flag.BoolVar(&versionRequested, "version", false, "Output version information and exit.")
	flag.BoolVar(&Verbose, "verbose", false, "Enables verbose logging.")
	flag.IntVar(&MetricsInterval, "interval", 60, "Set how frequently (in seconds) metrics will be published")

	flag.Parse()

	// If the version was requested, we can exit early here
	if versionRequested {
		fmt.Println("Version:", appVersion)
		os.Exit(0)
	}
}

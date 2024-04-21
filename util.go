package main

import (
	"fmt"

	palworld_api_client "github.com/jkueh/palworld-api-stats/client"
)

func PublishMetric(c *palworld_api_client.RESTAPIClient) {
	metrics := c.GetMetrics()
	fmt.Println("Current Server FPS:", metrics.ServerFPS)
}

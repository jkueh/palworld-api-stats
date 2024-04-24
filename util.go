package main

import (
	"fmt"

	palworld_api_client "github.com/jkueh/palworld-api-stats/palworld_api_client"
)

func PublishMetric(c *palworld_api_client.Client) {
	metrics := c.GetMetrics()
	fmt.Println("Current Server FPS:", metrics.ServerFPS)
}

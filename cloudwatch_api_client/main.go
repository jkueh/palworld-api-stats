package cloudwatch_api_client

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/jkueh/palworld-api-stats/types/responses"
)

type ClientConfig struct {
	MetricsNamespace string
}

type Client struct {
	config *ClientConfig
	client *cloudwatch.CloudWatch
}

// Returns the Client struct, populated with an active cloudwatch client
func New(config *ClientConfig) *Client {
	session := session.Must(session.NewSession())
	return &Client{
		config: config,
		client: cloudwatch.New(session),
	}
}

func (c *Client) PublishMetrics(m *responses.MetricsResponse) {
	fmt.Println("Server FPS:", m.ServerFPS)
	fmt.Println("Not yet implemented")
}

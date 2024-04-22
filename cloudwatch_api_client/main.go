package cloudwatch_api_client

import (
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/jkueh/palworld-api-stats/types/responses"
)

type ClientConfig struct {
	MetricsNamespace string
	Region           string
	Verbose          bool
	Debug            bool
}

type Client struct {
	config *ClientConfig
	client *cloudwatch.CloudWatch
}

// Returns the Client struct, populated with an active cloudwatch client
func New(config *ClientConfig) *Client {
	session := session.Must(session.NewSession(&aws.Config{
		Region: &config.Region,
	}))
	if config.Verbose {
		fmt.Println("Configuring Cloudwatch for region:", config.Region)
	}
	return &Client{
		config: config,
		client: cloudwatch.New(session),
	}
}

func (c *Client) PublishMetrics(m *responses.MetricsResponse) {
	var metricData []*cloudwatch.MetricDatum

	metricTime := time.Now()

	metricData = append(metricData, &cloudwatch.MetricDatum{
		MetricName: aws.String("ServerFPS"),
		Timestamp:  &metricTime,
		Unit:       aws.String("FPS"),
		Values:     []*float64{aws.Float64(float64(m.ServerFPS))},
	})

	_, err := c.client.PutMetricData(&cloudwatch.PutMetricDataInput{
		Namespace:  &c.config.MetricsNamespace,
		MetricData: metricData,
	})
	if err != nil {
		fmt.Println("An error occurred while publishing metrics:", err)
		os.Exit(127)
	}
}

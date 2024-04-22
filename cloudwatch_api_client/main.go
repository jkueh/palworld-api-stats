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
	MetricsNamespace            string
	Region                      string
	Verbose                     bool
	Debug                       bool
	CloudwatchStorageResolution int64
}

type Client struct {
	config *ClientConfig
	client *cloudwatch.CloudWatch
}

// Returns the Client struct, populated with an active cloudwatch client
func New(config *ClientConfig) *Client {
	session := session.Must(session.NewSession())
	if config.Verbose {
		fmt.Println("Configuring Cloudwatch for region:", config.Region)
	}
	return &Client{
		config: config,
		client: cloudwatch.New(session, &aws.Config{
			Region: &config.Region,
		}),
	}
}

func (c *Client) PublishMetrics(m *responses.MetricsResponse) {
	var metricData []*cloudwatch.MetricDatum

	metricTime := time.Now()

	metricData = append(metricData, &cloudwatch.MetricDatum{
		MetricName:        aws.String("ServerFPS"),
		Timestamp:         &metricTime,
		Values:            []*float64{aws.Float64(float64(m.ServerFPS))},
		StorageResolution: &c.config.CloudwatchStorageResolution,
	})
	metricData = append(metricData, &cloudwatch.MetricDatum{
		MetricName:        aws.String("CurrentPlayerNum"),
		Timestamp:         &metricTime,
		Values:            []*float64{aws.Float64(float64(m.CurrentPlayerNum))},
		StorageResolution: &c.config.CloudwatchStorageResolution,
	})
	metricData = append(metricData, &cloudwatch.MetricDatum{
		MetricName:        aws.String("ServerFrameTime"),
		Timestamp:         &metricTime,
		Values:            []*float64{aws.Float64(m.ServerFrameTime)},
		StorageResolution: &c.config.CloudwatchStorageResolution,
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

package palworld_api_client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/jkueh/palworld-api-stats/types/responses"
)

type RESTAPIClientConfig struct {
	Host     string
	Username string
	Password string
	Port     int
	Verbose  bool
}

type RESTAPIClient struct {
	config *RESTAPIClientConfig
	client http.Client
}

func New(c *RESTAPIClientConfig) RESTAPIClient {
	return RESTAPIClient{
		config: c,
		client: http.Client{},
	}
}

func (c *RESTAPIClient) DoNothing() {
	fmt.Println("Nothing: Done")
}

// A wrapper around http.Client.Do(), so we can inject authorisation details
func (c *RESTAPIClient) Do(method string, endpoint string) (*http.Response, error) {
	host := fmt.Sprintf("%s:%d", c.config.Host, c.config.Port)
	if c.config.Verbose {
		fmt.Println("Host:", host)
	}
	req := http.Request{
		Method: method,
		URL: &url.URL{
			// We're assuming no one's doing mutual TLS here, or something strange like exposing their API endpoints to the
			// public internet
			Scheme: "http",
			User:   url.UserPassword(c.config.Username, c.config.Password),
			Host:   host,
			Path:   fmt.Sprintf("/v1/api/%s", endpoint),
		},
		Header: map[string][]string{
			"content-type": {"application/json"},
		},
	}
	// TODO: Connection timeout using http.Transport
	// (e.g. in the case of high-frequency metrics publishing, it makes more sense to omit than publish old data late)
	return c.client.Do(&req)
}

// Convenience method for /info, great as a pre-flight check
func (c *RESTAPIClient) GetInfo() *responses.ServerInfoResponse {
	var respPayload responses.ServerInfoResponse
	// TODO: 	Find a good way to abstract this logic so it's not copy-paste every time, while retaining some degree of
	// 				decent error handling
	resp, restErr := c.Do("GET", "info")
	if restErr != nil {
		fmt.Fprintln(os.Stderr, "RESTAPIClient.GetInfo() - Error returned from API")
		fmt.Fprintln(os.Stderr, restErr)
		// Exit here if this fails - If this endpoint doesn't work, it's likely nothing else will
		os.Exit(127)
	}
	body, bodyReadErr := io.ReadAll(resp.Body)
	if bodyReadErr != nil {
		fmt.Fprintln(os.Stderr, "RESTAPIClient.GetInfo() - Unable to read response body")
		fmt.Fprintln(os.Stderr, bodyReadErr)
		os.Exit(128)
	}

	// Convert the (presumed) JSON body into the response payload struct
	unmarshalErr := json.Unmarshal(body, &respPayload)
	if unmarshalErr != nil {
		fmt.Fprintln(os.Stderr, "RESTAPIClient.GetInfo() - Unable to parse JSON response")
		fmt.Fprintln(os.Stderr, bodyReadErr)
		os.Exit(129)
	}

	return &respPayload
}

// Function to return the metrics response
func (c *RESTAPIClient) GetMetrics() *responses.MetricsResponse {
	var respPayload responses.MetricsResponse
	resp, restErr := c.Do("GET", "metrics")
	if restErr != nil {
		fmt.Fprintln(os.Stderr, "RESTAPIClient.GetMetrics() - Error returned from API")
		fmt.Fprintln(os.Stderr, restErr)
		// Exit here if this fails - If this endpoint doesn't work, it's likely nothing else will
		os.Exit(127)
	}
	body, bodyReadErr := io.ReadAll(resp.Body)
	if bodyReadErr != nil {
		fmt.Fprintln(os.Stderr, "RESTAPIClient.GetMetrics() - Unable to read response body")
		fmt.Fprintln(os.Stderr, bodyReadErr)
		os.Exit(128)
	}

	// Convert the (presumed) JSON body into the response payload struct
	unmarshalErr := json.Unmarshal(body, &respPayload)
	if unmarshalErr != nil {
		fmt.Fprintln(os.Stderr, "RESTAPIClient.GetMetrics() - Unable to parse JSON response")
		fmt.Fprintln(os.Stderr, bodyReadErr)
		os.Exit(129)
	}

	return &respPayload
}

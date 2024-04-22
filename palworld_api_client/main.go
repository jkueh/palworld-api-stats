package api_client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/jkueh/palworld-api-stats/types/responses"
)

type ClientConfig struct {
	Host     string
	Username string
	Password string
	Port     int
	Verbose  bool
	Debug    bool
}

type Client struct {
	config *ClientConfig
	client http.Client
}

func New(c *ClientConfig) Client {
	return Client{
		config: c,
		client: http.Client{},
	}
}

// A wrapper around http.Client.Do(), so we can inject authorisation details
func (c *Client) Do(method string, endpoint string) (*http.Response, error) {
	host := fmt.Sprintf("%s:%d", c.config.Host, c.config.Port)
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
func (c *Client) GetInfo() *responses.ServerInfoResponse {
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

	if resp.StatusCode != 200 {
		fmt.Println("Error code returned from server:", resp.Status)
		if len(body) > 0 {
			fmt.Println(body)
		}
		os.Exit(129)
	}

	// Convert the (presumed) JSON body into the response payload struct
	unmarshalErr := json.Unmarshal(body, &respPayload)
	if unmarshalErr != nil {
		fmt.Fprintln(os.Stderr, "RESTAPIClient.GetInfo() - Unable to parse JSON response")
		if c.config.Verbose {
			fmt.Println(string(body))
		}
		fmt.Fprintln(os.Stderr, bodyReadErr)
		os.Exit(130)
	}

	return &respPayload
}

// Function to return the metrics response
func (c *Client) GetMetrics() *responses.MetricsResponse {
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

	if resp.StatusCode != 200 {
		fmt.Println("Error code returned from server:", resp.Status)
		if len(body) > 0 {
			fmt.Println(body)
		}
		os.Exit(129)
	}

	if c.config.Debug {
		fmt.Println("Received", resp.StatusCode, "response from metrics endpoint")
	}

	// Convert the (presumed) JSON body into the response payload struct
	unmarshalErr := json.Unmarshal(body, &respPayload)
	if unmarshalErr != nil {
		fmt.Fprintln(os.Stderr, "RESTAPIClient.GetMetrics() - Unable to parse JSON response")
		if c.config.Verbose {
			fmt.Println(string(body))
		}
		fmt.Fprintln(os.Stderr, bodyReadErr)
		os.Exit(130)
	}

	return &respPayload
}

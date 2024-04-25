package api_client

import (
	"encoding/json"
	"errors"
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
func (c *Client) Do(method string, endpoint string) ([]byte, error) {
	host := fmt.Sprintf("%s:%d", c.config.Host, c.config.Port)

	// TODO: Connection timeout using http.Transport
	// (e.g. in the case of high-frequency metrics publishing, it makes more sense to omit than publish old data late)
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

	resp, respErr := c.client.Do(&req)
	if respErr != nil {
		fmt.Fprintf(os.Stderr, "palworld_api_client.Client.Do() - Error returned from '%s' endpoint\n", endpoint)
		fmt.Fprintln(os.Stderr, respErr)
		return []byte{}, respErr
	}

	// Check to see if we can read the response body
	body, bodyReadErr := io.ReadAll(resp.Body)
	if bodyReadErr != nil {
		fmt.Fprintf(
			os.Stderr,
			"palworld_api_client.Client.Do() - Unable to read response body from '%s' endpoint",
			endpoint,
		)
		fmt.Fprintln(os.Stderr, bodyReadErr)
		return []byte{}, &BodyReadError{error: bodyReadErr}
	}

	// Check for a non-200 response code
	if resp.StatusCode != 200 {
		fmt.Println("Error code returned from server:", resp.Status)
		if len(body) > 0 {
			fmt.Println(body)
		}
		return []byte{}, &Non200Error{statusCode: resp.StatusCode}
	}

	return body, nil
}

// Convenience method for /info, great as a pre-flight check
func (c *Client) GetInfo() *responses.ServerInfoResponse {
	var respPayload responses.ServerInfoResponse
	body, doErr := c.Do("GET", "info")
	if doErr != nil {
		if errors.As(doErr, &Non200Error{}) {
			os.Exit(128)
		} else if errors.As(doErr, &BodyReadError{}) {
			os.Exit(129)
		} else {
			fmt.Fprintln(os.Stderr, "GetInfo(): Unhandled error", doErr)
		}
	}

	// Convert the (presumed) JSON body into the response payload struct
	unmarshalErr := json.Unmarshal(body, &respPayload)
	if unmarshalErr != nil {
		fmt.Fprintln(os.Stderr, "RESTAPIClient.GetInfo() - Unable to parse JSON response")
		if c.config.Verbose {
			fmt.Println(string(body))
		}
		fmt.Fprintln(os.Stderr, unmarshalErr)
		os.Exit(130)
	}

	return &respPayload
}

// Function to return the metrics response
func (c *Client) GetMetrics() *responses.MetricsResponse {
	var respPayload responses.MetricsResponse
	body, doErr := c.Do("GET", "metrics")
	if doErr != nil {
		if errors.As(doErr, &Non200Error{}) {
			os.Exit(128)
		} else if errors.As(doErr, &BodyReadError{}) {
			os.Exit(129)
		} else {
			fmt.Fprintln(os.Stderr, "GetMetrics(): Unhandled error", doErr)
		}
	}

	// Convert the (presumed) JSON body into the response payload struct
	unmarshalErr := json.Unmarshal(body, &respPayload)
	if unmarshalErr != nil {
		fmt.Fprintln(os.Stderr, "GetMetrics() - Unable to parse JSON response")
		if c.config.Verbose {
			fmt.Println(string(body))
		}
		fmt.Fprintln(os.Stderr, unmarshalErr)
		os.Exit(130)
	}

	return &respPayload
}

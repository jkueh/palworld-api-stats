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
	req := http.Request{
		Method: method,
		URL: &url.URL{
			User: url.UserPassword(c.config.Username, c.config.Password),
			Host: c.config.Host,
			Path: fmt.Sprintf("/v1/api/%s", endpoint),
		},
	}
	req.Header.Add("Content-type", "application/json")
	return c.client.Do(&req)
}

// Convenience method for /info, great as a pre-flight check
func (c *RESTAPIClient) GetInfo() *responses.ServerInfoResponse {
	var respPayload responses.ServerInfoResponse
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

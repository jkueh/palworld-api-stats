package palworld_api_client

import (
	"fmt"
	"net/http"
	"net/url"
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

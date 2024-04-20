package palworld_api_client

import "fmt"

type RESTAPIClientConfig struct {
	Host     string
	Username string
	Password string
}

type RESTAPIClient struct {
	config *RESTAPIClientConfig
}

func New(c *RESTAPIClientConfig) RESTAPIClient {
	return RESTAPIClient{config: c}
}

func (c *RESTAPIClient) DoNothing() {
	fmt.Println("Nothing: Done")
}

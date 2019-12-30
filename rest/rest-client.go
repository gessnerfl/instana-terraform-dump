package rest

import (
	"errors"
	"fmt"
	"log"

	"gopkg.in/resty.v1"
)

//ErrEntityNotFound error message which is returned when the entity cannot be found at the server
var ErrEntityNotFound = errors.New("Failed to get resource from Instana API. 404 - Resource not found")

//NewClient creates a new instance of the Instana REST API client
func NewClient(apiToken string, host string) Client {
	restyClient := resty.New()

	return &clientImpl{
		apiToken:    apiToken,
		host:        host,
		restyClient: restyClient,
	}
}

//Client interface for the simplified REST client
type Client interface {
	Get(resourcePath string) ([]byte, error)
}

//RestClientImpl is a helper class to interact with Instana REST API
type clientImpl struct {
	apiToken    string
	host        string
	restyClient *resty.Client
}

var emptyResponse = make([]byte, 0)

//GetOne request the resource with the given ID
func (client *clientImpl) Get(resourcePath string) ([]byte, error) {
	url := client.buildURL(resourcePath)
	log.Printf("Call GET %s\n", url)
	resp, err := client.createRequest().Get(url)
	if err != nil {
		return emptyResponse, fmt.Errorf("failed to send HTTP GET request to Instana API; status code = %d; status message = %s, %s", resp.StatusCode(), resp.Status(), err)
	}
	statusCode := resp.StatusCode()
	if statusCode == 404 {
		return emptyResponse, ErrEntityNotFound
	}
	if statusCode < 200 || statusCode >= 300 {
		return emptyResponse, fmt.Errorf("failed to send HTTP GET request to Instana API; status code = %d; status message = %s", statusCode, resp.Status())
	}
	return resp.Body(), nil
}

func (client *clientImpl) createRequest() *resty.Request {
	return client.restyClient.R().SetHeader("Accept", "application/json").SetHeader("Authorization", fmt.Sprintf("apiToken %s", client.apiToken))
}

func (client *clientImpl) buildURL(resourcePath string) string {
	return fmt.Sprintf("https://%s%s", client.host, resourcePath)
}

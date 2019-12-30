package customevents

import (
	"encoding/json"
	"fmt"

	"github.com/gessnerfl/instana-terraform-dump/rest"
)

//CustomEventAPI Interface for the API access to Instana REST API for Custom Events
type CustomEventAPI interface {
	GetAll() ([]*CustomEventSpecification, error)
}

//NewCustomEventAPI creates a new instance of the CustomEventAPI
func NewCustomEventAPI(restClient rest.Client) CustomEventAPI {
	return &customEventAPIImpl{
		restClient: restClient,
	}
}

type customEventAPIImpl struct {
	restClient rest.Client
}

func (api *customEventAPIImpl) GetAll() ([]*CustomEventSpecification, error) {
	data, err := api.restClient.Get("/api/events/settings/event-specifications/custom")
	if err != nil {
		return make([]*CustomEventSpecification, 0), err
	}

	specs := make([]*CustomEventSpecification, 0)
	if err := json.Unmarshal(data, &specs); err != nil {
		return specs, fmt.Errorf("failed to parse json; %s", err)
	}
	return specs, nil
}

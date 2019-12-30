package modules

import (
	"bytes"

	"github.com/gessnerfl/instana-terraform-dump/modules/customevents"
	"github.com/gessnerfl/instana-terraform-dump/rest"
)

type customEventsModule struct {
	restClient rest.Client
}

func (m *customEventsModule) AppendInstanaResourcesTo(buffer *bytes.Buffer) error {
	customEventsAPI := customevents.NewCustomEventAPI(m.restClient)
	customEventSpecs, err := customEventsAPI.GetAll()
	if err != nil {
		return err
	}

	customEventsWriter := customevents.NewCustomEventWriter(buffer)
	err = customEventsWriter.Write(customEventSpecs)
	if err != nil {
		return err
	}
	return nil
}

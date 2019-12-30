package modules

import "bytes"
import "github.com/gessnerfl/instana-terraform-dump/rest"

//Module interface definition of an instana resource module
type Module interface {
	AppendInstanaResourcesTo(buffer *bytes.Buffer) error
}

//NewFactory creates a new instance of a ModuleFactory
func NewFactory(restClient rest.Client) ModuleFactory {
	return &modules{
		restClient: restClient,
	}
}

//ModuleFactory interface definition of a module factory which is used to create new instances of the different modules
type ModuleFactory interface {
	CustomEvents() Module
}

type modules struct {
	restClient rest.Client
}

func (mods *modules) CustomEvents() Module {
	return &customEventsModule{
		restClient: mods.restClient,
	}
}

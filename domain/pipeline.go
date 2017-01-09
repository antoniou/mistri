package domain

import (
	"fmt"
	"log"
)

type Pipeline interface {
	Create([]string) error
	Delete(interface{}) error
	Read(interface{}) interface{}
}

// BasePipeline
type BasePipeline struct {
	Name  string
	Steps []Actor
}

type PipelineFactory func(conf map[string]string) (Pipeline, error)

var pipelineFactories = make(map[string]PipelineFactory)

func Register(name string, factory PipelineFactory) {
	if factory == nil {
		log.Fatalf("Pipeline factory %s does not exist.", name)
	}
	_, registered := pipelineFactories[name]
	if registered {
		log.Fatalf("Pipeline factory %s already registered. Ignoring.", name)
	}
	pipelineFactories[name] = factory
}

func init() {
	Register("AWS_CP", NewAWSCodePipeline)
}

func NewPipeline(conf map[string]string) (Pipeline, error) {
	provider := conf["provider"]
	pipelineFactory, ok := pipelineFactories[provider]

	if !ok {
		return nil, fmt.Errorf("Invalid pipeline type %s.", provider)
	}

	return pipelineFactory(conf)
}

package domain

import (
	"fmt"
	"log"
)

type Pipeline interface {
	Create([]string) error
	Delete(interface{}) error
	Read(interface{}) interface{}
	createSteps()
}

// BasePipeline
type BasePipeline struct {
	Name   string
	Source Source
	Steps  []Actor
}

func (b *BasePipeline) Create(args []string) error {
	log.Printf("Creating Pipeline with name %s", b.Name)
	for _, step := range b.Steps {
		err := step.Run(nil)
		if err != nil {
			return fmt.Errorf("Error during creation of pipeline %s: %s", b.Name, err)
		}
	}

	return nil
}

func (b *BasePipeline) createSteps() []Actor {
	log.Fatal("Not implemented")
	return nil
}

func (b *BasePipeline) Delete(interface{}) error {
	log.Fatal("Not implemented")
	return nil
}

func (b *BasePipeline) Read(interface{}) interface{} {
	log.Fatal("Not implemented")
	return nil
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

package main

import (
	"log"
	"os"

	"github.com/antoniou/zero2Pipe/domain"
)

var pipelineFactories = make(map[string]domain.PipelineFactory)

func main() {
	cli, _ := New()
	cli.Run(os.Args)
}

func Register(name string, factory domain.PipelineFactory) {
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
	Register("AWS_CP", domain.NewAWSCodePipeline)
}

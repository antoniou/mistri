package domain

import (
	"log"
)

// AWSCodePipeline implements the Pipeline interface
type AWSCodePipeline struct {
	Steps []Actor
}

func (p *AWSCodePipeline) Create(args []string) error {
	log.Printf("Creating pipeline %s", args[0])

	p.Steps = []Actor{
		&CloudFormationActor{
			Template:  "templates/lambda-store.json",
			StackName: "s3-lambda-bucket",
		},
		&LambdaActor{
			S3Bucket:       "lambda-store-eu-west-2-329485089133",
			FunctionSource: "templates/lambda",
		},
		&CloudFormationActor{
			Template:  "templates/pipeline.json",
			StackName: "code-pipeline",
		},
	}

	for _, step := range p.Steps {
		step.Run(nil)
	}

	return nil

}

func (p *AWSCodePipeline) Delete(interface{}) error {
	return nil

}

func (p *AWSCodePipeline) Read(interface{}) interface{} {
	return nil

}

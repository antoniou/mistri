package domain

import (
	"log"
)

func NewAWSCodePipeline(conf map[string]string) (Pipeline, error) {
	return &AWSCodePipeline{
		BasePipeline{
			Name: conf["name"],
		},
	}, nil
}

// AWSCodePipeline implements the Pipeline interface
type AWSCodePipeline struct {
	BasePipeline
}

func (p *AWSCodePipeline) Create(args []string) error {
	log.Printf("Creating AWS CodePipeline with name %s", p.Name)

	p.Steps = []Actor{
		&CloudFormationActor{
			Template:  "templates/lambda-store.json",
			StackName: "s3-lambda-bucket",
		},
		&LambdaActor{
			S3Bucket:       "lambda-store-eu-west-1-329485089133",
			FunctionSource: "templates/lambda",
		},
		&CloudFormationActor{
			Template:  "templates/pipeline.json",
			StackName: p.Name,
			Parameters: map[string]string{
				"ApplicationRepository": "nevergreen-standalone",
				"PipelineName":          p.Name,
			},
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

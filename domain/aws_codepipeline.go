package domain

import (
	"log"
)

func NewAWSCodePipeline(conf map[string]string) (Pipeline, error) {
	source, err := NewPathSource(".")
	if err != nil {
		log.Fatal(err)
	}
	p := &AWSCodePipeline{
		BasePipeline{
			Name:   conf["name"],
			Source: source,
		},
	}
	p.createSteps()
	return p, nil
}

// AWSCodePipeline implements the Pipeline interface
type AWSCodePipeline struct {
	BasePipeline
}

func (p *AWSCodePipeline) createSteps() {
	p.Steps = []Actor{
		&CloudFormationActor{
			Template:  "/Users/nassos/workspace/go/src/github.com/antoniou/zero2Pipe/templates/lambda-store.json",
			StackName: "s3-lambda-bucket",
		},
		&LambdaInstallerActor{
			S3Bucket:       "lambda-store-eu-west-1-329485089133",
			S3KeyPrefix:    p.Name,
			FunctionSource: "/Users/nassos/workspace/go/src/github.com/antoniou/zero2Pipe/templates/lambda",
		},
		&CloudFormationActor{
			Template:  "/Users/nassos/workspace/go/src/github.com/antoniou/zero2Pipe/templates/pipeline.json",
			StackName: p.Name,
			Parameters: map[string]string{
				"ApplicationRepositoryName":       p.Source.Name(),
				"ApplicationRepositoryOwner":      p.Source.Owner(),
				"ApplicationRepositoryOAuthToken": p.Source.Auth(),
				"PipelineName":                    p.Name,
			},
		},
	}
}

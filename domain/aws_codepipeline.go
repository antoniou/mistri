package domain

import (
	"fmt"
	"log"
)

func NewAWSCodePipeline(conf map[string]string) (Pipeline, error) {
	source, err := NewPathSource(".")
	if err != nil {
		log.Fatal(err)
	}
	p := &AWSCodePipeline{
		BasePipeline: BasePipeline{
			Name:   conf["name"],
			Source: source,
		},
		Context: NewAWSCodePipelineContext(),
	}
	p.createSteps()
	return p, nil
}

// AWSCodePipeline implements the Pipeline interface
type AWSCodePipeline struct {
	BasePipeline
	Context *AWSCodePipelineContext
}

func (p *AWSCodePipeline) createSteps() {
	lambdaS3Bucket := fmt.Sprintf("lambda-store-%s-%s", p.Context.Props["region"], p.Context.Props["account"])

	p.Steps = []Actor{
		&CloudFormationActor{
			AWSActor: &AWSActor{
				Context: p.Context,
			},
			Template:  "/Users/nassos/workspace/go/src/github.com/antoniou/zero2Pipe/templates/lambda-store.json",
			StackName: fmt.Sprintf("%s-lambda-store", p.Name),
			Parameters: map[string]string{
				"LambdaBucketName": lambdaS3Bucket,
			},
		},
		&LambdaGeneratorActor{
			Generator: NewGenerator("AWSBuildspecGenerator"),
			params: map[string]string{
				"FunctionSource": "/Users/nassos/workspace/go/src/github.com/antoniou/zero2Pipe/templates/lambda/genBuildspec",
				"Template":       "/Users/nassos/workspace/go/src/github.com/antoniou/zero2Pipe/templates/buildspec.yml.tmpl",
				"pipelineName":   p.Name,
				"AWS_ACCOUNT":    p.Context.Props["account"],
				"AWS_REGION":     p.Context.Props["region"],
			},
		},
		&LambdaInstallerActor{
			S3Bucket:       lambdaS3Bucket,
			S3KeyPrefix:    p.Name,
			FunctionSource: "/Users/nassos/workspace/go/src/github.com/antoniou/zero2Pipe/templates/lambda",
		},
		&CloudFormationActor{
			AWSActor: &AWSActor{
				Context: p.Context,
			},
			Template:  "/Users/nassos/workspace/go/src/github.com/antoniou/zero2Pipe/templates/pipeline.json",
			StackName: p.Name,
			Parameters: map[string]string{
				"ApplicationRepositoryName":       p.Source.Name(),
				"ApplicationRepositoryOwner":      p.Source.Owner(),
				"ApplicationRepositoryOAuthToken": p.Source.Auth(),
				"PipelineName":                    p.Name,
				"LambdaBucketName":                lambdaS3Bucket,
			},
		},
	}
}

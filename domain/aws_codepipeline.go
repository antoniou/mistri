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
			Template:  "templates/lambda-store.json",
			StackName: fmt.Sprintf("%s-lambda-store", p.Name),
			Parameters: map[string]string{
				"LambdaBucketName": lambdaS3Bucket,
			},
		},
		&CustomActor{
			Generator: NewGenerator("AWSLambdaCodeExporter"),
			params: map[string]string{
				"dir": "templates",
			},
		},
		&GithubAuthorisationActor{
			AWSActor: &AWSActor{
				Context: p.Context,
			},
		},
		&CustomActor{
			Generator: NewGenerator("AWSTemplateGenerator"),
			params: map[string]string{
				"FunctionSource": "./.mistri/templates/lambda/genBuildspec",
				"Template":       "./.mistri/templates/buildspec.yml.tmpl",
				"pipelineName":   p.Name,
				"AWS_ACCOUNT":    p.Context.Props["account"],
				"AWS_REGION":     p.Context.Props["region"],
			},
		},
		&CustomActor{
			Generator: NewGenerator("AWSTemplateGenerator"),
			params: map[string]string{
				"FunctionSource": "./.mistri/templates/lambda/genBuildParams",
				"Template":       "./.mistri/templates/Dockerrun.aws.json.j2.tmpl",
				"pipelineName":   p.Name,
				"AWS_ACCOUNT":    p.Context.Props["account"],
				"AWS_REGION":     p.Context.Props["region"],
				"host_port":      "80",
				"container_port": "3000",
				"version":        "{{ version }}",
			},
		},
		&LambdaInstallerActor{
			S3Bucket:       lambdaS3Bucket,
			S3KeyPrefix:    p.Name,
			FunctionSource: "./.mistri/templates/lambda",
		},
		&CloudFormationActor{
			AWSActor: &AWSActor{
				Context: p.Context,
			},
			Template:  "templates/pipeline.json",
			StackName: p.Name,
			Parameters: map[string]string{
				"ApplicationRepositoryName":       p.Source.Name(),
				"ApplicationRepositoryOwner":      p.Source.Owner(),
				"PipelineName":                    p.Name,
				"LambdaBucketName":                lambdaS3Bucket,
				"ApplicationRepositoryOAuthToken": "",
			},
		},
	}
}

package domain

import (
	"io/ioutil"
	"log"
	"strings"

	"github.com/antoniou/zero2Pipe/lambda"
)

// LambdaInstallerActor implements Actor
// It installs all Lambda functions found under path
// FunctionSource to AWS Lambda
type LambdaInstallerActor struct {
	S3Bucket       string
	S3KeyPrefix    string
	FunctionSource string
}

// Run is the entrypoint to the Actor workload.
// LambdaInstallerActor Run finds all functions under
// FunctionSource path and uploads them to S3Bucket by
// prefixing them with S3KeyPrefix
func (l *LambdaInstallerActor) Run(interface{}) error {
	log.Println("Installing functions!")
	l.installFunctions()
	return nil
}

func (l *LambdaInstallerActor) installFunction(name string) error {
	log.Printf("[DEBUG] Installing function %s", name)
	f := lambda.NewFunction(map[string]string{
		"name":        name,
		"path":        strings.Join([]string{l.FunctionSource, name}, "/"),
		"s3bucket":    l.S3Bucket,
		"s3KeyPrefix": l.S3KeyPrefix,
	})
	return f.Setup()
}

func (l *LambdaInstallerActor) installFunctions() error {
	files, err := ioutil.ReadDir(l.FunctionSource)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if !file.IsDir() {
			continue
		}

		err := l.installFunction(file.Name())
		if err != nil {
			return err
		}
	}
	return nil
}

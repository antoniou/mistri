package domain

import (
	"io/ioutil"
	"log"
	"strings"

	"github.com/antoniou/zero2Pipe/lambda"
)

type Actor interface {
	Run(interface{}) error
}

// LambdaActor implements Actor
type LambdaActor struct {
	S3Bucket       string
	FunctionSource string
}

func (l *LambdaActor) Run(interface{}) error {
	log.Println("Installing functions!")
	l.installFunctions(l.FunctionSource, l.S3Bucket)
	return nil
}

func (l *LambdaActor) installFunctions(path string, s3bucket string) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if !file.IsDir() {
			continue
		}

		log.Printf("[DEBUG] Installing function %s", file.Name())
		f := lambda.NewFunction(map[string]string{
			"name":        file.Name(),
			"path":        strings.Join([]string{path, file.Name()}, "/"),
			"s3bucket":    s3bucket,
			"s3KeyPrefix": "SimplePipeline",
		})
		f.Setup()
	}
}

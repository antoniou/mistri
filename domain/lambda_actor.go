package domain

import (
	"io/ioutil"
	"log"
	"strings"

	"github.com/antoniou/zero2Pipe/lambda"
)

// LambdaActor implements Actor
type LambdaActor struct {
	S3Bucket       string
	S3KeyPrefix    string
	FunctionSource string
}

func (l *LambdaActor) Run(interface{}) error {
	log.Println("Installing functions!")
	l.installFunctions()
	return nil
}

func (l *LambdaActor) installFunction(name string) error {
	log.Printf("[DEBUG] Installing function %s", name)
	f := lambda.NewFunction(map[string]string{
		"name":        name,
		"path":        strings.Join([]string{l.FunctionSource, name}, "/"),
		"s3bucket":    l.S3Bucket,
		"s3KeyPrefix": l.S3KeyPrefix,
	})
	return f.Setup()
}

func (l *LambdaActor) installFunctions() error {
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

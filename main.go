package main

import (
	"io/ioutil"
	"log"
	"strings"

	"github.com/antoniou/zero2Pipe/lambda"
)

func installFunctions(path string, s3bucket string) {
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

func mains() {
	installFunctions("functions", "lambda-store-eu-west-1-329485089133")
}

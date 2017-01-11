package domain

import (
	"os"
	"path/filepath"
	"text/template"
)

var generators = make(map[string]Generator)

type Generator func(conf map[string]string) error

func NewGenerator(Type string) Generator {
	return generators[Type]
}

func AWSBuildspecGenerator(conf map[string]string) error {
	templates := template.Must(template.ParseGlob(conf["Template"]))
	fo, err := os.Create(filepath.Join(conf["FunctionSource"], "buildspec.yml"))
	if err != nil {
		return err
	}

	err = templates.Execute(fo, conf)
	if err != nil {
		return err
	}

	return nil
}

type LambdaGeneratorActor struct {
	FunctionSource string
	Generator      Generator
	params         map[string]string
}

func (l *LambdaGeneratorActor) Run(interface{}) error {
	return l.Generator(l.params)
}

func registerGenerator(name string, gen Generator) {
	generators[name] = gen
}

func init() {
	registerGenerator("AWSBuildspecGenerator", AWSBuildspecGenerator)
}

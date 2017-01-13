package domain

import (
	"log"
	"os"
	"path/filepath"
	"text/template"
)

var generators = make(map[string]Generator)

type Generator func(conf map[string]string) error

func NewGenerator(Type string) Generator {
	return generators[Type]
}

func AWSLambdaCodeExporter(conf map[string]string) error {
	return export(conf["dir"])
}

func export(file string) error {
	var lambdaDir []string
	lambdaDir, _ = AssetDir(file)
	os.MkdirAll(filepath.Join(".zero2Pipe", file), 0777)

	for _, f := range lambdaDir {
		absf := filepath.Join(".", file, f)
		_, err := AssetInfo(absf)
		if err != nil {
			export(absf)
			continue
		}

		newfName := filepath.Join(".zero2Pipe", file, f)

		newf, _ := os.Create(newfName)
		data, err := Asset(absf)
		if err != nil {
			log.Fatal(err)
		}

		_, err = newf.Write(data)
		if err != nil {
			log.Fatal(err.Error())
		}
		newf.Close()

	}

	return nil
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
	registerGenerator("AWSLambdaCodeExporter", AWSLambdaCodeExporter)
}

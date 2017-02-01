package domain

import (
	"fmt"
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
	os.MkdirAll(filepath.Join(".mistri", file), 0777)

	for _, f := range lambdaDir {
		absf := filepath.Join(".", file, f)
		_, err := AssetInfo(absf)
		if err != nil {
			export(absf)
			continue
		}

		newfName := filepath.Join(".mistri", file, f)

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

func AWSTemplateGenerator(conf map[string]string) error {
	templateName := conf["Template"]
	templates := template.Must(template.ParseGlob(templateName))
	fileExtension := filepath.Ext(templateName)
	renderedfileName := templateName[0 : len(templateName)-len(fileExtension)]
	fmt.Printf("Creating %s\n", filepath.Base(renderedfileName))
	fo, err := os.Create(filepath.Join(conf["FunctionSource"], filepath.Base(renderedfileName)))
	if err != nil {
		return err
	}

	err = templates.Execute(fo, conf)
	if err != nil {
		return err
	}

	return nil
}

type CustomActor struct {
	Generator Generator
	params    map[string]string
}

func (c *CustomActor) Run(interface{}) error {
	return c.Generator(c.params)
}

func registerGenerator(name string, gen Generator) {
	generators[name] = gen
}

func init() {
	registerGenerator("AWSTemplateGenerator", AWSTemplateGenerator)
	registerGenerator("AWSLambdaCodeExporter", AWSLambdaCodeExporter)
}

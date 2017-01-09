package client

import (
	"log"

	"github.com/antoniou/zero2Pipe/domain"
)

// CreateCommand is a Command implementation that
// launches the pipeline creation
type CreateCommand struct {
	BaseCommand
}

// NewCreateCommand returns a CreateCommand instance
func NewCreateCommand() *CreateCommand {

	return &CreateCommand{
		BaseCommand{
			Name:        "create",
			Usage:       "create a CD Pipeline",
			Description: "Create a CI/CD Pipeline",
			ArgsUsage:   "<pipeline>",
		},
	}
}

func (c *CreateCommand) Run(args []string) error {
	if len(args) == 0 {
		log.Fatalf("The %s command expects at least one argument", c.Name)
	}

	pipeline, err := domain.NewPipeline(map[string]string{
		"provider": "AWS_CP",
		"name":     args[0],
	})
	if err != nil {
		log.Fatal(err)
	}

	if err := pipeline.Create(args); err != nil {
		log.Fatal(err)
	}

	return nil
}

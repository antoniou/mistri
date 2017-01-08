package client

import (
	"fmt"

	"github.com/antoniou/zero2Pipe/domain"
)

// CreateCommand is a Command implementation that
// launches the pipeline creation
type CreateCommand struct {
	BaseCommand
}

// New returns a CreateCommand instance
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

func (c *CreateCommand) Run(args []string) int {
	pipeline := domain.AWSCodePipeline{}
	fmt.Println(pipeline.Create(args))
	return 0
}

func (c *CreateCommand) Help() string {
	return "Create command Help"
}

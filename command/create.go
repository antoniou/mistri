package command

import (
	"fmt"

	"github.com/antoniou/zero2Pipe/domain"
)

// CreateCommand is a Command implementation that
// launches the pipeline creation
type CreateCommand struct {
}

func (c *CreateCommand) Run(args []string) int {
	aws := domain.AWSCodePipeline{}
	fmt.Println(aws.Create(args))
	return 0
}

func (c *CreateCommand) Help() string {
	return "Create command Help"
}

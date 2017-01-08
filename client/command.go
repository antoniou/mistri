package client

// A Command is a runnable sub-command of a CLI.
type Command interface {
	Help() string
	Run(args []string) int
}

type BaseCommand struct {
	Name        string
	Description string
	Usage       string
	ArgsUsage   string
	Action      func(interface{}) error
}

func (c *BaseCommand) Run(args []string) int {
	return 0
}

func (c *BaseCommand) Help() string {
	return "Help"
}

package command

// A command is a runnable sub-command of a CLI.
type Command interface {
	Help() string
	Run(args []string) int
}

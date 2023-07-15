package commands

type CommandExecutor struct {
	command     Command
	currentPath string
}

func (c *CommandExecutor) ExecuteCommand() (int, error) {
	return c.command.Execute()
}

func (c *CommandExecutor) SetCommand(command Command) {
	c.command = command
}

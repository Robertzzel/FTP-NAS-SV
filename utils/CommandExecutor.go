package utils

import "FTP-NAS-SV/commands"

type CommandExecutor struct {
	command     commands.Command
	currentPath string
}

func (c *CommandExecutor) ExecuteCommand() (int, error) {
	return c.command.Execute()
}

func (c *CommandExecutor) SetCommand(command commands.Command) {
	c.command = command
}

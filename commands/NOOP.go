package commands

import "FTP-NAS-SV/codes"

type NOOP struct{}

func NewNOOPCommand() NOOP {
	return NOOP{}
}

func (cmd NOOP) Execute() (int, error) {
	return codes.CommandOkay, nil
}

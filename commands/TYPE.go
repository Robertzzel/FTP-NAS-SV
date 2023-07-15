package commands

import (
	"FTP-NAS-SV/codes"
	"FTP-NAS-SV/utils"
)

type TYPE struct {
	parameters       []string
	user             *utils.User
	transmissionType *int32
}

func NewTYPECommand(params []string, transmissionType *int32, user *utils.User) TYPE {
	return TYPE{parameters: params, transmissionType: transmissionType, user: user}
}

func (cmd TYPE) Execute() (int, error) {
	if !cmd.user.IsLogenIn() {
		return codes.NotLoggedIn, nil
	}
	if len(cmd.parameters) != 2 {
		return codes.SyntaxErrorParametersArguments, nil
	}

	switch cmd.parameters[1] {
	case "A":
		*cmd.transmissionType = codes.ASCII
	case "I":
		*cmd.transmissionType = codes.Image
	default:
		return codes.SyntaxErrorParametersArguments, nil
	}

	return codes.CommandOkay, nil
}

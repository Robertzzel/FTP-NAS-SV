package commands

import (
	"FTP-NAS-SV/codes"
	"FTP-NAS-SV/utils"
	"os"
	"path"
	"strings"
)

type RNTO struct {
	parameters   []string
	currentPath  string
	user         *utils.User
	selectedFile *string
}

func NewRNTOCommand(parameteres []string, currentPath string, user *utils.User, selectedFile *string) RNTO {
	return RNTO{
		parameters:   parameteres,
		currentPath:  currentPath,
		user:         user,
		selectedFile: selectedFile,
	}
}

func (cmd RNTO) Execute() (int, error) {
	if !cmd.user.IsLogenIn() {
		return codes.NotLoggedIn, nil
	}
	if len(cmd.parameters) != 2 {
		return codes.SyntaxErrorParametersArguments, nil
	}
	if *cmd.selectedFile == "" {
		return codes.BadSequenceOfCommands, nil
	}
	if !(strings.HasPrefix(cmd.parameters[1], "./") || strings.HasPrefix(cmd.parameters[1], "/")) {
		return codes.SyntaxErrorParametersArguments, nil
	}

	var filepath string
	if path.IsAbs(cmd.parameters[1]) {
		filepath = cmd.parameters[1]
	} else {
		filepath = path.Join(cmd.currentPath, cmd.parameters[1])
	}
	if !strings.HasPrefix(filepath, cmd.user.BasePath) {
		return codes.RequestedActionNotTakenFileName, nil
	}

	if err := os.Rename(*cmd.selectedFile, filepath); err != nil {
		return codes.RequestedActionNotTakenFileName, nil
	}

	*cmd.selectedFile = ""
	return codes.RequestedFileActionOkayCompleted, nil
}

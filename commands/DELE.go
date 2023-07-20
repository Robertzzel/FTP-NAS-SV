package commands

import (
	"FTP-NAS-SV/codes"
	"FTP-NAS-SV/utils"
	"os"
	"path"
	"strings"
)

type DELE struct {
	parameters  []string
	currentPath string
	user        *utils.User
}

func NewDELECommand(parameteres []string, currentPath string, user *utils.User) DELE {
	return DELE{
		parameters:  parameteres,
		currentPath: currentPath,
		user:        user,
	}
}

func (cmd DELE) Execute() (int, error) {
	if !cmd.user.IsLogenIn() {
		return codes.NotLoggedIn, nil
	}
	if len(cmd.parameters) != 2 {
		return codes.SyntaxErrorParametersArguments, nil
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
		return codes.RequestedActionNotTaken, nil
	}

	fileInfo, err := os.Stat(filepath)
	if err != nil {
		return codes.RequestedActionNotTaken, nil
	}
	if fileInfo.IsDir() {
		return codes.RequestedActionNotTaken, nil
	}

	if err := os.Remove(filepath); err != nil {
		return codes.RequestedActionNotTaken, nil
	}

	return codes.RequestedFileActionOkayCompleted, nil
}

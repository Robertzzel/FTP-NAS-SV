package commands

import (
	"FTP-NAS-SV/codes"
	"FTP-NAS-SV/utils"
	"os"
	"path"
	"strings"
)

type RNFR struct {
	parameters   []string
	currentPath  string
	user         *utils.User
	selectedFile *string
}

func NewRNFRCommand(parameteres []string, currentPath string, user *utils.User, selectedFile *string) RNFR {
	return RNFR{
		parameters:   parameteres,
		currentPath:  currentPath,
		user:         user,
		selectedFile: selectedFile,
	}
}

func (cmd RNFR) Execute() (int, error) {
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

	_, err := os.Stat(filepath)
	if err != nil {
		return codes.RequestedActionNotTaken, nil
	}

	*cmd.selectedFile = filepath

	return codes.RequestedFileActionPending, nil
}

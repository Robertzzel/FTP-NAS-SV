package commands

import (
	"FTP-NAS-SV/codes"
	"FTP-NAS-SV/utils"
	"os"
	"path"
	"strings"
)

type CDUP struct {
	currentPath *string
	parameters  []string
	user        *utils.User
}

func NewCDUPCommand(params []string, currentPath *string, user *utils.User) CDUP {
	return CDUP{
		parameters:  params,
		currentPath: currentPath,
		user:        user,
	}
}

func (cmd CDUP) Execute() (int, error) {
	if !cmd.user.IsLogenIn() {
		return codes.NotLoggedIn, nil
	}

	newCurrentPath := path.Dir(*cmd.currentPath)

	if !strings.HasPrefix(newCurrentPath, cmd.user.BasePath) {
		return codes.RequestedActionNotTaken, nil
	}
	if _, err := os.Stat(newCurrentPath); os.IsNotExist(err) {
		return codes.RequestedActionNotTaken, nil
	}

	*cmd.currentPath = newCurrentPath

	return codes.CommandOkay, nil
}

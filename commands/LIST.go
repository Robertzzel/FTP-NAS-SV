package commands

import "FTP-NAS-SV/utils"

type LIST struct {
	params      []string
	currentPath string
	user        *utils.User
}

func (cmd LIST) Execute() (int, error) {
	return -1, nil
}

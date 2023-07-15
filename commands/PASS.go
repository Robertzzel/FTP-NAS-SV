package commands

import (
	"FTP-NAS-SV/codes"
	"FTP-NAS-SV/database"
	"FTP-NAS-SV/utils"
	"path/filepath"
)

type PASS struct {
	parameters      []string
	user            *utils.User
	databaseManager database.DatabaseManager
	currentPath     *string
}

func NewPASSCommand(parameters []string, user *utils.User, dbManager database.DatabaseManager, currentPath *string) PASS {
	return PASS{
		parameters:      parameters,
		user:            user,
		databaseManager: dbManager,
		currentPath:     currentPath,
	}
}

func (cmd PASS) Execute() (int, error) {
	if cmd.user.IsLogenIn() {
		return codes.UserLoggedInProceed, nil
	}
	if len(cmd.parameters) != 2 {
		return codes.SyntaxErrorParametersArguments, nil
	}
	if cmd.user.Name == "" {
		return codes.BadSequenceOfCommands, nil
	}

	password := utils.Hash(cmd.parameters[1])
	isPasswordCorrect, err := cmd.databaseManager.Login(cmd.user.Name, password)
	if err != nil {
		return -1, err
	}

	if isPasswordCorrect {
		cmd.user.Password = password
		*cmd.currentPath = filepath.Join(*cmd.currentPath, cmd.user.Name)
		cmd.user.BasePath = *cmd.currentPath
		return codes.UserLoggedInProceed, nil
	}

	return codes.NeedAccountForLogin, nil
}

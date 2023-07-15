package commands

import (
	"FTP-NAS-SV/database"
	"FTP-NAS-SV/status_codes"
	"FTP-NAS-SV/utils"
)

type PASS struct {
	parameters      []string
	user            *utils.User
	databaseManager database.DatabaseManager
}

func NewPASSCommand(parameters []string, user *utils.User, dbManager database.DatabaseManager) PASS {
	return PASS{
		parameters:      parameters,
		user:            user,
		databaseManager: dbManager,
	}
}

func (cmd PASS) Execute() (int, error) {
	if cmd.user.IsLogenIn() {
		return status_codes.UserLoggedInProceed, nil
	}
	if len(cmd.parameters) != 2 {
		return status_codes.SyntaxErrorParametersArguments, nil
	}
	if cmd.user.Name == "" {
		return status_codes.BadSequenceOfCommands, nil
	}

	password := utils.Hash(cmd.parameters[1])
	isPasswordCorrect, err := cmd.databaseManager.Login(cmd.user.Name, password)
	if err != nil {
		return -1, err
	}

	if isPasswordCorrect {
		cmd.user.Password = password
		return status_codes.UserLoggedInProceed, nil
	}

	return status_codes.NeedAccountForLogin, nil
}

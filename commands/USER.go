package commands

import (
	"FTP-NAS-SV/database"
	"FTP-NAS-SV/status_codes"
	"FTP-NAS-SV/utils"
)

type USER struct {
	parameters      []string
	user            *utils.User
	databaseManager database.DatabaseManager
}

func NewUSERCommand(parameters []string, user *utils.User, dbManager database.DatabaseManager) USER {
	return USER{
		parameters:      parameters,
		user:            user,
		databaseManager: dbManager,
	}
}

func (cmd USER) Execute() (int, error) {
	if cmd.user.IsLogenIn() {
		return status_codes.UserLoggedInProceed, nil
	}
	if len(cmd.parameters) != 2 {
		return status_codes.SyntaxErrorParametersArguments, nil
	}
	userExists, err := cmd.databaseManager.CheckUsernameExists(cmd.parameters[1])
	if err != nil {
		return -1, err
	}

	if userExists {
		cmd.user.Name = cmd.parameters[1]
		return status_codes.UserNameOkayNeedPassword, nil
	}

	return status_codes.NeedAccountForLogin, nil
}

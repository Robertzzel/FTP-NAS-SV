package commands

import (
	"FTP-NAS-SV/codes"
	"FTP-NAS-SV/connection_management"
	"FTP-NAS-SV/utils"
	"fmt"
)

type PWD struct {
	currentPath string
	conn        *connection_management.TcpConnectionWrapper
	user        *utils.User
}

func NewPWDCommand(conn *connection_management.TcpConnectionWrapper, user *utils.User, currentPath string) PWD {
	return PWD{
		currentPath: currentPath,
		conn:        conn,
		user:        user,
	}
}

func (cmd PWD) Execute() (int, error) {
	if !cmd.user.IsLogenIn() {
		return codes.NotLoggedIn, nil
	}

	if err := cmd.conn.WriteMessage([]byte(fmt.Sprintf("%d\n%s", codes.PathnameCreated, cmd.currentPath))); err != nil {
		return codes.ServiceNotAvailable, err
	}

	return -1, nil
}

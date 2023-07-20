package commands

import (
	"FTP-NAS-SV/codes"
	. "FTP-NAS-SV/connection_management"
	"FTP-NAS-SV/utils"
	"fmt"
	"github.com/google/uuid"
	"os"
	"path"
)

type STOU struct {
	parameters  []string
	currentPath string
	user        *utils.User
	controlConn TcpConnectionWrapper
}

func NewSTOUCommand(parameteres []string, currentPath string, controlConn TcpConnectionWrapper, user *utils.User) STOU {
	return STOU{
		parameters:  parameteres,
		currentPath: currentPath,
		user:        user,
		controlConn: controlConn,
	}
}

func (cmd STOU) Execute() (int, error) {
	if !cmd.user.IsLogenIn() {
		return codes.NotLoggedIn, nil
	}
	if cmd.user.PassiveModeState() != codes.ClientConnected {
		return codes.CantOpenDataConnection, nil
	}

	filepath := path.Join(cmd.currentPath, "./"+uuid.New().String())
	fileDescriptor, err := os.Create(filepath)
	if err != nil {
		return codes.RequestedActionNotTaken, nil
	}

	if err = cmd.controlConn.WriteStatusCode(codes.DataConnectionAlreadyOpen); err != nil {
		return codes.ServiceNotAvailable, err
	}
	if err = ReadMessageToFile(cmd.user.DTConnection, fileDescriptor); err != nil {
		return codes.ConnectionClosedTransferAborted, err
	}
	if err = cmd.controlConn.WriteMessage(
		[]byte(fmt.Sprintf(
			"%d\n%s",
			codes.ClosingDataConnection, filepath,
		))); err != nil {
		return codes.ServiceNotAvailable, err
	}

	_ = cmd.user.ClosePassiveMode()
	return -1, nil
}

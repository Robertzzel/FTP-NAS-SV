package commands

import (
	"FTP-NAS-SV/codes"
	. "FTP-NAS-SV/connection_management"
	"FTP-NAS-SV/utils"
	"os"
	"path"
	"strings"
)

type RETR struct {
	parameters  []string
	currentPath string
	user        *utils.User
	controlConn TcpConnectionWrapper
}

func NewRETRCommand(parameteres []string, currentPath string, controlConn TcpConnectionWrapper, user *utils.User) RETR {
	return RETR{
		parameters:  parameteres,
		currentPath: currentPath,
		user:        user,
		controlConn: controlConn,
	}
}

func (cmd RETR) Execute() (int, error) {
	if !cmd.user.IsLogenIn() {
		return codes.NotLoggedIn, nil
	}
	if cmd.user.PassiveModeState() != codes.ClientConnected {
		return codes.CantOpenDataConnection, nil
	}
	if len(cmd.parameters) != 2 {
		return codes.SyntaxErrorParametersArguments, nil
	}

	var filepath string
	if !(strings.HasPrefix(cmd.parameters[1], "./") || strings.HasPrefix(cmd.parameters[1], "/")) {
		return codes.SyntaxErrorParametersArguments, nil
	}
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
	fileDescriptor, err := os.Open(filepath)
	if err != nil {
		return codes.RequestedActionNotTaken, nil
	}

	if err = cmd.controlConn.WriteStatusCode(codes.DataConnectionAlreadyOpen); err != nil {
		return codes.ServiceNotAvailable, err
	}
	if err = WriteFileMessage(cmd.user.DTConnection, fileDescriptor); err != nil {
		return codes.ConnectionClosedTransferAborted, err
	}
	if err = cmd.controlConn.WriteStatusCode(codes.ClosingDataConnection); err != nil {
		return codes.ServiceNotAvailable, err
	}

	_ = cmd.user.ClosePassiveMode()
	return -1, nil
}

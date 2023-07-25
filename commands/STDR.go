package commands

import (
	"FTP-NAS-SV/codes"
	. "FTP-NAS-SV/connection_management"
	"FTP-NAS-SV/utils"
	uuid "github.com/google/uuid"
	"os"
	"path"
	"strings"
)

type STDR struct {
	parameters  []string
	currentPath string
	user        *utils.User
	controlConn TcpConnectionWrapper
}

func (cmd STDR) Execute() (int, error) {
	if !cmd.user.IsLogenIn() {
		return codes.NotLoggedIn, nil
	}
	if cmd.user.PassiveModeState() != codes.ClientConnected {
		return codes.CantOpenDataConnection, nil
	}
	if len(cmd.parameters) != 2 {
		return codes.SyntaxErrorParametersArguments, nil
	}

	var dirPath string
	if !(strings.HasPrefix(cmd.parameters[1], "./") || strings.HasPrefix(cmd.parameters[1], "/")) {
		return codes.SyntaxErrorParametersArguments, nil
	}
	if path.IsAbs(cmd.parameters[1]) {
		dirPath = cmd.parameters[1]
	} else {
		dirPath = path.Join(cmd.currentPath, cmd.parameters[1])
	}
	if !strings.HasPrefix(dirPath, cmd.user.BasePath) {
		return codes.RequestedActionNotTaken, nil
	}

	pathSplitted := strings.Split(dirPath, "/")
	dirname := pathSplitted[len(pathSplitted)-1]
	if err := os.MkdirAll(dirPath+"/", 0770); err != nil {
		return codes.RequestedActionNotTaken, nil
	}

	randomName := path.Join(cmd.currentPath, uuid.New().String())
	fileDescriptor, err := os.Create(randomName)
	if err != nil {
		return codes.RequestedActionNotTaken, nil
	}
	defer func() {
		_ = fileDescriptor.Close()
		_ = os.Remove(randomName)
	}()

	if err = cmd.controlConn.WriteStatusCode(codes.DataConnectionAlreadyOpen); err != nil {
		return codes.ServiceNotAvailable, err
	}
	if err = ReadMessageToFile(cmd.user.DTConnection, fileDescriptor); err != nil {
		return codes.ConnectionClosedTransferAborted, err
	}

	if err = utils.Unzip(randomName, path.Join(cmd.currentPath, dirname)); err != nil {
		return codes.RequestedActionNotTaken, nil
	}

	if err = cmd.controlConn.WriteStatusCode(codes.ClosingDataConnection); err != nil {
		return codes.ServiceNotAvailable, err
	}

	_ = cmd.user.ClosePassiveMode()
	return -1, nil
}

func NewSTDRCommand(parameteres []string, currentPath string, controlConn TcpConnectionWrapper, user *utils.User) STDR {
	return STDR{
		parameters:  parameteres,
		currentPath: currentPath,
		user:        user,
		controlConn: controlConn,
	}
}

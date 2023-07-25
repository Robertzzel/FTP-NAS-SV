package commands

import (
	"FTP-NAS-SV/codes"
	. "FTP-NAS-SV/connection_management"
	"FTP-NAS-SV/utils"
	"github.com/google/uuid"
	"log"
	"os"
	"path"
	"strings"
)

type REDR struct {
	currentPath string
	parameters  []string
	user        *utils.User
	controlConn TcpConnectionWrapper
}

func NewREDRCommand(params []string, currentPath string, controlConn TcpConnectionWrapper, user *utils.User) REDR {
	return REDR{parameters: params, currentPath: currentPath, user: user, controlConn: controlConn}
}

func (cmd REDR) Execute() (int, error) {
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
	if !fileInfo.IsDir() {
		return codes.RequestedActionNotTaken, nil
	}

	outputPath := path.Join(cmd.currentPath, uuid.New().String())
	if err = utils.Zip(filepath, outputPath); err != nil {
		return codes.RequestedActionNotTaken, nil
	}
	defer func() {
		_ = os.Remove(outputPath)
	}()
	fileDescriptor, err := os.Open(outputPath)
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

	if err = os.Remove(outputPath); err != nil {
		log.Println(outputPath, "could not be removed.")
	}

	_ = cmd.user.ClosePassiveMode()
	return -1, nil
}

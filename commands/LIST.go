package commands

import (
	"FTP-NAS-SV/codes"
	. "FTP-NAS-SV/connection_management"
	"FTP-NAS-SV/utils"
	"encoding/json"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type LIST struct {
	parameters  []string
	currentPath string
	user        *utils.User
	controlConn TcpConnectionWrapper
}

func NewLISTCommand(parameteres []string, currentPath string, controlConn TcpConnectionWrapper, user *utils.User) LIST {
	return LIST{
		parameters:  parameteres,
		currentPath: currentPath,
		user:        user,
		controlConn: controlConn,
	}
}

func (cmd LIST) Execute() (int, error) {
	if !cmd.user.IsLogenIn() {
		return codes.NotLoggedIn, nil
	}
	if cmd.user.PassiveModeState() != codes.ClientConnected {
		return codes.CantOpenDataConnection, nil
	}

	directoryPath := cmd.currentPath
	if len(cmd.parameters) == 2 {
		if !(strings.HasPrefix(cmd.parameters[1], "./") || strings.HasPrefix(cmd.parameters[1], "/")) {
			return codes.SyntaxErrorParametersArguments, nil
		}
		if path.IsAbs(cmd.parameters[1]) {
			directoryPath = cmd.parameters[1]
		} else {
			directoryPath = path.Join(cmd.currentPath, cmd.parameters[1])
		}
		if !strings.HasPrefix(directoryPath, cmd.user.BasePath) {
			return codes.RequestedActionNotTaken, nil
		}
	}

	fileInfo, err := os.Stat(directoryPath)
	if err != nil {
		return codes.RequestedActionAborted, nil
	}

	if fileInfo.IsDir() {
		files, err := os.ReadDir(directoryPath)
		if err != nil {
			return codes.RequestedActionAborted, err
		}

		var contents []utils.FileDetails
		for _, file := range files {
			fileDetails := utils.FileDetails{Size: 0, Name: file.Name(), IsDir: file.IsDir()}

			fileType, _ := utils.GetFileType(filepath.Join(directoryPath, file.Name()))
			fileDetails.Type = fileType
			if strings.Contains(fileType, "image") {
				fileDetails.ImageData, err = utils.Resize(filepath.Join(directoryPath, file.Name()))
				if err != nil {
					fileDetails.ImageData = nil
				}
			}

			info, err := file.Info()
			if err != nil {
				fileDetails.Size = -1
			} else {
				fileDetails.Size = info.Size()
			}
			contents = append(contents, fileDetails)
		}

		var sendData []byte
		if contents != nil {
			sendData, err = json.Marshal(contents)
			if err != nil {
				return codes.RequestedActionAborted, err
			}
		} else {
			sendData = []byte("")
		}

		if err := cmd.controlConn.WriteStatusCode(codes.DataConnectionAlreadyOpen); err != nil {
			return codes.ServiceNotAvailable, err
		}

		if err := WriteMessage(cmd.user.DTConnection, sendData); err != nil {
			return codes.ConnectionClosedTransferAborted, err
		}
	} else {
		fileType, _ := utils.GetFileType(directoryPath)
		fileDetails := utils.FileDetails{
			Name:  fileInfo.Name(),
			Size:  fileInfo.Size(),
			IsDir: fileInfo.IsDir(),
			Type:  fileType,
		}
		sendData, err := json.Marshal(fileDetails)
		if err != nil {
			return codes.RequestedActionAborted, nil
		}

		if err = cmd.controlConn.WriteStatusCode(codes.DataConnectionAlreadyOpen); err != nil {
			return codes.ServiceNotAvailable, err
		}

		if err = WriteMessage(cmd.user.DTConnection, sendData); err != nil {
			return codes.ConnectionClosedTransferAborted, err
		}
	}

	if err := cmd.controlConn.WriteStatusCode(codes.ClosingDataConnection); err != nil {
		return codes.ServiceNotAvailable, err
	}

	_ = cmd.user.ClosePassiveMode()

	return -1, nil
}

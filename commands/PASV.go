package commands

import (
	"FTP-NAS-SV/codes"
	"FTP-NAS-SV/connection_management"
	"FTP-NAS-SV/utils"
	"fmt"
	"strconv"
	"strings"
)

type PASV struct {
	user        *utils.User
	conn        *connection_management.TcpConnectionWrapper
	connManager *connection_management.ConnectionManager
}

func NewPASVCommand(conn *connection_management.TcpConnectionWrapper, connManager *connection_management.ConnectionManager, user *utils.User) PASV {
	return PASV{
		user:        user,
		conn:        conn,
		connManager: connManager,
	}
}

func (cmd PASV) Execute() (int, error) {
	if !cmd.user.IsLogenIn() {
		return codes.NotLoggedIn, nil
	}
	passiveModeState := cmd.user.PassiveModeState()
	if passiveModeState == codes.ClientConnected || passiveModeState == codes.WaitingForClientToConnect {
		_ = cmd.user.ClosePassiveMode()
	}

	listener, err := cmd.connManager.ListenToAvailablePort()
	if err != nil {
		return codes.ServiceNotAvailable, err
	}

	publicIp, err := ipStringToFTPFormat(cmd.connManager.GetPublicIP())
	if err != nil {
		return codes.ServiceNotAvailable, err
	}
	port, err := strconv.Atoi(strings.Split(listener.Addr().String(), ":")[1])
	if err != nil {
		return codes.ServiceNotAvailable, err
	}
	err = cmd.conn.WriteMessage(
		[]byte(
			fmt.Sprintf(
				"%d Entering Passive Mode %d,%d,%d,%d,%d,%d",
				codes.EnteringPassiveMode,
				publicIp[0], publicIp[1], publicIp[2], publicIp[3],
				port/256, port%256,
			),
		),
	)
	if err != nil {
		return codes.ServiceNotAvailable, err
	}

	cmd.user.SetUserPASVMode(listener)
	return -1, nil
}

func ipStringToFTPFormat(ip string) ([4]int, error) {
	var rezIp [4]int
	var err error

	ipSliced := strings.Split(ip, ".")
	for i := 0; i < 4; i++ {
		rezIp[i], err = strconv.Atoi(ipSliced[i])
		if err != nil {
			return rezIp, err
		}
	}
	return rezIp, nil
}

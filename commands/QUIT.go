package commands

import (
	"FTP-NAS-SV/codes"
	"FTP-NAS-SV/connection_management"
)

type QUIT struct {
	conn *connection_management.TcpConnectionWrapper
}

func NewQUITCommand(conn *connection_management.TcpConnectionWrapper) QUIT {
	return QUIT{conn: conn}
}

func (cmd QUIT) Execute() (int, error) {
	if err := cmd.conn.WriteStatusCode(codes.ServiceClosingControlConnection); err != nil {
		return -1, err
	}

	return -1, cmd.conn.Close()
}

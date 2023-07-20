package utils

import (
	"FTP-NAS-SV/codes"
	"crypto/tls"
	"log"
	"net"
)

type User struct {
	Name         string
	Password     string
	BasePath     string
	DTListener   net.Listener
	DTConnection net.Conn
}

func (user *User) IsLogenIn() bool {
	return user.Name != "" && user.Password != ""
}

func (user *User) SetUserPASVMode(dTListener net.Listener) {
	go func() {
		var err error
		user.DTListener = dTListener
		user.DTConnection, err = dTListener.Accept()
		if err != nil {
			_ = user.ClosePassiveMode()
			log.Println("Error while connecting:", err)
			return
		}
		if err = user.DTConnection.(*tls.Conn).Handshake(); err != nil {
			log.Println("Handshake error:", err)
		}
	}()
}

func (user *User) PassiveModeState() int {
	if user.DTListener == nil {
		return codes.NotInitiated
	}
	if user.DTConnection == nil {
		return codes.WaitingForClientToConnect
	}
	return codes.ClientConnected
}

func (user *User) ClosePassiveMode() error {
	if user.DTListener != nil {
		if err := user.DTListener.Close(); err != nil {
			return err
		}
		user.DTListener = nil
	}
	if user.DTConnection != nil {
		if err := user.DTConnection.Close(); err != nil {
			return err
		}
		user.DTConnection = nil
	}
	return nil
}

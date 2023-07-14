package main

import (
	"FTP-NAS-SV/commands"
	. "FTP-NAS-SV/connection_management"
	"FTP-NAS-SV/database"
	"FTP-NAS-SV/status_codes"
	"FTP-NAS-SV/utils"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

func handleClientConnection(conn TcpConnectionWrapper, dbManager database.DatabaseManager) error {
	user := utils.User{}
	commandExecutor := utils.CommandExecutor{}

	if err := conn.WriteStatusCode(status_codes.ServiceReadyForNewUser); err != nil {
		return errors.New(fmt.Sprintln("Error on", conn.RemoteAddr().String(), ",err : ", err))
	}

	for {
		message, err := conn.ReadMessage()
		if err != nil {
			return errors.New(fmt.Sprintln("Error on", conn.RemoteAddr().String(), ",err : ", err))
		}

		messageComponents := strings.Split(string(message), " ")
		switch messageComponents[0] {
		case "USER":
			cmd := commands.NewUSERCommand(messageComponents, user, dbManager)
			commandExecutor.SetCommand(cmd)
			statusCode, err := commandExecutor.ExecuteCommand()
			if err != nil {
				_ = conn.WriteStatusCode(status_codes.ServiceNotAvailable)
				continue
			}
			_ = conn.WriteStatusCode(statusCode)
		case "PASS":
			break
		case "CWD":
			break
		case "CDUP":
			break
		case "QUIT":
			break
		case "TYPE":
			break
		case "RMD":
			break
		case "MKD":
			break
		case "PWD":
			break
		case "LIST":
			break
		case "STAT":
			break
		case "HELP":
			break
		case "RETR":
			break
		case "STOR":
			break
		case "STOU":
			break
		case "RNFR":
			break
		case "RNTO":
			break
		case "DELE":
			break
		case "PORT":
			break
		case "ALLO":
			break
		case "NOOP":
			break
		}
	}
}

func main() {
	certFile := os.Args[1] // server.crt
	keyFile := os.Args[2]  // server.key

	dbManager, err := database.NewDatabase()
	if err != nil {
		panic(err)
	}

	connManager, err := NewConnectionManager(certFile, keyFile)
	if err != nil {
		panic(err)
	}

	listener, err := connManager.ListenForClients()
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection", err)
		}

		go func() {
			err := handleClientConnection(TcpConnectionWrapper{Conn: conn}, dbManager)
			if err != nil {
				log.Println(err)
			}
		}()
	}
}

package main

import (
	"FTP-NAS-SV/codes"
	"FTP-NAS-SV/commands"
	. "FTP-NAS-SV/connection_management"
	"FTP-NAS-SV/database"
	"FTP-NAS-SV/utils"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	//UsbDrivePath   = "/home/robert/usb_drives/usb1/"
	UsbDrivePath = "/home/robert/Downloads"
)

func handleClientConnection(conn TcpConnectionWrapper, dbManager database.DatabaseManager, connManager *ConnectionManager) error {
	user := utils.User{}
	commandExecutor := commands.CommandExecutor{}
	currentPath := UsbDrivePath
	transmissionType := codes.Image

	if err := conn.WriteStatusCode(codes.ServiceReadyForNewUser); err != nil {
		return errors.New(fmt.Sprintln("Error on", conn.RemoteAddr().String(), ",err : ", err))
	}

	for {
		message, err := conn.ReadMessage()
		if err != nil {
			return errors.New(fmt.Sprintln("Error on", conn.RemoteAddr().String(), ",err : ", err))
		}

		messageComponents := strings.Split(string(message), " ")
		var cmd commands.Command
		switch messageComponents[0] {
		case "USER":
			cmd = commands.NewUSERCommand(messageComponents, &user, dbManager)
		case "PASS":
			cmd = commands.NewPASSCommand(messageComponents, &user, dbManager, &currentPath)
		case "CWD":
			cmd = commands.NewCWDCommand(messageComponents, &currentPath, &user)
		case "CDUP":
			cmd = commands.NewCDUPCommand(messageComponents, &currentPath, &user)
		case "QUIT":
			cmd = commands.NewQUITCommand(&conn)
		case "TYPE":
			cmd = commands.NewTYPECommand(messageComponents, &transmissionType, &user)
		case "RMD":
			cmd = commands.NewRMDCommand(messageComponents, currentPath, &user)
		case "MKD":
			cmd = commands.NewMKDCommand(messageComponents, currentPath, &user)
		case "PWD":
			cmd = commands.NewPWDCommand(&conn, &user, currentPath)
		case "PASV":
			cmd = commands.NewPASVCommand(&conn, connManager, &user)
		case "LIST":
			cmd = commands.NewLISTCommand(messageComponents, currentPath, conn, &user)
		case "STAT":
			break
		case "HELP":
			break
		case "RETR":
			cmd = commands.NewRETRCommand(messageComponents, currentPath, conn, &user)
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
			cmd = commands.NewNOOPCommand()
		default:
			_ = conn.WriteStatusCode(codes.SyntaxErrorCommandUnrecognized)
			continue
		}

		commandExecutor.SetCommand(cmd)
		statusCode, err := commandExecutor.ExecuteCommand()

		if statusCode == -1 && err == nil {
			continue
		}

		if err == nil {
			_ = conn.WriteStatusCode(statusCode)
		} else {
			_ = conn.WriteStatusCode(codes.ServiceNotAvailable)
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

	listener, err := connManager.ListenForClientsPI()
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	fmt.Println("listening on ", listener.Addr().String())
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection", err)
		}

		go func() {
			err := handleClientConnection(TcpConnectionWrapper{Conn: conn}, dbManager, &connManager)
			if err != nil {
				log.Println(err)
			}
		}()
	}
}

package main

import (
	database "FTP-NAS-SV/database"
	"FTP-NAS-SV/utils"
	"log"

	"os"
)

func main1() {
	if len(os.Args) < 4 {
		log.Fatal("go run register.go <username> <password> <email>")
	}

	db, err := database.NewDatabase()
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("insert into User (Name, Password, Email) values (?, ?, ?)", os.Args[1], utils.Hash(os.Args[2]), os.Args[3])
	if err != nil {
		log.Fatal(err)
	}
}

package database

import (
	"FTP-NAS-SV/utils"
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
)

type DatabaseManager struct {
	*sql.DB
}

func NewDatabase() (DatabaseManager, error) {
	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		return DatabaseManager{}, err
	}

	dm := DatabaseManager{db}
	if err = dm.migrateDatabase(); err != nil {
		return DatabaseManager{}, err
	}

	return dm, nil
}

func (db *DatabaseManager) Login(username, password string) (bool, error) {
	var cnt int
	err := db.QueryRow(`select count(*) from User where Name = ? and Password = ? LIMIT 1`, username, utils.Hash(password)).Scan(&cnt)
	if err != nil {
		return false, errors.New("database problem")
	}
	return cnt != 0, nil
}

func (db *DatabaseManager) CheckUsernameExists(username string) (bool, error) {
	var cnt int
	err := db.QueryRow(`select count(*) from User where Name = ? LIMIT 1`, username).Scan(&cnt)
	if err != nil {
		return false, errors.New("database problem: " + err.Error())
	}
	return cnt == 0, nil
}

func (db *DatabaseManager) migrateDatabase() error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS User(
    	Id integer PRIMARY KEY,
		Name varchar(255) UNIQUE NOT NULL,
		Email varchar(255),
		Password varchar(255) NOT NULL
    )`)
	if err != nil {
		return err
	}
	return nil
}

func (db *DatabaseManager) Close() {
	db.Close()
}

package mysql

import (
	"database/sql"
	"fmt"

	// go sql driver
	_ "github.com/go-sql-driver/mysql"
)

// Database means database
type Database struct {
	db *sql.DB
}

// Connect to database
func Connect(host string, port uint16, user, password string) (Database, error) {
	str := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, port, DatabaseName)
	db, err := sql.Open("mysql", str)
	return Database{db}, err
}

// Initialize database, create database and tables
func (d Database) Initialize() error {
	if err := d.createDatabase(); err != nil {
		return err
	}
	if err := d.createUserTable(); err != nil {
		return err
	}
	if err := d.createFileTable(); err != nil {
		return err
	}
	return nil
}

// DestroyTables drop user and file table, only available for dev
func (d Database) DestroyTables() error {
	if err := d.dropFileTable(); err != nil {
		return err
	}
	if err := d.dropUserTable(); err != nil {
		return err
	}
	return nil
}

func (d Database) createDatabase() error {
	_, err := d.db.Exec(createDatabase)
	return err
}

func (d Database) createUserTable() error {
	_, err := d.db.Exec(createUserTable)
	return err
}

func (d Database) dropUserTable() error {
	_, err := d.db.Exec("drop table " + userTable)
	return err
}

func (d Database) createFileTable() error {
	_, err := d.db.Exec(createFileTable)
	return err
}

func (d Database) dropFileTable() error {
	_, err := d.db.Exec("drop table " + fileTable)
	return err
}

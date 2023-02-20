package database

import (
	"database/sql"
	"fmt"
	"log"
)

var connString string

func GetDb() (conn *sql.DB, err error) {
	conn, err = sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatal("Open connection failed:", err.Error())
	}
	return conn, err
}
func SetParam(server, user, password, port, database string) {
	connString = fmt.Sprintf("server=%s;user id=%s;password=%s;port=%s;database=%s;",
		server, user, password, port, database)
}

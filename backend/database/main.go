package database

import (
	"database/sql"
	"fmt"
	"log"
)

func GetDb(connStr string) (conn *sql.DB, err error) {

	conn, err = sql.Open("sqlserver", connStr)
	if err != nil {
		log.Fatal("Open connection failed:", err.Error())
	}
	return conn, err
}
func SetParam(server, user, password, port, database string) string {
	return fmt.Sprintf("server=%s;user id=%s;password=%s;port=%s;database=%s;",
		server, user, password, port, database)
}

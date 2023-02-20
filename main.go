package main

import (
	"database/sql"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"log"
	"net/http"
	"os"
	"sqlutils/backend/route"
)

var server = "localhost"

//var port = 1433
var port = 1433
var user = "sa"
var password = "masterkey"
var database = "utils"

func main() {
	var err error

	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
		server, user, password, port, database)
	fmt.Println(server, user, password, database)
	conn, err := sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatal("Open connection failed:", err.Error())
	}
	fmt.Printf("Connected!\n")
	defer conn.Close()
	err = conn.Ping()
	if err != nil {
		log.Fatal(err.Error())
	}
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("/", route.Home)
	host, _ := os.Hostname()
	log.Println("Запуск сервера на http://" + host + ":8080")
	err = http.ListenAndServe(":8080", mux)

	if err != nil {
		log.Fatal(err.Error())
	}

	//query := `select name from test`
	//rows, err := conn.Query(query)
	//for rows.Next() {
	//	var t string
	//	rows.Scan(&t)
	//	fmt.Println(t)
	//}
	//
	////var (
	////	userid   = flag.String("U", "sa", "login_id")
	////	password = flag.String("P", "masterkey", "password")
	////	server   = flag.String("S", "192.168.1.2\\SQLEXPRESS", "server_name[\\instance_name]")
	////	database = flag.String("d", "utils", "db_name")
	////	port     = flag.String("port", "1433", "the database port")
	////)
	//flag.Parse()
	//
	////dsn := "server=" + *server + ";user id=" + *userid + ";password=" + *password + ";port=" + *port + ";database=" + *database
	////	dsn := "server=" + *server + ";user id=" + *userid + ";password=" + *password + ";port=" + *port + ";database=" + *database
	//connString := "sqlserver://sa:masterkey@home:1433/SQLEXPRESS?database=utils"
	////db, err := sql.Open("mssql", connString)
	//db, err := sql.Open("sqlserver", connString)
	//if err != nil {
	//	fmt.Println("Cannot connect: ", err.Error())
	//	return
	//}
	//err = db.Ping()
	//if err != nil {
	//	fmt.Println("Cannot connect: ", err.Error())
	//	return
	//}
	//defer db.Close()
	//
	//fmt.Printf("bye\n")
}

package main

import (
	_ "github.com/denisenkom/go-mssqldb"
	"log"
	"net/http"
	"os"
	"sqlutils/backend/route"
)

func main() {

	//var server = ""
	//var port = "1433"
	//var user = "sa"
	//var password = "masterkey"
	//var dbName = "utils"

	var err error

	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("/", route.Login)
	mux.HandleFunc("/main", route.DoLogin)
	mux.HandleFunc("/task1", route.Task1)

	//mux.HandleFunc("/login", route.Login)
	host, _ := os.Hostname()
	log.Println("Сервер запущен: http://" + host + ":8080")
	err = http.ListenAndServe(":8080", mux)

	if err != nil {
		log.Fatal(err.Error())
	}

}

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
	os.RemoveAll("tmp")
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("/", route.Login)
	mux.HandleFunc("/main", route.DoLogin)
	mux.HandleFunc("/task1", route.Task1Handler)
	mux.HandleFunc("/task-list", route.TaskListHandler)
	mux.HandleFunc("/catalog-list", route.CatalogListHandler)
	mux.HandleFunc("/task-create", route.TaskCreateHandler)
	mux.HandleFunc("/task-edit/", route.TaskEditHandler)

	mux.HandleFunc("/catalog-create", route.CatalogCreateHandler)

	//api
	mux.HandleFunc("/api/task-list", route.GetTaskListHandler)
	mux.HandleFunc("/api/task-params", route.GetTaskParamsHandler)
	mux.HandleFunc("/api/task-save", route.TaskSaveHandler)
	mux.HandleFunc("/api/task-save-params", route.TaskSaveParamsHandler)
	mux.HandleFunc("/api/upload", route.TaskUploadHandler)
	mux.HandleFunc("/api/task-exe", route.TaskExeHandler)
	mux.HandleFunc("/api/catalog-list", route.GetCatalogListHandler)
	mux.HandleFunc("/api/catalog-save", route.CatalogSaveHandler)
	//mux.HandleFunc("/login", route.Login)
	host, _ := os.Hostname()
	log.Println("Сервер запущен: http://" + host + ":8080")
	err = http.ListenAndServe(":8080", mux)

	if err != nil {
		log.Fatal(err.Error())
	}

}

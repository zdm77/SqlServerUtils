package main

import (
	"bufio"
	_ "github.com/denisenkom/go-mssqldb"
	"log"
	"net/http"
	"os"
	"sqlutils/backend/route"
	catalog_route "sqlutils/backend/route/catalog-route"
	catalog_work_route "sqlutils/backend/route/catalog-work-route"
	"sqlutils/backend/route/task-route"
	"strings"
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
	//*********Каталог****************////////////////
	//*********Задачи***//
	mux.HandleFunc("/task-list-catalog", task_route.TaskListCatalogHandler)
	mux.HandleFunc("/task-catalog-create", task_route.TaskCatalogCreateHandler)
	mux.HandleFunc("/task-catalog-edit/", task_route.TaskCatalogEditHandler)

	//****Каталог Задачи api
	mux.HandleFunc("/api/task-list", task_route.GetTaskCatalogListHandler)
	mux.HandleFunc("/api/task-params", task_route.GetTaskParamsHandler)
	mux.HandleFunc("/api/task-save", task_route.TaskCatalogSaveHandler)
	mux.HandleFunc("/api/task-save-params", task_route.TaskSaveParamsHandler)

	//********Справочники*************//
	mux.HandleFunc("/catalog-list", catalog_route.CatalogListHandler)
	mux.HandleFunc("/catalog-create", catalog_route.CatalogCreateHandler)
	mux.HandleFunc("/catalog-edit/", catalog_route.CatalogEditHandler)

	//*******Справочники api
	mux.HandleFunc("/api/catalog-list", catalog_route.GetCatalogListHandler)
	mux.HandleFunc("/api/catalog-save", catalog_route.CatalogSaveHandler)
	mux.HandleFunc("/api/get-db-fields", catalog_route.GetDbFieldsHandler)
	mux.HandleFunc("/api/save-db-fields", catalog_route.SaveDbFieldsHandler)

	mux.HandleFunc("/task1", task_route.Task1Handler)

	//Заполнение справочников
	mux.HandleFunc("/catalog-work-list/", catalog_work_route.CatalogWorkListHandler)
	mux.HandleFunc("/catalog-work-create/", catalog_work_route.CatalogWorkCreateHandler)
	mux.HandleFunc("/api/catalog-work-save", catalog_work_route.CatalogWorkSaveHandler)
	mux.HandleFunc("/api/catalog-work-list", catalog_work_route.CatalogGetWorkListHandler)

	mux.HandleFunc("/api/upload", task_route.TaskUploadHandler)
	mux.HandleFunc("/api/task-exe", task_route.TaskExeHandler)

	//mux.HandleFunc("/login", route.Login)
	host, _ := os.Hostname()

	file, err := os.Open("settings.txt")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer file.Close()
	sc := bufio.NewScanner(file)
	port := ""
	for sc.Scan() {
		cont := strings.Contains(sc.Text(), "WebPort")
		if cont {
			port = strings.Split(sc.Text(), "=")[1]
		}

	}
	file.Close()
	log.Println("Сервер запущен: http://" + host + ":" + port)
	err = http.ListenAndServe(":"+port, mux)

	if err != nil {
		log.Fatal(err.Error())
	}

}

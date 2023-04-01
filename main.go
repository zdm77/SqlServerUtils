package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"html/template"
	"log"
	"net/http"
	"os"
	"sqlutils/backend/database"
	"sqlutils/backend/route"
	catalog_route "sqlutils/backend/route/catalog-route"
	catalog_work_route "sqlutils/backend/route/catalog-work-route"
	script_route "sqlutils/backend/route/script-route"
	"sqlutils/backend/route/task-route"
	"sqlutils/backend/session"
	"strconv"
	"strings"
)

func main() {

	var err error
	os.RemoveAll("tmp")
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("/", route.Login)
	mux.HandleFunc("/main", route.DoLogin)
	mux.HandleFunc("/api/get-access", GetAccess)
	mux.HandleFunc("/test", Test)
	mux.HandleFunc("/test2", Test2)
	mux.HandleFunc("/api/test/", GetTestDataForLocal)
	//*********Каталог****************////////////////
	//*********Задачи***//
	mux.HandleFunc("/task-list-catalog", task_route.TaskListCatalogHandler)
	mux.HandleFunc("/task-catalog-create", task_route.TaskCatalogCreateHandler)
	mux.HandleFunc("/task-catalog-edit/", task_route.TaskCatalogEditHandler)

	//****Каталог Задачи api
	mux.HandleFunc("/api/task-list", task_route.GetTaskCatalogListHandler)
	mux.HandleFunc("/api/task-list/delete/", task_route.TaskCatalogDeleteListHandler)
	mux.HandleFunc("/api/task-params", task_route.GetTaskParamsHandler)
	mux.HandleFunc("/api/task-save", task_route.TaskCatalogSaveHandler)
	mux.HandleFunc("/api/task-save-params", task_route.TaskSaveParamsHandler)

	//********Справочники*************//
	mux.HandleFunc("/catalog-list/", catalog_route.CatalogListHandler)
	mux.HandleFunc("/catalog-create/", catalog_route.CatalogCreateHandler)
	mux.HandleFunc("/catalog-edit/", catalog_route.CatalogEditHandler)

	//*******Справочники api
	mux.HandleFunc("/api/get-catalog-fields/", catalog_route.GetCatalogFieldsHandler)
	mux.HandleFunc("/api/catalog-list/delete/", catalog_route.CatalogListDeleteListHandler)
	mux.HandleFunc("/api/catalog-list/", catalog_route.GetCatalogListHandler)
	mux.HandleFunc("/api/catalog-list-link/", catalog_route.GetLinkListHandler)
	mux.HandleFunc("/api/catalog-save", catalog_route.CatalogSaveHandler)
	mux.HandleFunc("/api/get-db-fields", catalog_route.GetDbFieldsHandler)
	mux.HandleFunc("/api/save-db-fields", catalog_route.SaveDbFieldsHandler)
	//Новое поле
	mux.HandleFunc("/api/create-db-field", catalog_route.CreateDbFieldHandler)

	mux.HandleFunc("/task1", task_route.Task1Handler)

	//Заполнение справочников
	mux.HandleFunc("/api/catalog-work-list/delete/", catalog_work_route.CatalogWorkListDeleteListHandler)
	mux.HandleFunc("/catalog-work-list/", catalog_work_route.CatalogWorkListHandler)
	mux.HandleFunc("/catalog-work-create/", catalog_work_route.CatalogWorkCreateHandler)
	mux.HandleFunc("/catalog-work-edit/", catalog_work_route.CatalogWorkEditHandler)
	mux.HandleFunc("/api/catalog-work-save", catalog_work_route.CatalogWorkSaveHandler)
	mux.HandleFunc("/api/catalog-work-list", catalog_work_route.CatalogGetWorkListHandler)

	//Скрипты
	//Справочник
	mux.HandleFunc("/catalog/script-list/", script_route.ScriptListHandler)
	mux.HandleFunc("/catalog/script-list-work/", script_route.ScriptListWorkHandler)
	mux.HandleFunc("/api/catalog/script-list/", script_route.GetCatalogScriptListHandler)
	mux.HandleFunc("/script-edit/", script_route.ScriptEditHandler)
	mux.HandleFunc("/api/script-save", script_route.ScriptSaveHandler)
	mux.HandleFunc("/api/script-exe/", script_route.ScriptExeHandler)

	//
	mux.HandleFunc("/api/upload", task_route.TaskUploadHandler)
	mux.HandleFunc("/api/task-exe", task_route.TaskExeHandler)
	//доступ
	mux.HandleFunc("/api/get-access-catalog-record", catalog_work_route.CatalogAccessRecordHandler)
	mux.HandleFunc("/api/save-access-catalog-record", catalog_work_route.SaveAccessRecordHandler)
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
func GetAccess(w http.ResponseWriter, r *http.Request) {
	user := session.GetSessionData(r)
	if user != nil {
		var data []byte
		if user.Login == user.SuperAdmin {
			data, _ = json.Marshal(route.Message{Text: "isAdmin"})
		} else {
			data, _ = json.Marshal(route.Message{Text: "isUser"})
		}
		w.Write(data)
	} else {
		data, _ := json.Marshal(route.Message{Text: "not-login"})
		w.Write(data)
	}
}
func Home(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "Первая страничка!")
}
func Test(w http.ResponseWriter, r *http.Request) {

	files := []string{
		"./ui/html/test.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/top.layout.tmpl",
		"./ui/html/controls/create.tmpl",
		"./ui/html/controls/table.tmpl",
		"./ui/html/controls/list-panel.tmpl",
	}

	tpl, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
	}

	tpl.Execute(w, nil)
}
func Test2(w http.ResponseWriter, r *http.Request) {
	user := session.GetSessionData(r)
	db, _ := database.GetDb(user.ConnString)
	defer db.Close()

	for i := 0; i < 100000; i++ {
		query := `insert into test2 (name) values ('test` + strconv.Itoa(i+1) + `')`

		_, err := db.Exec(query)
		if err != nil {
			log.Println(err.Error())
		}
	}

}
func GetTestData(w http.ResponseWriter, r *http.Request) {
	keys := r.URL.Query()

	user := session.GetSessionData(r)
	type JsonT struct {
		Start int `json:"start"`
	}

	st := keys["start"][0]
	order := keys["order[0][dir]"]
	orderCol := keys["order[0][column]"]
	fmt.Println(orderCol, order)
	start, _ := strconv.Atoi(st)
	//	decoder := json.NewDecoder(r.Body)
	param := JsonT{
		Start: start,
	}

	//decoder.Decode(&param)
	//log.Println(param)
	list := database.Test(user, param.Start)
	data, _ := json.Marshal(list)
	w.Write(data)
}
func GetTestDataForLocal(w http.ResponseWriter, r *http.Request) {
	//user := session.GetSessionData(r)
	//
	////list := database.Test(user, 0)
	//_, list := db_catalog_work.GetCatalogWorkListById(user, 5071, true)
	//
	//data, _ := json.Marshal(list)
	//w.Write(data)
}

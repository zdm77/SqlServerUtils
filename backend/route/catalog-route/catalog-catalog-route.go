package catalog_route

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"sqlutils/backend/database/db_catalog"
	"sqlutils/backend/model"
	"sqlutils/backend/route"
	"sqlutils/backend/session"
	"strconv"
)

func CatalogListHandler(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/catalog/catalog-list.page.tmpl",
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
func GetCatalogListHandler(w http.ResponseWriter, r *http.Request) {
	user := session.GetSessionData(r)
	if user != nil {

		list := db_catalog.GetCatalogList(user)
		data, _ := json.Marshal(list)
		w.Write(data)
	} else {
		data, _ := json.Marshal(route.Message{Text: "not-login"})
		w.Write(data)
	}
}
func CatalogCreateHandler(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/catalog/catalog-create.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/top.layout.tmpl",
		"./ui/html/controls/save.tmpl",
	}

	tpl, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
	}

	tpl.Execute(w, nil)
}
func CatalogSaveHandler(w http.ResponseWriter, r *http.Request) {
	user := session.GetSessionData(r)
	if user != nil {
		decoder := json.NewDecoder(r.Body)
		var param model.Catalog
		err := decoder.Decode(&param)
		if err != nil {
			log.Println(err.Error())
		} else {
			err, id := db_catalog.SaveCatalog(user, param)
			if err != nil && err.Error() != "sql: no rows in result set" {
				data, _ := json.Marshal(route.Message{Text: err.Error()})
				w.Write(data)
			} else {
				data, _ := json.Marshal(route.Message{Text: "ok-" + strconv.Itoa(id)})
				w.Write(data)
			}
		}

	}
}
func CatalogEditHandler(w http.ResponseWriter, r *http.Request) {
	user := session.GetSessionData(r)
	if user != nil {
		keys := r.URL.Query()
		id, err := strconv.Atoi(keys.Get("id"))
		data := db_catalog.GetCatalogById(user, id)
		files := []string{
			"./ui/html/catalog/catalog-create.page.tmpl",
			"./ui/html/base.layout.tmpl",
			"./ui/html/top.layout.tmpl",
			"./ui/html/controls/save.tmpl",
		}

		tpl, err := template.ParseFiles(files...)
		if err != nil {
			log.Println(err.Error())
		}

		tpl.Execute(w, data)
	} else {
		data, _ := json.Marshal(route.Message{Text: "not-login"})
		w.Write(data)
	}
}
func GetDbFieldsHandler(w http.ResponseWriter, r *http.Request) {
	user := session.GetSessionData(r)
	if user != nil {
		type S struct {
			Name string `json:"name"`
		}
		decoder := json.NewDecoder(r.Body)
		var param S
		err := decoder.Decode(&param)
		if err != nil {

		}
		list, err := db_catalog.GetDbTableFields(user, param.Name)
		data, _ := json.Marshal(list)
		w.Write(data)
	} else {
		data, _ := json.Marshal(route.Message{Text: "not-login"})
		w.Write(data)
	}
}

func SaveDbFieldsHandler(w http.ResponseWriter, r *http.Request) {
	user := session.GetSessionData(r)
	if user != nil {

		decoder := json.NewDecoder(r.Body)
		var param []model.Field
		err := decoder.Decode(&param)
		if err != nil {

		}
		err = db_catalog.SaveCatalogFields(user, param)
		if err != nil {

		}
		data, _ := json.Marshal("ok")
		w.Write(data)
	} else {
		data, _ := json.Marshal(route.Message{Text: "not-login"})
		w.Write(data)
	}
}
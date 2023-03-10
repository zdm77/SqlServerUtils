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
	keys := r.URL.Query()
	typeId := keys.Get("id")
	tpl, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
	}

	tpl.Execute(w, typeId)
}
func GetCatalogListHandler(w http.ResponseWriter, r *http.Request) {
	user := session.GetSessionData(r)
	if user != nil {
		keys := r.URL.Query()
		typeId := keys.Get("id")
		list := db_catalog.GetCatalogList(user, typeId)
		data, _ := json.Marshal(list)
		w.Write(data)
	} else {
		data, _ := json.Marshal(route.Message{Text: "not-login"})
		w.Write(data)
	}
}
func GetLinkListHandler(w http.ResponseWriter, r *http.Request) {
	user := session.GetSessionData(r)
	if user != nil {
		keys := r.URL.Query()
		id, err := strconv.Atoi(keys.Get("id"))
		err, list := db_catalog.GetLinkList(user, id)
		var data []byte
		if err != nil {
			data, _ = json.Marshal(route.Message{Text: err.Error()})
		} else {
			data, _ = json.Marshal(list)

		}
		w.Write(data)
	} else {
		data, _ := json.Marshal(route.Message{Text: "not-login"})
		w.Write(data)
	}
}
func CatalogCreateHandler(w http.ResponseWriter, r *http.Request) {
	keys := r.URL.Query()
	typeId := keys.Get("typeId")
	files := []string{
		"./ui/html/catalog/catalog-create.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/top.layout.tmpl",
		"./ui/html/controls/save.tmpl",
	}

	catalog := model.Catalog{
		Id:         0,
		EntityId:   0,
		Name:       "",
		TableName:  "",
		Fields:     nil,
		TypeEntity: typeId,
	}
	tpl, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
	}

	tpl.Execute(w, catalog)
	//tpl.Execute(w, nil)
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
		data := db_catalog.GetCatalogById(user, id, false)
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
		list, err := db_catalog.GetDbTableFields(user, param.Name, true)
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
func CatalogListDeleteListHandler(w http.ResponseWriter, r *http.Request) {
	user := session.GetSessionData(r)
	if user != nil {
		keys := r.URL.Query()
		id, err := strconv.Atoi(keys.Get("id"))
		err = db_catalog.DeleteCatalogList(user, id)
		var data []byte
		if err != nil {
			data, _ = json.Marshal(route.Message{Text: err.Error()})
		} else {
			data, _ = json.Marshal(route.Message{Text: "ok"})
		}
		w.Write(data)
	} else {
		data, _ := json.Marshal(route.Message{Text: "not-login"})
		w.Write(data)
	}
}

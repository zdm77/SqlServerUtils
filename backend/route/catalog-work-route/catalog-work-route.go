package catalog_work_route

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	db_catalog_work "sqlutils/backend/database/db-catalog-work"
	"sqlutils/backend/database/db_catalog"
	"sqlutils/backend/model"
	"sqlutils/backend/route"
	"sqlutils/backend/session"
	"strconv"
)

func CatalogWorkListHandler(w http.ResponseWriter, r *http.Request) {
	user := session.GetSessionData(r)
	if user != nil {
		keys := r.URL.Query()
		id, err := strconv.Atoi(keys.Get("id"))
		//	_, data := db_catalog_work.GetCatalogWorkListById(user, id)
		files := []string{
			"./ui/html/catalog-work/catalog-work-list-json.tmpl",
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

		tpl.Execute(w, id)
	} else {
		data, _ := json.Marshal(route.Message{Text: "not-login"})
		w.Write(data)
	}
}
func CatalogGetWorkListHandler(w http.ResponseWriter, r *http.Request) {
	user := session.GetSessionData(r)
	if user != nil {
		decoder := json.NewDecoder(r.Body)
		type Id struct {
			Id int `json:"id"`
		}
		var id Id
		err := decoder.Decode(&id)
		if err != nil {
			log.Println(err.Error())
		}
		//_, list := db_catalog_work.GetCatalogWorkListById(user, id.Id, true)
		_, list := db_catalog_work.GetCatalogWorkListByIdJson(user, id.Id, true)
		data, _ := json.Marshal(list)
		w.Write(data)
	} else {
		data, _ := json.Marshal(route.Message{Text: "not-login"})
		w.Write(data)
	}
}

func CatalogWorkCreateHandler(w http.ResponseWriter, r *http.Request) {
	user := session.GetSessionData(r)
	if user != nil {
		keys := r.URL.Query()
		id, err := strconv.Atoi(keys.Get("id"))
		data := db_catalog.GetCatalogById(user, id, true, false, false)

		files := []string{
			"./ui/html/catalog-work/catalog-work-crate.page.tmpl",
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
func CatalogWorkEditHandler(w http.ResponseWriter, r *http.Request) {
	user := session.GetSessionData(r)
	if user != nil {
		keys := r.URL.Query()
		catalogId, err := strconv.Atoi(keys.Get("id"))
		entityId, err := strconv.Atoi(keys.Get("entityId"))

		err, data := db_catalog_work.GetEntityByCatalogId(user, catalogId, entityId)
		files := []string{
			"./ui/html/catalog-work/catalog-work-crate.page.tmpl",
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
func CatalogWorkSaveHandler(w http.ResponseWriter, r *http.Request) {
	user := session.GetSessionData(r)
	if user != nil {
		decoder := json.NewDecoder(r.Body)
		var param model.Catalog
		err := decoder.Decode(&param)
		if err != nil {
			log.Println(err.Error())
		} else {
			err, id := db_catalog_work.SaveCatalogWork(user, param)
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
func CatalogWorkListDeleteListHandler(w http.ResponseWriter, r *http.Request) {
	user := session.GetSessionData(r)
	if user != nil {
		keys := r.URL.Query()
		id, err := strconv.Atoi(keys.Get("id"))
		catalogId, err := strconv.Atoi(keys.Get("catalogId"))
		err = db_catalog_work.DeleteCatalogWorkList(user, id, catalogId)
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

func CatalogAccessRecordHandler(w http.ResponseWriter, r *http.Request) {
	user := session.GetSessionData(r)
	if user != nil {
		decoder := json.NewDecoder(r.Body)
		type Param struct {
			Id    int    `json:"id"`
			Table string `json:"table"`
		}
		var param Param
		err := decoder.Decode(&param)
		if err != nil {
			log.Println(err.Error())
		}
		_, list := db_catalog_work.GetCatalogAccessRecord(user, param.Id, param.Table)
		data, _ := json.Marshal(list)
		w.Write(data)
	} else {
		data, _ := json.Marshal(route.Message{Text: "not-login"})
		w.Write(data)
	}
}
func SaveAccessRecordHandler(w http.ResponseWriter, r *http.Request) {
	user := session.GetSessionData(r)
	if user != nil {
		decoder := json.NewDecoder(r.Body)
		type Param struct {
			Id     int    `json:"id"`
			Table  string `json:"table"`
			Access string `json:"access"`
		}
		var param Param
		err := decoder.Decode(&param)
		if err != nil {
			log.Println(err.Error())
		}
		err = db_catalog_work.SaveCatalogAccessRecord(user, param.Id, param.Table, param.Access)
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

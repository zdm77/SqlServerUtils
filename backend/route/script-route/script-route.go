package script_route

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

func ScriptAccessHandler(w http.ResponseWriter, r *http.Request) {
	user := session.GetSessionData(r)
	if user != nil {
		decoder := json.NewDecoder(r.Body)
		type Param struct {
			Id int `json:"id"`
		}
		var param Param
		err := decoder.Decode(&param)
		if err != nil {
			log.Println(err.Error())
		}
		_, list := db_catalog.GetAccessScript(user, param.Id)
		data, _ := json.Marshal(list)
		w.Write(data)
	} else {
		data, _ := json.Marshal(route.Message{Text: "not-login"})
		w.Write(data)
	}
}
func ScriptListHandler(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/catalog/script-list.tmpl",
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
func ScriptListWorkHandler(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/script-work/script-work.tmpl",
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
func GetCatalogScriptListHandler(w http.ResponseWriter, r *http.Request) {
	user := session.GetSessionData(r)
	if user != nil {
		list := db_catalog.GetCatalogScript(user)
		data, _ := json.Marshal(list)
		w.Write(data)
	} else {
		data, _ := json.Marshal(route.Message{Text: "not-login"})
		w.Write(data)
	}
}
func ScriptEditHandler(w http.ResponseWriter, r *http.Request) {
	user := session.GetSessionData(r)
	if user != nil {
		keys := r.URL.Query()
		id, err := strconv.Atoi(keys.Get("id"))
		data := db_catalog.GetScriptById(user, id)
		files := []string{
			"./ui/html/catalog/script-create.tmpl",
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
func ScriptSaveHandler(w http.ResponseWriter, r *http.Request) {
	user := session.GetSessionData(r)
	if user != nil {
		decoder := json.NewDecoder(r.Body)
		var param model.Script
		err := decoder.Decode(&param)
		if err != nil {
			log.Println(err.Error())
		} else {
			db_catalog.SaveScript(user, param)
			if err != nil && err.Error() != "sql: no rows in result set" {
				data, _ := json.Marshal(route.Message{Text: err.Error()})
				w.Write(data)
			} else {
				data, _ := json.Marshal(route.Message{Text: "ok"})
				w.Write(data)
			}
		}

	}
}

func ScriptExeHandler(w http.ResponseWriter, r *http.Request) {
	user := session.GetSessionData(r)
	if user != nil {
		keys := r.URL.Query()
		id, _ := strconv.Atoi(keys.Get("id"))

		_, isErr := db_catalog.ExeScript(user, id)
		if isErr {
			data, _ := json.Marshal(route.Message{Text: "Ошибка"})
			w.Write(data)
		} else {
			data, _ := json.Marshal(route.Message{Text: "ok"})
			w.Write(data)
		}

	} else {
		data, _ := json.Marshal(route.Message{Text: "not-login"})
		w.Write(data)
	}
}
func SaveAccessScriptHandler(w http.ResponseWriter, r *http.Request) {
	user := session.GetSessionData(r)
	if user != nil {
		decoder := json.NewDecoder(r.Body)
		type Param struct {
			Id     int    `json:"id"`
			Access string `json:"access"`
		}
		var param Param
		err := decoder.Decode(&param)
		if err != nil {
			log.Println(err.Error())
		}
		err = db_catalog.SaveScriptAccess(user, param.Id, param.Access)
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

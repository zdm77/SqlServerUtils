package task_route

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"sqlutils/backend/model"
	"sqlutils/backend/route"
	"sqlutils/backend/session"
	"sqlutils/backend/task"
)

func Task1Handler(w http.ResponseWriter, r *http.Request) {

	//	if user != nil {
	files := []string{
		"./ui/html/task/task-exe.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/top.layout.tmpl",
	}

	tpl, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
	}
	user := session.GetSessionData(r)
	var data []byte
	if user != nil {
		data, _ = json.Marshal(route.Message{Text: "not-login"})

	}
	tpl.Execute(w, data)
}

func TaskExeHandler(w http.ResponseWriter, r *http.Request) {
	user := session.GetSessionData(r)
	var data []byte
	if user != nil {
		decoder := json.NewDecoder(r.Body)
		var tsk model.Result
		err := decoder.Decode(&tsk)
		if err != nil {
			log.Println(err.Error())
		}
		err = task.ExecTaskFromExcel(user, tsk)
		if err == nil {
			data, _ = json.Marshal(route.Message{Text: "ok"})
			w.Write(data)
		} else {
			data, _ = json.Marshal(route.Message{Text: "Ошибка: " + err.Error()})
			w.Write(data)
		}
		//data, _ = json.Marshal(params)

	} else {
		data, _ = json.Marshal(route.Message{Text: "not-login"})
		w.Write(data)
	}

}

package task_route

import (
	"encoding/json"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	db_task "sqlutils/backend/database/db-task"
	"sqlutils/backend/model"
	"sqlutils/backend/route"
	"sqlutils/backend/session"
	"sqlutils/backend/task"
	"strconv"
	"time"
)

func TaskListCatalogHandler(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/catalog/task-list.page.tmpl",
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
func GetTaskCatalogListHandler(w http.ResponseWriter, r *http.Request) {
	user := session.GetSessionData(r)
	if user != nil {

		tasks := db_task.GetTaskCatalogList(user)
		data, _ := json.Marshal(tasks)
		w.Write(data)
	} else {
		data, _ := json.Marshal(route.Message{Text: "not-login"})
		w.Write(data)
	}
}
func TaskCatalogCreateHandler(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/catalog/task-create.page.tmpl",
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

func TaskSaveParamsHandler(w http.ResponseWriter, r *http.Request) {
	user := session.GetSessionData(r)
	if user != nil {
		decoder := json.NewDecoder(r.Body)
		var task []model.TaskParams
		err := decoder.Decode(&task)
		if err != nil {
			log.Println(err.Error())
		} else {
			err = db_task.SaveTaskCatalogParams(user, task)
			if err != nil {
				data, _ := json.Marshal(route.Message{Text: err.Error()})
				w.Write(data)
			} else {
				data, _ := json.Marshal(route.Message{Text: "ok"})
				w.Write(data)
			}
		}

	}
}
func GetTaskParamsHandler(w http.ResponseWriter, r *http.Request) {
	user := session.GetSessionData(r)
	if user != nil {
		decoder := json.NewDecoder(r.Body)
		var task model.Task
		err := decoder.Decode(&task)
		if err != nil {
			log.Println(err.Error())
		}
		params := db_task.GetTaskCatalogParams(user, task.Id, false)
		data, _ := json.Marshal(params)
		w.Write(data)
	} else {
		data, _ := json.Marshal(route.Message{Text: "not-login"})
		w.Write(data)
	}
}
func TaskUploadHandler(w http.ResponseWriter, r *http.Request) {
	user := session.GetSessionData(r)
	tmpDir := path.Join("tmp", time.Now().Local().Format("02012006-150405"))
	os.MkdirAll("tmp", 0777)
	os.MkdirAll(tmpDir, 0777)
	var taskId int
	var strHeader int
	onlyHeader := true
	//fmt.Println(taskId)
	var fileExec string
	if user != nil {
		m, err := r.MultipartReader()
		if err != nil {
			log.Println(err.Error())
		} else {
			//var files []string
			for {
				part, err := m.NextPart()
				if err == io.EOF {
					break
				}
				if part.FileName() != "" {
					fileExec = filepath.Join(tmpDir, part.FileName())
					dst, err := os.Create(fileExec)
					_, err = io.Copy(dst, part)
					if err != nil {
						log.Println(err.Error())
						break
					}
					err = part.Close()
					err = dst.Close()

				} else {

					val, _ := ioutil.ReadAll(part)
					switch part.FormName() {
					case "task_id":
						{
							val, _ := strconv.Atoi(string(val))
							taskId = val
						}
					case "only_headers":
						{
							val, _ := strconv.ParseBool(string(val))
							onlyHeader = val
						}
					case "str_header":
						{
							val, _ := strconv.Atoi(string(val))
							strHeader = val
						}
					}
				}
			}

		}
		//task.ExecTaskFromExcel(user, fileExec, taskId)
	}
	if onlyHeader {
		params := task.GetHeaders(fileExec, strHeader)
		data, _ := json.Marshal(params)
		w.Write(data)
	} else {
		params := task.GetData(user, fileExec, taskId)
		data, _ := json.Marshal(params)
		w.Write(data)
	}

	//os.RemoveAll(tmpDir)

}
func TaskCatalogSaveHandler(w http.ResponseWriter, r *http.Request) {
	user := session.GetSessionData(r)
	if user != nil {
		decoder := json.NewDecoder(r.Body)
		var task model.Task
		err := decoder.Decode(&task)
		if err != nil {
			log.Println(err.Error())
		} else {
			err, id := db_task.SaveTaskCatalog(user, task)
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
func TaskCatalogEditHandler(w http.ResponseWriter, r *http.Request) {
	user := session.GetSessionData(r)
	if user != nil {
		keys := r.URL.Query()
		id, err := strconv.Atoi(keys.Get("id"))
		task := db_task.GetTaskCatalogById(user, id)
		files := []string{
			"./ui/html/catalog/task-create.page.tmpl",
			"./ui/html/base.layout.tmpl",
			"./ui/html/top.layout.tmpl",
			"./ui/html/controls/save.tmpl",
		}

		tpl, err := template.ParseFiles(files...)
		if err != nil {
			log.Println(err.Error())
		}

		tpl.Execute(w, task)
	} else {
		data, _ := json.Marshal(route.Message{Text: "not-login"})
		w.Write(data)
	}
}

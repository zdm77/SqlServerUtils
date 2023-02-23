package route

import (
	"bufio"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sqlutils/backend/database"
	db_task "sqlutils/backend/database/db-task"
	"sqlutils/backend/model"
	"sqlutils/backend/session"
	"sqlutils/backend/task"
	"strconv"
	"strings"
	"time"
)

var MainPath, _ = os.Getwd()

type Message struct {
	Text string `json:"text"`
}

func Home(w http.ResponseWriter, r *http.Request) {
	var tpl = template.Must(template.ParseFiles("./ui/html/home.page.tmpl"))

	tpl.Execute(w, nil)
}
func Login(w http.ResponseWriter, r *http.Request) {
	var tpl = template.Must(template.ParseFiles("./ui/html/login.page.tmpl"))
	//var tpl = template.Must(template.ParseFiles("./ui/html/test.page.html"))
	data := ""
	tpl.Execute(w, data)
}
func DoLogin(w http.ResponseWriter, r *http.Request) {

	user := session.GetSessionData(r)
	if user != nil && user.IsLogin {
		fmt.Println("")
	}
	login := r.FormValue("login")
	password := r.FormValue("passwd")
	server := ""
	port := ""
	dbName := ""
	file, err := os.Open("settings.txt")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer file.Close()
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		cont := strings.Contains(sc.Text(), "DBHost")
		if cont {
			server = strings.Split(sc.Text(), "=")[1]
		}
		cont = strings.Contains(sc.Text(), "Port")
		if cont {
			port = strings.Split(sc.Text(), "=")[1]
		}
		cont = strings.Contains(sc.Text(), "DataBase")
		if cont {
			dbName = strings.Split(sc.Text(), "=")[1]
		}
	}
	connStr := database.SetParam(server, login, password, port, dbName)
	db, err := database.GetDb(connStr)
	defer db.Close()
	err = db.Ping()

	if err != nil {
		fmt.Println(err.Error())
		var tpl = template.Must(template.ParseFiles("./ui/html/login.page.tmpl"))
		data := "Не верное имя или пароль"
		tpl.Execute(w, data)
	} else {
		var userSession model.User
		userSession.IsLogin = true
		userSession.Login = login
		userSession.ConnString = connStr
		session.Save(userSession, w, r)
		files := []string{
			"./ui/html/home.page.tmpl",
			"./ui/html/base.layout.tmpl",
			"./ui/html/top.layout.tmpl",
		}
		tpl, err := template.ParseFiles(files...)
		if err != nil {
			log.Println(err.Error())
		}

		tpl.Execute(w, nil)
	}
}
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
		data, _ = json.Marshal(Message{Text: "not-login"})

	}
	tpl.Execute(w, data)
}
func ReturnToLogin(w http.ResponseWriter) {
	var tpl = template.Must(template.ParseFiles("./ui/html/login.page.tmpl"))
	data := ""
	tpl.Execute(w, data)
}
func GetTaskListHandler(w http.ResponseWriter, r *http.Request) {
	user := session.GetSessionData(r)
	if user != nil {

		tasks := db_task.GetTaskList(user)
		data, _ := json.Marshal(tasks)
		w.Write(data)
	} else {
		data, _ := json.Marshal(Message{Text: "not-login"})
		w.Write(data)
	}
}
func TaskListHandler(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/catalog/task-list.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/top.layout.tmpl",
	}

	tpl, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
	}

	tpl.Execute(w, nil)
}
func TaskCreateHandler(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/catalog/task-create.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/top.layout.tmpl",
	}

	tpl, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
	}

	tpl.Execute(w, nil)
}
func TaskEditHandler(w http.ResponseWriter, r *http.Request) {
	user := session.GetSessionData(r)
	if user != nil {
		keys := r.URL.Query()
		id, err := strconv.Atoi(keys.Get("id"))
		task := db_task.GetTaskById(user, id)
		files := []string{
			"./ui/html/catalog/task-create.page.tmpl",
			"./ui/html/base.layout.tmpl",
			"./ui/html/top.layout.tmpl",
		}

		tpl, err := template.ParseFiles(files...)
		if err != nil {
			log.Println(err.Error())
		}
		//		data, _ := json.Marshal(task)

		tpl.Execute(w, task)
	} else {
		data, _ := json.Marshal(Message{Text: "not-login"})
		w.Write(data)
	}
}
func TaskSaveHandler(w http.ResponseWriter, r *http.Request) {
	user := session.GetSessionData(r)
	if user != nil {
		decoder := json.NewDecoder(r.Body)
		var task model.Task
		err := decoder.Decode(&task)
		if err != nil {
			log.Println(err.Error())
		} else {
			err, id := db_task.SaveTask(user, task)
			if err != nil && err.Error() != "sql: no rows in result set" {
				data, _ := json.Marshal(Message{Text: err.Error()})
				w.Write(data)
			} else {
				data, _ := json.Marshal(Message{Text: "ok-" + strconv.Itoa(id)})
				w.Write(data)
			}
		}

	}
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
			err = db_task.SaveTaskParams(user, task)
			if err != nil {
				data, _ := json.Marshal(Message{Text: err.Error()})
				w.Write(data)
			} else {
				data, _ := json.Marshal(Message{Text: "ok"})
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
		params := db_task.GetTaskParams(user, task.Id, false)
		data, _ := json.Marshal(params)
		w.Write(data)
	} else {
		data, _ := json.Marshal(Message{Text: "not-login"})
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
		params := task.GetHeaders(user, fileExec, strHeader)
		data, _ := json.Marshal(params)
		w.Write(data)
	} else {
		params := task.GetData(user, fileExec, taskId)
		data, _ := json.Marshal(params)
		w.Write(data)
	}

	//os.RemoveAll(tmpDir)

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
			data, _ = json.Marshal(Message{Text: "ok"})
			w.Write(data)
		} else {
			data, _ = json.Marshal(Message{Text: "Ошибка: " + err.Error()})
			w.Write(data)
		}
		//data, _ = json.Marshal(params)

	} else {
		data, _ = json.Marshal(Message{Text: "not-login"})
		w.Write(data)
	}

}

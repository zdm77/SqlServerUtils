package route

import (
	"bufio"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"sqlutils/backend/database"
	"sqlutils/backend/model"
	"sqlutils/backend/session"
	"strings"
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
		cont = strings.Contains(sc.Text(), "DBPort")
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

func ReturnToLogin(w http.ResponseWriter) {
	var tpl = template.Must(template.ParseFiles("./ui/html/login.page.tmpl"))
	data := ""
	tpl.Execute(w, data)
}

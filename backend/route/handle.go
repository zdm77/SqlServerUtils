package route

import (
	"bufio"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"sqlutils/backend/database"
	"strings"
)

var MainPath, _ = os.Getwd()

func Home(w http.ResponseWriter, r *http.Request) {
	var tpl = template.Must(template.ParseFiles("./ui/html/home.page.tmpl"))

	tpl.Execute(w, nil)
}
func Login(w http.ResponseWriter, r *http.Request) {
	var tpl = template.Must(template.ParseFiles("./ui/html/login.page.tmpl"))
	data := ""
	tpl.Execute(w, data)
}
func DoLogin(w http.ResponseWriter, r *http.Request) {
	user := r.FormValue("login")
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
	database.SetParam(server, user, password, port, dbName)
	db, err := database.GetDb()
	defer db.Close()
	err = db.Ping()
	type MessageObject struct {
		Id int
	}
	if err != nil {
		fmt.Println(err.Error())
		var tpl = template.Must(template.ParseFiles("./ui/html/login.page.tmpl"))
		data := "Не верное имя или пароль"
		tpl.Execute(w, data)
	} else {
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
func Task1(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/task1.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/top.layout.tmpl",
	}

	tpl, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
	}

	tpl.Execute(w, nil)
}

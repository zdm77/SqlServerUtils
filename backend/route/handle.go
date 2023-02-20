package route

import (
	"html/template"
	"log"
	"net/http"
	"os"
)

var MainPath, _ = os.Getwd()

func Home(w http.ResponseWriter, r *http.Request) {

	os.Chdir(MainPath)
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	ts, err := template.ParseFiles("./ui/html/home.page.tmpl")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
	err = ts.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}

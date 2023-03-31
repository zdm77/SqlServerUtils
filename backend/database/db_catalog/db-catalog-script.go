package db_catalog

import (
	"os"
	"path/filepath"
	"sqlutils/backend/database"
	"sqlutils/backend/model"
	"strconv"
)

func GetCatalogScript(user *model.User) (result []model.Script) {
	dir := user.ScriptCatalog
	type File struct {
		Dir   string
		Name  string
		PathF string
	}
	var files []File
	filepath.Walk(dir,
		func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				files = append(files, File{
					Dir:   filepath.Dir(path),
					Name:  filepath.Base(path),
					PathF: path,
				})
			}
			return nil
		})
	db, _ := database.GetDb(user.ConnString)
	defer db.Close()
	var query string
	for _, file := range files {
		var id int
		query = `select id from utils_script where script_name='` + file.PathF + `'`
		db.QueryRow(query).Scan(&id)

		if id == 0 {
			query = `insert into utils_script (name, script_name) values ('` + file.Name[:len(file.Name)-4] + `','` + file.PathF + `')`
			db.Exec(query)
		}
	}
	query = `select id, name, script_name from utils_script order by script_name`
	rows, _ := db.Query(query)
	for rows.Next() {
		var r model.Script
		rows.Scan(&r.Id, &r.Name, &r.ScriptName)
		result = append(result, r)
	}

	return result
}
func GetScriptById(user *model.User, id int) (r model.Script) {
	db, _ := database.GetDb(user.ConnString)
	defer db.Close()
	query := `select id, name, script_name from utils_script where id =` + strconv.Itoa(id)
	db.QueryRow(query).Scan(&r.Id, &r.Name, &r.ScriptName)
	return r
}
func SaveScript(user *model.User, script model.Script) () {
	db, _ := database.GetDb(user.ConnString)
	defer db.Close()
	query := `update utils_script set name ='` + script.Name + `' where id = ` + strconv.Itoa(script.Id)
	db.Exec(query)
}

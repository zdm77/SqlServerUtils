package db_catalog

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sqlutils/backend/database"
	"sqlutils/backend/model"
	"strconv"
	"strings"
)

func GetCatalogScript(user *model.User) (result []model.Script) {
	dir := user.ScriptCatalog
	type File struct {
		Dir   string
		Name  string
		PathF string
	}
	var files []File
	dirO, _ := os.ReadDir(dir)
	for _, d := range dirO {
		if d.IsDir() {
			dr, _ := os.ReadDir(filepath.Join(dir, d.Name()))
			for _, dd := range dr {
				if !dd.IsDir() {
					files = append(files, File{
						Dir:   filepath.Join(dir, d.Name()),
						Name:  dd.Name(),
						PathF: filepath.Join(dir, d.Name(), dd.Name()),
					})
				}
			}
		}
	}
	//filepath.Walk(dir,
	//	func(path string, info os.FileInfo, err error) error {
	//		if !info.IsDir() {
	//			files = append(files, File{
	//				Dir:   filepath.Dir(path),
	//				Name:  filepath.Base(path),
	//				PathF: path,
	//			})
	//		}
	//		return nil
	//	})
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
func ExeScript(user *model.User, scriptId int) (err error, isErr bool) {
	db, _ := database.GetDb(user.ConnString)
	defer db.Close()
	query := `select  script_name  from utils_script  where id = ` + strconv.Itoa(scriptId)
	var script string
	err = db.QueryRow(query).Scan(&script)
	if err != nil {
		log.Println(err.Error())
		return err, true
	}
	exe := user.PythonExe
	cmd := exec.Command(exe, script)
	out := bytes.Buffer{}
	cmd.Stderr = &out
	err = cmd.Run()

	isErr = strings.Contains(strings.ToUpper(out.String()), "ERROR")
	if isErr {
		fmt.Println(out.String())
	}

	if isErr {

		log.Println("Буфер:" + out.String())
		err = fmt.Errorf("Ошибка выполнения скрипта")
		return err, isErr
	}
	return err, isErr
}

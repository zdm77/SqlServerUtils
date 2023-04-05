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

func DeleteNotPyScript(conn string) {
	db, err := database.GetDb(conn)
	if err != nil {
		log.Println(err.Error())
	}
	defer db.Close()
	query := `delete from utils_script where script_name not like '%py'`
	_, err = db.Exec(query)
	if err != nil {
		log.Println(err.Error())
	}
}
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
					arr := strings.Split(dd.Name(), ".")
					if arr[1] == "py" {
						files = append(files, File{
							Dir:   filepath.Join(dir, d.Name()),
							Name:  dd.Name(),
							PathF: filepath.Join(dir, d.Name(), dd.Name()),
						})
					}

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
	query = `select id, name, script_name, access from utils_script order by script_name`
	rows, _ := db.Query(query)

	for rows.Next() {
		var r model.Script
		isAccess := false
		rows.Scan(&r.Id, &r.Name, &r.ScriptName, &r.Access)
		if user.SuperAdmin != user.Login {
			arr := strings.Split(r.Access, ",")
			for _, a := range arr {
				if a == user.Login {
					isAccess = true
					break
				}
			}
		} else {
			isAccess = true
		}
		if isAccess {
			result = append(result, r)
		}
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
func SaveScript(user *model.User, script model.Script) {
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
func GetAccessScript(user *model.User, id int) (err error, result []model.AccessRecord) {
	db, _ := database.GetDb(user.ConnString)
	defer db.Close()
	query := `select access from utils_script where id=` + strconv.Itoa(id)

	var usersInTableStr string
	err = db.QueryRow(query).Scan(&usersInTableStr)
	usersInTable := strings.Split(usersInTableStr, ",")
	query = `select name from master.sys.server_principals where type_desc ='SQL_LOGIN'  and is_disabled=0 order by name`
	rows, err := db.Query(query)
	if err != nil {
		log.Println(err.Error())
	}
	defer rows.Close()
	var usersAll []string
	for rows.Next() {
		var n string
		err = rows.Scan(&n)
		if err != nil {
			log.Println(err.Error())
		}
		usersAll = append(usersAll, n)
	}
	for _, user := range usersAll {
		isAccess := false
		for _, u := range usersInTable {
			if u == user {
				isAccess = true
				break
			}
		}
		result = append(result, model.AccessRecord{
			UserName: user,
			Access:   isAccess,
		})
	}

	return err, result
}
func SaveScriptAccess(user *model.User, id int, access string) (err error) {
	db, _ := database.GetDb(user.ConnString)
	defer db.Close()

	query := `update utils_script set access = '` + access + `'  where id = ` + strconv.Itoa(id)
	_, err = db.Exec(query)
	if err != nil {
		log.Println(err.Error())
	}
	return err
}

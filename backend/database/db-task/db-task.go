package db_task

import (
	"database/sql"
	"log"
	"sqlutils/backend/database"
	"sqlutils/backend/model"
)

func GetTaskList(user *model.User) (result []model.Task) {
	db, _ := database.GetDb(user.ConnString)
	defer db.Close()
	query := `select  id, name, table_db from utils_task order by name`
	rows, err := db.Query(query)
	defer rows.Close()
	if err != nil {
		log.Println(err.Error())
	}
	for rows.Next() {
		var r model.Task
		err = rows.Scan(&r.Id, &r.Name, &r.TableDb)
		if err != nil {

		}
		result = append(result, r)

	}
	return result
}
func GetTaskById(user *model.User, id int) (r model.Task) {
	db, _ := database.GetDb(user.ConnString)
	defer db.Close()
	query := `select  id, name, table_db from utils_task where id = @Id`

	stmt, err := db.Prepare(query)
	row := stmt.QueryRow(sql.Named("Id", id))
	err = row.Scan(&r.Id, &r.Name, &r.TableDb)
	if err != nil {
		log.Println(err.Error())
	}

	return r
}
func SaveTask(user *model.User, task model.Task) (err error) {
	db, _ := database.GetDb(user.ConnString)
	defer db.Close()
	var query string
	var stmt *sql.Stmt
	db.Query(query)
	if task.Id == 0 {
		query = `insert into utils_task (name, table_db) values (@Name, @TableDb)`
	} else {
		query = `update  utils_task set name = @Name, table_db = @TableDb where id = @Id`
	}
	stmt, _ = db.Prepare(query)
	_, err = stmt.Exec(sql.Named("Name", task.Name),
		sql.Named("TableDb", task.TableDb),
		sql.Named("Id", task.Id))
	if err != nil {
		log.Print(err.Error())
	}

	return err
}

package db_task

import (
	"database/sql"
	"log"
	"sqlutils/backend/database"
	"sqlutils/backend/model"
	"strconv"
)

func GetTaskList(user *model.User) (result []model.Task) {
	db, _ := database.GetDb(user.ConnString)
	defer db.Close()
	query := `select  id, name, table_db,  str_header from utils_task order by id`
	rows, err := db.Query(query)
	defer rows.Close()
	if err != nil {
		log.Println(err.Error())
	}
	for rows.Next() {
		var r model.Task
		err = rows.Scan(&r.Id, &r.Name, &r.TableDb, &r.StrHeader)
		if err != nil {

		}
		result = append(result, r)

	}
	return result
}
func GetTaskById(user *model.User, id int) (r model.Task) {
	db, _ := database.GetDb(user.ConnString)
	defer db.Close()
	query := `select  id, name, table_db, coalesce(str_header, 1) from utils_task where id = @Id`

	stmt, err := db.Prepare(query)
	row := stmt.QueryRow(sql.Named("Id", id))
	err = row.Scan(&r.Id, &r.Name, &r.TableDb, &r.StrHeader)
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

	if task.Id == 0 {
		query = `insert into utils_task (name, table_db, str_header) values (@Name, @TableDb, @StrHeader)`
	} else {
		query = `update  utils_task set name = @Name, table_db = @TableDb, str_header = @StrHeader where id = @Id`
	}
	stmt, _ = db.Prepare(query)
	_, err = stmt.Exec(sql.Named("Name", task.Name),
		sql.Named("TableDb", task.TableDb),
		sql.Named("Id", task.Id),
		sql.Named("StrHeader", task.StrHeader),
	)
	if err != nil {
		log.Print(err.Error())
	}

	return err
}

func SaveTaskParams(user *model.User, params []model.TaskParams) (err error) {
	db, _ := database.GetDb(user.ConnString)
	defer db.Close()
	var query string
	var stmt *sql.Stmt
	query = `delete from settings_utils_task where task_id=` + strconv.Itoa(params[0].TaskId)
	db.Exec(query)
	query = `insert into settings_utils_task (task_id, field_excel, field_db) values (@task_id, @field_excel, @field_db)`

	for _, param := range params {
		stmt, _ = db.Prepare(query)
		_, err = stmt.Exec(sql.Named("task_id", param.TaskId),
			sql.Named("field_excel", param.FieldExcel),
			sql.Named("field_db", param.FieldDb),
		)
		if err != nil {
			log.Print(err.Error())
			return err
		}
	}

	return err
}

func GetTaskParams(user *model.User, id int, isValue bool) (result []model.TaskParams) {
	db, _ := database.GetDb(user.ConnString)
	defer db.Close()
	query := `select  task_id, field_excel, field_db from settings_utils_task where task_id=` + strconv.Itoa(id)
	if isValue {
		query += ` and field_db!='' and field_db is not null`
	}
	rows, err := db.Query(query)
	defer rows.Close()
	if err != nil {
		log.Println(err.Error())
	}
	for rows.Next() {
		var r model.TaskParams
		err = rows.Scan(&r.TaskId, &r.FieldExcel, &r.FieldDb)
		if err != nil {

		}
		result = append(result, r)

	}
	return result
}

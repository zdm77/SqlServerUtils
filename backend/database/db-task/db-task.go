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
	query := `select  id, name, table_db,  str_header from utils_task_catalog order by id`
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
	query := `select  id, name, table_db, coalesce(str_header, 1) from utils_task_catalog where id = @Id`

	stmt, err := db.Prepare(query)
	row := stmt.QueryRow(sql.Named("Id", id))
	err = row.Scan(&r.Id, &r.Name, &r.TableDb, &r.StrHeader)
	if err != nil {
		log.Println(err.Error())
	}

	return r
}
func SaveTask(user *model.User, task model.Task) (err error, id int) {
	db, _ := database.GetDb(user.ConnString)
	defer db.Close()
	var query string
	var stmt *sql.Stmt

	if task.Id == 0 {
		query = `insert into utils_task_catalog (name, table_db, str_header)
					values (@Name, @TableDb, @StrHeader);  SELECT SCOPE_IDENTITY()`
	} else {
		query = `update  utils_task_catalog set name = @Name, table_db = @TableDb, str_header = @StrHeader where id = @Id`
	}
	stmt, _ = db.Prepare(query)
	err = stmt.QueryRow(sql.Named("Name", task.Name),
		sql.Named("TableDb", task.TableDb),
		sql.Named("Id", task.Id),
		sql.Named("StrHeader", task.StrHeader),
	).Scan(&id)

	if err != nil && err.Error() != "sql: no rows in result set" {
		log.Print(err.Error())
	}

	return err, id
}

func SaveTaskParams(user *model.User, params []model.TaskParams) (err error) {
	db, _ := database.GetDb(user.ConnString)
	defer db.Close()
	var query string
	var stmt *sql.Stmt
	query = `delete from utils_settings_task where task_id=` + strconv.Itoa(params[0].TaskId)
	db.Exec(query)
	query = `insert into utils_settings_task (task_id, field_excel, field_db, col_number, field_type)
values (@task_id, @field_excel, @field_db, @col_number, @field_type)`

	for _, param := range params {
		stmt, _ = db.Prepare(query)
		_, err = stmt.Exec(sql.Named("task_id", param.TaskId),
			sql.Named("field_excel", param.FieldExcel),
			sql.Named("field_db", param.FieldDb),
			sql.Named("col_number", param.Id),
			sql.Named("field_type", param.FieldType),
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
	query := `select   task_id, field_excel, field_db, col_number, field_type from utils_settings_task where task_id=` + strconv.Itoa(id)
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
		err = rows.Scan(&r.TaskId, &r.FieldExcel, &r.FieldDb, &r.Id, &r.FieldType)
		if err != nil {

		}
		result = append(result, r)

	}
	return result
}

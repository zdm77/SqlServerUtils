package db_task

import (
	"database/sql"
	"log"
	"sqlutils/backend/database"
	"sqlutils/backend/model"
	"strconv"
)

func GetTaskCatalogList(user *model.User) (result []model.Task) {
	db, _ := database.GetDb(user.ConnString)
	defer db.Close()
	query := `select  id, name, table_db,  str_header, sheet_number from utils_task_catalog order by id`
	rows, err := db.Query(query)
	defer rows.Close()
	if err != nil {
		log.Println(err.Error())
	}
	for rows.Next() {
		var r model.Task
		err = rows.Scan(&r.Id, &r.Name, &r.TableDb, &r.StrHeader, &r.SheetNumber)
		if err != nil {

		}
		result = append(result, r)

	}
	return result
}
func GetTaskCatalogById(user *model.User, id int) (r model.Task) {
	db, _ := database.GetDb(user.ConnString)
	defer db.Close()
	query := `select  id, name, table_db, coalesce(str_header, 1) , coalesce(sheet_number, 1)  from utils_task_catalog where id = @Id`

	stmt, err := db.Prepare(query)
	row := stmt.QueryRow(sql.Named("Id", id))
	err = row.Scan(&r.Id, &r.Name, &r.TableDb, &r.StrHeader, &r.SheetNumber)
	if err != nil {
		log.Println(err.Error())
	}

	return r
}
func SaveTaskCatalog(user *model.User, task model.Task) (err error, id int) {
	db, _ := database.GetDb(user.ConnString)
	defer db.Close()
	var query string
	var stmt *sql.Stmt

	if task.Id == 0 {
		query = `insert into utils_task_catalog (name, table_db, str_header, sheet_number)
					values (@Name, @TableDb, @StrHeader, @SheetNumber );  SELECT SCOPE_IDENTITY()`
	} else {
		query = `update  utils_task_catalog set name = @Name, table_db = @TableDb, str_header = @StrHeader,  sheet_number = @SheetNumber where id = @Id`
	}
	stmt, _ = db.Prepare(query)
	err = stmt.QueryRow(sql.Named("Name", task.Name),
		sql.Named("TableDb", task.TableDb),
		sql.Named("Id", task.Id),
		sql.Named("StrHeader", task.StrHeader),
		sql.Named("SheetNumber", task.SheetNumber),
	).Scan(&id)

	if err != nil && err.Error() != "sql: no rows in result set" {
		log.Print(err.Error())
	}

	return err, id
}

func SaveTaskCatalogParams(user *model.User, params []model.TaskParams) (err error) {
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

func GetTaskCatalogParams(user *model.User, id int, isValue bool) (result []model.TaskParams) {
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

func DeleteTaskCatalogList(user *model.User, id int) (err error) {
	db, _ := database.GetDb(user.ConnString)
	defer db.Close()
	query := `delete from utils_task_catalog where id = ` + strconv.Itoa(id)
	_, err = db.Exec(query)

	return err
}

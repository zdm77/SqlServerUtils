package db_catalog

import (
	"database/sql"
	"fmt"
	"log"
	"sqlutils/backend/database"
	"sqlutils/backend/model"
	"strconv"
	"strings"
)

func GetDbTableFields(user *model.User, tableName string) (fields []model.Field, err error) {
	query := `SELECT 
    c.name ,
    t.Name ,
    c.max_length,
    c.precision ,
    c.scale ,
    c.is_nullable,
     c.is_identity,
    ISNULL(i.is_primary_key, 0) ,
     c.is_nullable
FROM    
    sys.columns c
INNER JOIN 
    sys.types t ON c.user_type_id = t.user_type_id
LEFT OUTER JOIN 
    sys.index_columns ic ON ic.object_id = c.object_id AND ic.column_id = c.column_id
LEFT OUTER JOIN 
    sys.indexes i ON ic.object_id = i.object_id AND ic.index_id = i.index_id
WHERE
    c.object_id = OBJECT_ID('` + tableName + `')`
	db, err := database.GetDb(user.ConnString)
	defer db.Close()
	rows, err := db.Query(query)
	defer rows.Close()
	for rows.Next() {
		var f model.Field
		err = rows.Scan(&f.NameDb, &f.NameType, &f.MaxLength, &f.Precision, &f.Scale, &f.IsNullable, &f.IsIdentity, &f.IsPrimaryKey, &f.IsNullableDb)
		fields = append(fields, f)
	}

	return fields, err
}
func GetCatalogList(user *model.User) (result []model.Catalog) {
	db, _ := database.GetDb(user.ConnString)
	defer db.Close()
	query := `select  id, name, table_name from utils_catalog_list order by id`
	rows, err := db.Query(query)
	defer rows.Close()
	if err != nil {
		log.Println(err.Error())
	}
	for rows.Next() {
		var r model.Catalog
		err = rows.Scan(&r.Id, &r.Name, &r.TableName)
		if err != nil {

		}
		result = append(result, r)

	}
	return result
}
func SaveCatalog(user *model.User, param model.Catalog) (err error, id int) {
	db, _ := database.GetDb(user.ConnString)
	defer db.Close()
	var query string
	var stmt *sql.Stmt

	if param.Id == 0 {
		query = `insert into utils_catalog_list (name, table_name)
					values (@Name, @TableDb);  SELECT SCOPE_IDENTITY()`
	} else {
		query = `update  utils_catalog_list set name = @Name, table_name = @TableDb where id = @Id`
	}
	stmt, _ = db.Prepare(query)
	err = stmt.QueryRow(sql.Named("Name", param.Name),
		sql.Named("TableDb", param.TableName),
		sql.Named("Id", param.Id),
	).Scan(&id)

	if err != nil && err.Error() != "sql: no rows in result set" {
		log.Print(err.Error())
	}

	return err, id
}
func GetCatalogById(user *model.User, id int) (r model.Catalog) {
	db, _ := database.GetDb(user.ConnString)
	defer db.Close()
	//основа
	query := `select  id, name, table_name from utils_catalog_list where id = @Id`

	stmt, err := db.Prepare(query)
	row := stmt.QueryRow(sql.Named("Id", id))

	err = row.Scan(&r.Id, &r.Name, &r.TableName)
	if err != nil {
		log.Println(err.Error())
	}
	fieldsDb, err := GetDbTableFields(user, r.TableName)
	//табличная часть
	query = `select id, name, catalog_id, name_db, name_type, max_length, precision, scale, is_nullable, is_identity, is_primary_key, is_nullable_db 
			from utills_catalog_fields where catalog_id=` + strconv.Itoa(id)
	rows, err := db.Query(query)
	defer rows.Close()
	var names []string
	for rows.Next() {
		var f model.Field
		err = rows.Scan(&f.Id, &f.Name, &f.CatalogId, &f.NameDb, &f.NameType, &f.MaxLength, &f.Precision, &f.Scale, &f.IsNullable, &f.IsIdentity, &f.IsPrimaryKey, &f.IsNullableDb)
		name := f.Name
		isNullDb := f.IsNullableDb
		catalogId := f.CatalogId
		entityId := f.Id
		for _, field := range fieldsDb {
			if field.NameDb == f.NameDb {
				f = field
				f.Name = name
				f.IsNullableDb = isNullDb
				f.CatalogId = catalogId
				f.Id = entityId
			}
		}
		r.Fields = append(r.Fields, f)
		names = append(names, f.NameDb)
	}
	//если появились новые поля

	for _, field := range fieldsDb {
		isNew := implContains(names, field.NameDb)
		if isNew == -1 {
			var newId int
			var stmt *sql.Stmt
			query = `insert into utills_catalog_fields (name, catalog_id, name_db, name_type, max_length, precision, scale, is_nullable, is_identity, is_primary_key, is_nullable_db) values
					(@name, @catalog_id, @name_db, @name_type, @max_length, @precision, @scale, @is_nullable, @is_identity, @is_primary_key, @is_nullable_db);  SELECT SCOPE_IDENTITY()`
			stmt, _ = db.Prepare(query)
			err = stmt.QueryRow(
				sql.Named("name", field.Name),
				sql.Named("catalog_id", id),
				sql.Named("name_db", field.NameDb),
				sql.Named("name_type", field.NameType),
				sql.Named("max_length", field.MaxLength),
				sql.Named("precision", field.Precision),
				sql.Named("scale", field.Scale),
				sql.Named("is_nullable", field.IsNullable),
				sql.Named("is_identity", field.IsIdentity),
				sql.Named("is_primary_key", field.IsPrimaryKey),
				sql.Named("is_nullable_db", field.IsNullableDb),
			).Scan(&newId)
			field.Id = newId
			r.Fields = append(r.Fields, field)
			if err != nil {
				log.Println(err.Error())
			}
		}
	}

	return r
}

func implContains(sl []string, name string) int {

	for index, value := range sl {
		if value == name {

			return index
		}
	}

	return -1
}
func SaveCatalogFields(user *model.User, fields []model.Field) (err error) {
	db, _ := database.GetDb(user.ConnString)
	//tx, _ := db.Begin()
	defer db.Close()
	var query string
	var catalogId = fields[0].CatalogId
	var isNew int
	query = `select count(id) from utills_catalog_fields where catalog_id=` + strconv.Itoa(catalogId)
	err = db.QueryRow(query).Scan(&isNew)
	if isNew == 0 {
		query = `insert into utills_catalog_fields (name, catalog_id, name_db, name_type, max_length, precision, scale, is_nullable, is_identity, is_primary_key, is_nullable_db) values
					(@name, @catalog_id, @name_db, @name_type, @max_length, @precision, @scale, @is_nullable, @is_identity, @is_primary_key, @is_nullable_db)`
	} else {
		query = `update utills_catalog_fields set name=@name, catalog_id=@catalog_id, name_db= @name_db, name_type = @name_type,
                                 max_length = @max_length, precision = @precision, scale = @scale, is_nullable = @is_nullable , 
                                 is_identity = @is_identity,  is_nullable_db=@is_nullable_db, is_primary_key = @is_primary_key where id=@id`
	}

	var stmt *sql.Stmt
	for _, field := range fields {
		stmt, _ = db.Prepare(query)
		_, err = stmt.Exec(
			sql.Named("id", field.Id),
			sql.Named("name", field.Name),
			sql.Named("catalog_id", catalogId),
			sql.Named("name_db", field.NameDb),
			sql.Named("name_type", field.NameType),
			sql.Named("max_length", field.MaxLength),
			sql.Named("precision", field.Precision),
			sql.Named("scale", field.Scale),
			sql.Named("is_nullable", field.IsNullable),
			sql.Named("is_identity", field.IsIdentity),
			sql.Named("is_primary_key", field.IsPrimaryKey),
			sql.Named("is_nullable_db", field.IsNullableDb),
		)
		if err != nil {
			log.Println(err.Error())
		}
	}
	return err
}

//редактирование справочника  - записи
func GetEntityByCatalogId(user *model.User, catalogId, entityId int) (err error, data model.Catalog) {
	db, _ := database.GetDb(user.ConnString)
	defer db.Close()
	catalog := GetCatalogById(user, catalogId)

	var fields []string
	query := `select `
	for _, field := range catalog.Fields {
		field.NameDb = "coalesce(cast(" + field.NameDb + " as varchar(255)), '') " + field.NameDb
		fields = append(fields, field.NameDb)
	}
	query += strings.Join(fields, ",") + ` from ` + catalog.TableName + ` where id = ` + strconv.Itoa(entityId) + ` for json auto `
	var jsonStr string
	err = db.QueryRow(query).Scan(&jsonStr)
	if err != nil {
		log.Println(err.Error())
		return err, data
	}
	jsonStr = strings.ReplaceAll(jsonStr, "[{", "")
	jsonStr = strings.ReplaceAll(jsonStr, "}]", "")
	jsonStr = strings.ReplaceAll(jsonStr, "\"", "")
	jsonArr := strings.Split(jsonStr, ",")
	for _, js := range jsonArr {
		arr := strings.Split(js, ":")
		for idx, field := range catalog.Fields {

			if field.NameDb == arr[0] {
				if field.NameType == "bit" {
					fmt.Println("")
				}
				if field.NameType == "bit" && arr[1] == "1" {
					catalog.Fields[idx].Value = "checked"
				} else {
					catalog.Fields[idx].Value = arr[1]
				}
			}
		}

	}
	catalog.EntityId = entityId
	return err, catalog
}

package db_catalog_work

import (
	"log"
	"sqlutils/backend/database"
	"sqlutils/backend/database/db_catalog"
	"sqlutils/backend/model"
	"strconv"
	"strings"
)

type FieldVals struct {
	CatalogId   int        `json:"catalog_id"`
	FieldId     string     `json:"field_id"`
	NameCatalog string     `json:"name_catalog"`
	Headers     []string   `json:"headers"`
	Fields      []string   `json:"fields"`
	Vals        [][]string `json:"vals"`
	ValuesId    []string   `json:"values_id"`
	IsList      []bool     `json:"is_list"`
}

func GetTableLinkById(user *model.User, id int) (err error, table string) {
	db, _ := database.GetDb(user.ConnString)
	defer db.Close()
	query := `select table_name from utils_catalog_list where id=` + strconv.Itoa(id)
	err = db.QueryRow(query).Scan(&table)
	return err, table
}
func GetCatalogWorkListById(user *model.User, id int) (err error, result FieldVals) {

	db, _ := database.GetDb(user.ConnString)
	defer db.Close()

	catalog := db_catalog.GetCatalogById(user, id, false)
	var data [][]string
	var headers []string
	var isList []bool
	query := `select `
	var fields []string
	var fieldId string
	var valuesId []string
	//var catalogRes model.Catalog
	var joinTables string
	queryAccess := `select id, access from ` + catalog.TableName
	type Access struct {
		id    string
		user  []string
		users string
	}
	//var access []Access
	var ids []string
	rows, err := db.Query(queryAccess)
	isAccessInTable := false
	if err == nil {
		isAccessInTable = true
		for rows.Next() {
			var a Access
			rows.Scan(&a.id, &a.users)
			users := strings.Split(a.users, ",")
			for _, u := range users {
				u = strings.TrimSpace(u)
				if u == user.Login {
					ids = append(ids, a.id)
				}
			}
		}
	}
	for _, cat := range catalog.Fields {

		var field string
		if cat.IsPrimaryKey {
			fieldId = cat.NameDb
		}
		if cat.IsIdentity || cat.Name != "" || !cat.IsNullable || !cat.IsNullableDb {
			if cat.LinkTableId != 0 {
				_, table := GetTableLinkById(user, cat.LinkTableId)
				joinTables += " inner join " + table + " on " + catalog.TableName + "." + cat.NameDb + "=" + table + ".id"
				field = table + ".name as " + cat.Name
			} else {
				field = "coalesce(cast(" + catalog.TableName + "." + cat.NameDb + " as varchar(255)), '') " + cat.NameDb
			}

			fields = append(fields, field)
			isList = append(isList, cat.IsList)
			if !cat.IsIdentity {
				headers = append(headers, cat.Name)
			} else {
				headers = append(headers, cat.NameDb)
			}

		}
	}
	query += strings.Join(fields, ",") + " from " + catalog.TableName
	query += joinTables
	if isAccessInTable {
		query += ` where ` + catalog.TableName + `.id in (` + strings.Join(ids, ",") + `) `
	}
	query += "  for json auto "
	var jsonString string
	err = db.QueryRow(query).Scan(&jsonString)
	jsonArr := strings.Split(jsonString, "{")
	//header := make(map[string]bool)
	for idx, js := range jsonArr {
		if idx != 0 {
			str := strings.ReplaceAll(js, "},", "")
			arr := strings.Split(str, ",")
			var dt []string
			for _, a := range arr {
				vl := strings.Split(a, ":")
				nameDb := vl[0]

				val := strings.ReplaceAll(vl[1], `}]`, "")
				for _, f := range catalog.Fields {
					if f.NameDb == strings.ReplaceAll(nameDb, "\"", "") && f.NameType == "bit" {
						if val == `"1"` {
							val = "checked"
						} else {
							val = "unchecked"
						}
					}
				}
				//if catalog.Fields[idxa].NameType == "bit" {
				//	if val == `"1"` {
				//		val = "checked"
				//	} else {
				//		val = "unchecked"
				//	}
				//}
				if nameDb == `"`+fieldId+`"` {
					valuesId = append(valuesId, val)
				}
				dt = append(dt, val[1:len(val)-1])
				//	header[strings.ReplaceAll(vl[0], `"`, "")] = true

			}
			data = append(data, dt)
		}

	}
	//for key, _ := range header {
	//	headers = append(headers, key)
	//}

	result = FieldVals{
		CatalogId:   id,
		FieldId:     fieldId,
		NameCatalog: catalog.Name,
		Headers:     headers,
		Fields:      fields,
		Vals:        data,
		ValuesId:    valuesId,
		IsList:      isList,
	}
	return err, result
}
func SaveCatalogWork(user *model.User, catalog model.Catalog) (err error, id int) {
	db, _ := database.GetDb(user.ConnString)
	defer db.Close()
	var fieldsA []string
	for _, cat := range catalog.Fields {
		if cat.Name != "" && cat.Value != "" {
			fieldsA = append(fieldsA, cat.NameDb)
		}

	}

	type DefaultField struct {
		isUserCreate, isUserModify, isDateCreate, isDateModify bool
		fieldName                                              string
	}
	var (
		query string
	)
	var df []DefaultField
	//дефолтные поля

	query = `select name_db, is_user_create, is_user_modify, is_date_create, is_date_modify from utills_catalog_fields where catalog_id = ` + strconv.Itoa(catalog.Id)
	rows, err := db.Query(query)
	defer rows.Close()
	for rows.Next() {
		var f DefaultField
		err = rows.Scan(&f.fieldName, &f.isUserCreate, &f.isUserModify, &f.isDateCreate, &f.isDateModify)
		if f.isUserCreate || f.isUserModify || f.isDateCreate || f.isDateModify {
			df = append(df, f)
		}
		//if f.isUserCreate && !df.isUserCreate {
		//	df.isUserCreate = true
		//}
		//if f.isUserModify && !df.isUserModify {
		//	df.isUserModify = true
		//}
		//if f.isDateCreate && !df.isDateCreate {
		//	df.isDateCreate = true
		//}
		//if f.isDateModify && !df.isDateModify {
		//	df.isDateModify = true
		//}
	}
	if catalog.EntityId == 0 {
		for _, f := range df {
			if f.isUserCreate || f.isDateCreate {
				fieldsA = append(fieldsA, f.fieldName)
			}

		}
		query = `insert into ` + catalog.TableName + ` (` + strings.Join(fieldsA, ",") + `) values (`
	} else {
		for _, f := range df {
			if f.isUserModify || f.isDateModify {
				fieldsA = append(fieldsA, f.fieldName)
			}

		}
		query = `update ` + catalog.TableName + ` set `
	}
	var vals []string
	for idx, cat := range catalog.Fields {
		if cat.Name != "" && cat.Value != "" {
			val := cat.Value
			if cat.NameType == "bit" {
				//	fmt.Println(val)
			}
			if cat.NameType != "int" {

				val = "'" + cat.Value + "'"
			}
			if catalog.EntityId != 0 {
				val = catalog.Fields[idx].NameDb + "=" + val
			}
			vals = append(vals, val)

		}
		//if cat.NameType != ''
	}
	if catalog.EntityId == 0 {
		for _, f := range df {
			if f.isUserCreate {
				vals = append(vals, "original_login()")
			}
			if f.isDateCreate {
				vals = append(vals, "GETDATE()")
			}
		}
	} else {
		for _, f := range df {
			if f.isUserModify {

				vals = append(vals, f.fieldName+"= original_login()")
			}
			if f.isDateModify {
				vals = append(vals, f.fieldName+"= GETDATE()")
			}
		}
	}
	query += strings.Join(vals, ",")
	if catalog.EntityId != 0 {
		query += ` where id = ` + strconv.Itoa(catalog.EntityId)
	} else {
		query += ") ;  SELECT SCOPE_IDENTITY()"
	}
	if catalog.EntityId != 0 {
		_, err = db.Exec(query)
	} else {
		err = db.QueryRow(query).Scan(&catalog.EntityId)

	}

	if err != nil {
		log.Println(err.Error())
	}
	return err, catalog.EntityId
}

//редактирование справочника  - записи
func GetEntityByCatalogId(user *model.User, catalogId, entityId int) (err error, data model.Catalog) {
	db, _ := database.GetDb(user.ConnString)
	defer db.Close()
	catalog := db_catalog.GetCatalogById(user, catalogId, true)

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
	//jsonStr = strings.ReplaceAll(jsonStr, "\"", "")
	jsonArr := strings.Split(jsonStr, ",")
	for _, js := range jsonArr {
		arr := strings.Split(js, ":")
		fld := arr[0]
		fld = fld[1 : len(fld)-1]
		for idx, field := range catalog.Fields {

			if field.NameDb == fld {
				if field.NameType == "bit" {

				}
				if field.NameType == "bit" && arr[1] == "1" {
					catalog.Fields[idx].Value = "checked"
				} else {
					val := arr[1]
					val = val[1 : len(val)-1]
					val = strings.ReplaceAll(val, "\\", "")
					catalog.Fields[idx].Value = val
				}
			}
		}

	}
	catalog.EntityId = entityId
	return err, catalog
}
func DeleteCatalogWorkList(user *model.User, id, catalogId int) (err error) {

	catalog := db_catalog.GetCatalogById(user, catalogId, false)
	db, _ := database.GetDb(user.ConnString)
	defer db.Close()
	query := `delete from  ` + catalog.TableName + ` where id = ` + strconv.Itoa(id)

	_, err = db.Exec(query)

	return err
}

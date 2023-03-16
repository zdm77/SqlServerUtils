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
	Data        [][]string `json:"data"`
	Json        string     `json:"json"`
}

func GetTableLinkById(user *model.User, id int) (err error, table string) {
	db, _ := database.GetDb(user.ConnString)
	defer db.Close()
	query := `select table_name from utils_catalog_list where id=` + strconv.Itoa(id)
	err = db.QueryRow(query).Scan(&table)
	return err, table
}

func GetCatalogWorkListByIdJson(user *model.User, id int, isJSONParse bool) (err error, result FieldVals) {

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
	var joinTables string
	type Access struct {
		id    string
		user  []string
		users string
	}
	isAccessInTable := false
	var whereAccess []string
	if user.SuperAdmin != user.Login {
		//ищем списки для проверки прав
		for _, cat := range catalog.Fields {
			if cat.LinkTableId != 0 && cat.IsAccessCheck {
				_, tableLink := GetTableLinkById(user, cat.LinkTableId)

				var ids []string
				queryAcc := `select id, access from ` + tableLink
				rows, err := db.Query(queryAcc)
				if err == nil {
					isAccessInTable = true
					defer rows.Close()
					for rows.Next() {

						var a string
						var wId string
						err = rows.Scan(&wId, &a)
						arr := strings.Split(a, ",")
						for _, a1 := range arr {
							a1 = strings.TrimSpace(a1)
							if a1 == user.Login {
								ids = append(ids, wId)
							}
						}

					}
					if len(ids) > 0 {
						wh := " and " + tableLink + ".id in (" + strings.Join(ids, ",") + ")"
						whereAccess = append(whereAccess, wh)
					}

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
				joinTables += " left join " + table + " on " + catalog.TableName + "." + cat.NameDb + "=" + table + ".id"
				field = "coalesce(cast(" + table + ".name as varchar(255)), '') " + cat.Name
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
	if user.SuperAdmin != user.Login {
		if isAccessInTable {
			if len(whereAccess) == 0 {
				return err, result
			}
			query += " where 1=1 "
			for _, data := range whereAccess {
				query += data
			}
			//query += ` where ` + catalog.TableName + `.id in (` + strings.Join(ids, ",") + `) `
		}
	}
	var jsonString string
	if isJSONParse {
		query += "  for json auto "

		err = db.QueryRow(query).Scan(&jsonString)
		if err != nil {
			log.Println(err.Error())
			return err, result
		}
	}
	result = FieldVals{
		CatalogId:   id,
		FieldId:     fieldId,
		NameCatalog: catalog.Name,
		Headers:     headers,
		Fields:      fields,
		Vals:        data,
		ValuesId:    valuesId,
		IsList:      isList,
		Data:        data,
		Json:        jsonString,
	}
	return err, result
}

func GetCatalogWorkListById(user *model.User, id int, isJSONParse bool) (err error, result FieldVals) {

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
	//queryAccess := `select id, access from ` + catalog.TableName
	type Access struct {
		id    string
		user  []string
		users string
	}
	//var access []Access
	//var ids []string
	isAccessInTable := false
	//if user.SuperAdmin != user.Login {
	//	rows, err := db.Query(queryAccess)
	//
	//	if err == nil {
	//		isAccessInTable = true
	//		for rows.Next() {
	//			var a Access
	//			rows.Scan(&a.id, &a.users)
	//			users := strings.Split(a.users, ",")
	//			for _, u := range users {
	//				u = strings.TrimSpace(u)
	//				if u == user.Login {
	//					ids = append(ids, a.id)
	//				}
	//			}
	//		}
	//	}
	//}
	var whereAccess []string
	if user.SuperAdmin != user.Login {
		//ищем списки для проверки прав
		for _, cat := range catalog.Fields {
			if cat.LinkTableId != 0 && cat.IsAccessCheck {
				_, tableLink := GetTableLinkById(user, cat.LinkTableId)

				var ids []string
				queryAcc := `select id, access from ` + tableLink
				rows, err := db.Query(queryAcc)
				if err == nil {
					isAccessInTable = true
					defer rows.Close()
					for rows.Next() {

						var a string
						var wId string
						err = rows.Scan(&wId, &a)
						arr := strings.Split(a, ",")
						for _, a1 := range arr {
							a1 = strings.TrimSpace(a1)
							if a1 == user.Login {
								ids = append(ids, wId)
							}
						}

					}

					if len(ids) > 0 {
						wh := " and " + tableLink + ".id in (" + strings.Join(ids, ",") + ")"
						whereAccess = append(whereAccess, wh)
					}

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
				joinTables += " left join " + table + " on " + catalog.TableName + "." + cat.NameDb + "=" + table + ".id"
				field = "coalesce(cast(" + table + ".name as varchar(255)), '') " + cat.Name
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
	if user.SuperAdmin != user.Login {
		if isAccessInTable {
			if len(whereAccess) == 0 {
				return err, result
			}
			query += " where 1=1 "
			for _, data := range whereAccess {
				query += data
			}
			//query += ` where ` + catalog.TableName + `.id in (` + strings.Join(ids, ",") + `) `
		}
	}
	var jsonString string
	if isJSONParse {
		query += "  for json auto "

		err = db.QueryRow(query).Scan(&jsonString)
		if err != nil {
			log.Println(err.Error())
			return err, result
		}
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
					if nameDb == `"`+fieldId+`"` {
						valuesId = append(valuesId, val)
					}
					dt = append(dt, val[1:len(val)-1])

				}
				data = append(data, dt)
			}

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
		Data:        data,
		Json:        jsonString,
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
	queryAccess := `select access from ` + catalog.TableName
	isAccessTable := false
	var acc string
	err = db.QueryRow(queryAccess).Scan(&acc)
	if err == nil {
		isAccessTable = true
	}
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
		if isAccessTable {
			fieldsA = append(fieldsA, "access")
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
	if isAccessTable {
		vals = append(vals, "'"+user.Login+"'")
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

func GetCatalogAccessRecord(user *model.User, id int, table string) (err error, result []model.AccessRecord) {
	db, _ := database.GetDb(user.ConnString)
	defer db.Close()
	query := `select access from ` + table + ` where id = ` + strconv.Itoa(id)
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

func SaveCatalogAccessRecord(user *model.User, id int, table, access string) (err error) {
	db, _ := database.GetDb(user.ConnString)
	defer db.Close()
	query := `update ` + table + ` set access = '` + access + `'  where id = ` + strconv.Itoa(id)
	_, err = db.Exec(query)
	if err != nil {
		log.Println(err.Error())
	}
	return err
}

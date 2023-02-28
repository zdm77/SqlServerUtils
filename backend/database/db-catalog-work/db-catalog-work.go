package db_catalog_work

import (
	"fmt"
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
}

func GetCatalogWorkListById(user *model.User, id int) (err error, result FieldVals) {

	db, _ := database.GetDb(user.ConnString)
	defer db.Close()

	catalog := db_catalog.GetCatalogById(user, id)
	var data [][]string
	var headers []string
	query := `select `
	var fields []string
	var fieldId string
	var valuesId []string
	for _, cat := range catalog.Fields {
		if cat.IsPrimaryKey {
			fieldId = cat.NameDb
		}
		if cat.IsIdentity || cat.Name != "" {
			field := "coalesce(cast(" + cat.NameDb + " as varchar(255)), '') " + cat.NameDb
			fields = append(fields, field)
			if !cat.IsIdentity {
				headers = append(headers, cat.Name)
			} else {
				headers = append(headers, cat.NameDb)
			}
		}
	}
	query += strings.Join(fields, ",") + " from " + catalog.TableName + "  for json auto "
	var jsonString string
	err = db.QueryRow(query).Scan(&jsonString)
	jsonArr := strings.Split(jsonString, "{")
	//header := make(map[string]bool)
	for idx, js := range jsonArr {
		if idx != 0 {
			str := strings.ReplaceAll(js, "},", "")
			arr := strings.Split(str, ",")
			var dt []string
			for idxa, a := range arr {
				vl := strings.Split(a, ":")
				nameDb := vl[0]

				val := strings.ReplaceAll(vl[1], `}]`, "")
				if catalog.Fields[idxa].NameType == "bit" {
					if val == `"1"` {
						val = "checked"
					}
				}
				if nameDb == `"`+fieldId+`"` {
					valuesId = append(valuesId, val)
				}
				dt = append(dt, val)
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
	var query string
	if catalog.EntityId == 0 {
		query = `insert into ` + catalog.TableName + ` (` + strings.Join(fieldsA, ",") + `) values (`
	} else {
		query = `update ` + catalog.TableName + ` set `
	}
	var vals []string
	for idx, cat := range catalog.Fields {
		if cat.Name != "" && cat.Value != "" {
			val := cat.Value
			if cat.NameType == "bit" {
				fmt.Println(val)
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
	catalog := db_catalog.GetCatalogById(user, catalogId)

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

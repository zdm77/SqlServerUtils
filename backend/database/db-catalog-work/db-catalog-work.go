package db_catalog_work

import (
	"log"
	"sqlutils/backend/database"
	"sqlutils/backend/database/db_catalog"
	"sqlutils/backend/model"
	"strings"
)

type FieldVals struct {
	CatalogId   int        `json:"catalog_id"`
	FieldId     string     `json:"field_id"`
	NameCatalog string     `json:"name_catalog"`
	Headers     []string   `json:"headers"`
	Fields      []string   `json:"fields"`
	Vals        [][]string `json:"vals"`
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
			for _, a := range arr {
				vl := strings.Split(a, ":")
				//val := strings.ReplaceAll(vl[1], `"`, "")
				val := strings.ReplaceAll(vl[1], `}]`, "")
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
	query := `insert into ` + catalog.TableName + ` (` + strings.Join(fieldsA, ",") + `) values (`
	var vals []string
	for _, cat := range catalog.Fields {
		if cat.Name != "" && cat.Value != "" {
			val := cat.Value
			if cat.NameType != "int" {
				val = "'" + cat.Value + "'"
			}
			vals = append(vals, val)
		}
		//if cat.NameType != ''
	}
	query += strings.Join(vals, ",") + ")"
	_, err = db.Exec(query)
	if err != nil {
		log.Println(err.Error())
	}
	return err, 0
}

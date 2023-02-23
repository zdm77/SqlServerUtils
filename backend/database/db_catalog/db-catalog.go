package db_catalog

import (
	"fmt"
	"log"
	"sqlutils/backend/database"
	"sqlutils/backend/model"
)

func GetDbTableFields(tableName string) {
	query := `SELECT 
    c.name col_name,
    t.Name col_type,
    c.max_length 'Max Length',
    c.precision ,
    c.scale ,
    c.is_nullable,
    ISNULL(i.is_primary_key, 0) 'Primary Key'
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
	fmt.Println(query)
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

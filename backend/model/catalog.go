package model

type Catalog struct {
	Id        int     `json:"id"`
	Name      string  `json:"name"`
	TableName string  `json:"table_name"`
	Fields    []Field `json:"fields"`
}
type Field struct {
	Id           int    `json:"id"`
	CatalogId    int    `json:"catalog_id"`
	NameDb       string `json:"name_db"`
	NameType     string `json:"name_type"`
	MaxLength    int    `json:"max_length"`
	Precision    int    `json:"precision"`
	Scale        int    `json:"scale"`
	IsNullable   bool   `json:"is_nullable"`
	IsNullableDb bool   `json:"is_nullable_db"`
	//автоинкремент
	IsIdentity   bool   `json:"is_identity"`
	IsPrimaryKey bool   `json:"is_primary_key"`
	Name         string `json:"name"`
}

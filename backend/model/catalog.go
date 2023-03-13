package model

type Catalog struct {
	Id        int     `json:"id"`
	EntityId  int     `json:"entity_id"`
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
	Value        string `json:"value"`
	IsList       bool   `json:"is_list"`
	IsForm       bool   `json:"is_form"`
	LinkTableId  int    `json:"link_table_id"`
}
type LinkTable struct {
	Id int `json:"id"`
}

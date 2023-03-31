package model

type Catalog struct {
	Id                int     `json:"id"`
	EntityId          int     `json:"entity_id"`
	Name              string  `json:"name"`
	TableName         string  `json:"table_name"`
	Fields            []Field `json:"fields"`
	TypeEntity        string  `json:"type_entity"`
	IsCatalogTask     bool    `json:"is_catalog_task"`
	OrderByDefault    string  `json:"order_by_default"`
	OrderByDefaultAsc string  `json:"order_by_default_asc"`
	PKeyName          string  `json:"p_key_name"`
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
	IsIdentity     bool   `json:"is_identity"`
	IsPrimaryKey   bool   `json:"is_primary_key"`
	Name           string `json:"name"`
	Value          string `json:"value"`
	IsList         bool   `json:"is_list"`
	IsForm         bool   `json:"is_form"`
	LinkTableId    int    `json:"link_table_id"`
	LinkFieldView  string `json:"link_field_view"`
	IsUserCreate   bool   `json:"is_user_create"`
	IsUserModify   bool   `json:"is_user_modify"`
	IsDateCreate   bool   `json:"is_date_create"`
	IsDateModify   bool   `json:"is_date_modify"`
	TableName      string `json:"table_name"`
	IsAccessCheck  bool   `json:"is_access_check"`
	IsForeignField bool   `json:"is_foreign_field"`
	OrderBy        int    `json:"order_by"`
	OrderByForm    int    `json:"order_by_form"`
	IsNewField     bool   `json:"is_new_field"`
}
type LinkTable struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Access string `json:"access"`
}
type Script struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	ScriptName string `json:"script_name"`
}

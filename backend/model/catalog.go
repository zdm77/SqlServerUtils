package model

type Catalog struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	TableName string `json:"table_name"`
}

package model

type Task struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	TableDb string `json:"table_db"`
}

package model

type Task struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	TableDb   string `json:"table_db"`
	StrHeader int    `json:"str_header"`
}
type TaskParams struct {
	Id         int    `json:"id"`
	TaskId     int    `json:"task_id"`
	FieldExcel string `json:"field_excel"`
	FieldDb    string `json:"field_db"`
}

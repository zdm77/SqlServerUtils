package model

type Task struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	TableDb     string `json:"table_db"`
	StrHeader   int    `json:"str_header"`
	SheetNumber int    `json:"sheet_number"`
	FileUpload  string `json:"file_upload"`
}
type TaskParams struct {
	Id         int    `json:"id"`
	TaskId     int    `json:"task_id"`
	FieldExcel string `json:"field_excel"`
	FieldDb    string `json:"field_db"`
	Dir        string `json:"dir"`
	FieldType  string `json:"field_type"`
	Value      string `json:"value"`
	IsError    bool   `json:"is_error"`
}
type Result struct {
	Headers    []TaskParams   `json:"headers"`
	Data       [][]TaskParams `json:"data"`
	FileUpload string         `json:"file_upload"`
	Task       Task           `json:"task"`
}

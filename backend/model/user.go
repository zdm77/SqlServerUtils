package model

type User struct {
	Login      string
	IsLogin    bool
	ConnString string
	DbName     string
	SuperAdmin string
}

type AccessRecord struct {
	UserName string `json:"user_name"`
	Access   bool   `json:"access"`
}

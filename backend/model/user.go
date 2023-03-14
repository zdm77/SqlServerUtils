package model

type User struct {
	Login      string
	IsLogin    bool
	ConnString string
	DbName     string
}

type AccessRecord struct {
	UserName string
	Access   bool
}

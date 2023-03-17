package database

import (
	"sqlutils/backend/model"
	"strconv"
)

type TestS struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	F    string `json:"f"`
}
type TestData struct {
	RecordsTotal    int     `json:"recordsTotal"`
	RecordsFiltered int     `json:"recordsFiltered"`
	Data            []TestS `json:"data"`
	//Length       int     `json:"length"`
}

func Test(user *model.User, start int) (data TestData) {
	var dt []TestS
	//for i := start; i < start+10; i++ {
	for i := 0; i < 1000; i++ {
		var f string

		if i < 20 {
			f = "2023-03-03"
		} else {
			f = "2023-03-20"
		}
		dt = append(dt, TestS{
			Id:   i,
			Name: "t" + strconv.Itoa(i),

			F: f,
		})
		//var dt TestS
		//dt = append(dt, TestS{
		//	Id:   i,
		//	Name: "Тест " + strconv.Itoa(i),
		//})
	}
	data = TestData{
		RecordsTotal:    50000,
		Data:            dt,
		RecordsFiltered: 50000,
		//Length:       100000,
	}
	return data
}

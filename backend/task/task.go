package task

import (
	"github.com/tealeg/xlsx/v3"
	"log"
	"sqlutils/backend/database"
	db_task "sqlutils/backend/database/db-task"
	"sqlutils/backend/model"
	"strconv"
	"strings"
	"time"
)

func ExecTaskFromExcel(user *model.User, fileExe string, taskId int) {
	task := db_task.GetTaskById(user, taskId)
	params := db_task.GetTaskParams(user, taskId, true)

	wb, err := xlsx.OpenFile(fileExe)
	if err != nil {
		log.Println(err.Error())
	}
	sh := wb.Sheets[0]
	//заголовки
	type Data struct {
		Col     int
		FieldDb string
	}

	//	var dataParam []Data
	var fields []string
	var values []string
	for _, p := range params {
		fields = append(fields, p.FieldDb)
	}
	db, _ := database.GetDb(user.ConnString)
	defer db.Close()
	for i := task.StrHeader; i < sh.MaxRow; i++ {
		values = nil
		query := `insert into ` + task.TableDb
		for _, p := range params {
			valEx, _ := sh.Cell(i, p.Id)
			vl := valEx.String()

			if valEx.NumFmt != "general" && valEx.NumFmt != "" && valEx.NumFmt != "0.00" && valEx.NumFmt != "0" {
				t, err := strconv.ParseInt(vl, 10, 64)
				if err == nil {
					tUx := (t - 25569) * 86400
					dt := time.Unix(tUx, 0)
					//valIns = dt.Format("01-02-2006")
					values = append(values, `'`+dt.Format("01-02-2006")+`'`)
				}
			} else {

				_, err := strconv.Atoi(vl)
				if err != nil {

					values = append(values, `'`+vl+`'`)
				} else {
					values = append(values, vl)
				}
			}

		}
		query += `(` + strings.Join(fields, ",") + `) values (` + strings.Join(values, ",") + `)`
		//	fmt.Println(query)
		_, err = db.Exec(query)
		if err != nil {
			return
			log.Println(err.Error())
		}
	}

	//for i := 0; i < sh.MaxCol; i++ {
	//	head, _ := sh.Cell(task.StrHeader-1, i)
	//	fmt.Println(head.String())
	//	for _, p := range params {
	//		if p.FieldExcel == head.String() {
	//
	//			dataParam = append(dataParam, Data{
	//				Col:     i,
	//				FieldDb: p.FieldDb,
	//			})
	//			continue
	//		}
	//	}
	//}
	//for i := task.StrHeader; i < sh.MaxRow; i++ {
	//	query := `insert into ` + task.TableDb
	//	//var fields []string
	//
	//	for _, val := range dataParam {
	//		fields = append(fields, val.FieldDb)
	//		v, _ := sh.Cell(i, val.Col)
	//		vl := v.String()
	//		_, err := strconv.Atoi(vl)
	//		if err != nil {
	//			values = append(values, `'`+v.String()+`'`)
	//		} else {
	//			values = append(values, v.String())
	//		}
	//
	//	}
	//	query += `(` + strings.Join(fields, ",") + `) values (` + strings.Join(values, ",") + `)`
	//	fmt.Println(query)
	//	db, _ := database.GetDb(user.ConnString)
	//	defer db.Close()
	//	_, err = db.Exec(query)
	//	if err != nil {
	//		return
	//		log.Println(err.Error())
	//	}
	//
	//}
	//for i := task.StrHeader - 1; i < sh.MaxRow; i++ {
	//	var qqq []string
	//	for j := 0; j < sh.MaxCol; j++ {
	//		val, _ := sh.Cell(i, j)
	//		valIns := val.String()
	//		if val.NumFmt != "general" && val.NumFmt != "" && val.NumFmt != "0.00" && val.NumFmt != "0" {
	//			t, err := strconv.ParseInt(valIns, 10, 64)
	//			if err == nil {
	//				tUx := (t - 25569) * 86400
	//				dt := time.Unix(tUx, 0)
	//				valIns = dt.Format("01-02-2006")
	//				fmt.Println(valIns)
	//			}
	//		}
	//		qqq = append(qqq, valIns)
	//	}
	//
	//fmt.Println(qqq)

	//}
}
func GetHeaders(user *model.User, fileExe string, taskId int) (params []model.TaskParams) {
	task := db_task.GetTaskById(user, taskId)
	wb, err := xlsx.OpenFile(fileExe)
	if err != nil {
		log.Println(err.Error())
	}

	sh := wb.Sheets[0]
	//заголовки
	for i := 0; i < sh.MaxCol; i++ {
		v, _ := sh.Cell(task.StrHeader-1, i)
		params = append(params, model.TaskParams{
			Id:         i,
			TaskId:     taskId,
			FieldExcel: v.String(),
			FieldDb:    "",
		})
	}

	return params
}

type Result struct {
	Headers []model.TaskParams   `json:"headers"`
	Data    [][]model.TaskParams `json:"data"`
}

func GetData(user *model.User, fileExe string, taskId int) (data Result) {
	task := db_task.GetTaskById(user, taskId)
	wb, err := xlsx.OpenFile(fileExe)
	if err != nil {
		log.Println(err.Error())
	}

	sh := wb.Sheets[0]
	var headers []model.TaskParams
	params := db_task.GetTaskParams(user, taskId, true)
	//заголовки
	for i := 0; i < sh.MaxCol; i++ {
		v, _ := sh.Cell(task.StrHeader-1, i)
		for _, p := range params {
			if p.FieldExcel == v.String() && p.Id == i {
				headers = append(headers, p)
			}
		}

	}
	//данные
	var dataO [][]model.TaskParams
	for i := task.StrHeader; i < sh.MaxRow; i++ {
		var dt []model.TaskParams
		for j := 0; j < sh.MaxCol; j++ {

			v, _ := sh.Cell(i, j)
			for _, p := range params {
				if p.Id == j {
					dt = append(dt, model.TaskParams{
						Id:         j,
						TaskId:     taskId,
						FieldExcel: v.String(),
						FieldDb:    p.FieldDb,
						Dir:        "",
						FieldType:  p.FieldType,
						Value:      v.String(),
					})
				}
			}

		}
		dataO = append(dataO, dt)
	}
	data = Result{
		Headers: headers,
		Data:    dataO,
	}
	return data
}

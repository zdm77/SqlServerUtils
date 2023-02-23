package task

import (
	"fmt"
	"github.com/tealeg/xlsx/v3"
	"log"
	"sqlutils/backend/database"
	db_task "sqlutils/backend/database/db-task"
	"sqlutils/backend/model"
	"strconv"
	"strings"
	"time"
)

func ExecTaskFromExcel(user *model.User, result model.Result) (err error) {
	//task := db_task.GetTaskById(user, taskId)
	//params := db_task.GetTaskParams(user, taskId, true)

	//	var dataParam []Data
	var fields []string
	var values []string
	table := result.Task.TableDb
	db, _ := database.GetDb(user.ConnString)
	tx, _ := db.Begin()
	defer db.Close()
	for _, head := range result.Headers {
		fields = append(fields, head.FieldDb)
	}

	for _, res := range result.Data {
		values = nil
		query := `insert into  ` + table
		for _, d := range res {
			val := d.Value
			if d.FieldType != "Число" {
				val = `'` + val + `'`
			}
			values = append(values, val)
		}
		query += `(` + strings.Join(fields, ",") + `) values (` + strings.Join(values, ",") + `)`
		fmt.Println(query)
		_, err = tx.Exec(query)
		if err != nil {
			tx.Rollback()
			return err
			log.Println(err.Error())
		}
	}
	tx.Commit()
	return err
	//
	//
	//		if valEx.NumFmt != "general" && valEx.NumFmt != "" && valEx.NumFmt != "0.00" && valEx.NumFmt != "0" {
	//			t, err := strconv.ParseInt(vl, 10, 64)
	//			if err == nil {
	//				tUx := (t - 25569) * 86400
	//				dt := time.Unix(tUx, 0)
	//				//valIns = dt.Format("01-02-2006")
	//				values = append(values, `'`+dt.Format("01-02-2006")+`'`)
	//			}
	//		} else {
	//
	//			_, err := strconv.Atoi(vl)
	//			if err != nil {
	//
	//				values = append(values, `'`+vl+`'`)
	//			} else {
	//				values = append(values, vl)
	//			}
	//		}
	//
	//	}
	//	query += `(` + strings.Join(fields, ",") + `) values (` + strings.Join(values, ",") + `)`
	//	//	fmt.Println(query)
	//	_, err = db.Exec(query)
	//	if err != nil {
	//		return
	//		log.Println(err.Error())
	//	}
	//}

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
func GetHeaders(user *model.User, fileExe string, strHeader int) (params []model.TaskParams) {
	//task := db_task.GetTaskById(user, taskId)
	wb, err := xlsx.OpenFile(fileExe)
	if err != nil {
		log.Println(err.Error())
	}

	sh := wb.Sheets[0]
	//заголовки
	for i := 0; i < sh.MaxCol; i++ {
		v, _ := sh.Cell(strHeader-1, i)
		params = append(params, model.TaskParams{
			Id:         i,
			TaskId:     0,
			FieldExcel: v.String(),
			FieldDb:    "",
		})
	}

	return params
}

func GetData(user *model.User, fileExe string, taskId int) (data model.Result) {
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
			val := v.String()
			for _, p := range params {
				if p.Id == j {
					isError := false
					//fmt.Println(p.FieldType)
					if p.FieldType == "Дата" {
						t, err := strconv.ParseInt(val, 10, 64)
						if err == nil {
							tUx := (t - 25569) * 86400
							dt := time.Unix(tUx, 0)
							//valIns = dt.Format("01-02-2006")
							//values = append(values, `'`+dt.Format("01-02-2006")+`'`)
							val = dt.Format("01-02-2006")
						} else {
							isError = true
						}
					}
					if p.FieldType == "Число" {
						_, err = strconv.Atoi(val)
						if err != nil {
							//возможно вещественное
							_, err = strconv.ParseFloat(val, 64)
							if err != nil {
								isError = true
							}
						}
					}
					dt = append(dt, model.TaskParams{
						Id:         j,
						TaskId:     taskId,
						FieldExcel: p.FieldExcel,
						FieldDb:    p.FieldDb,
						Dir:        "",
						FieldType:  p.FieldType,
						Value:      val,
						IsError:    isError,
					})
				}
			}

		}
		dataO = append(dataO, dt)
	}
	data = model.Result{
		Headers:    headers,
		Data:       dataO,
		FileUpload: fileExe,
	}
	return data
}

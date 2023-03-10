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
}
func GetHeaders(fileExe string, strHeader int, sheetNumber int) (params []model.TaskParams) {

	wb, err := xlsx.OpenFile(fileExe)
	if err != nil {
		log.Println(err.Error())
	}

	sh := wb.Sheets[sheetNumber-1]
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
	task := db_task.GetTaskCatalogById(user, taskId)
	wb, err := xlsx.OpenFile(fileExe)
	if err != nil {
		log.Println(err.Error())
	}

	sh := wb.Sheets[task.SheetNumber-1]
	var headers []model.TaskParams
	params := db_task.GetTaskCatalogParams(user, taskId, true)
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

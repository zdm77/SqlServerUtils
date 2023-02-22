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
	mapHeader := make(map[string]Data)

	for i := 0; i < sh.MaxCol; i++ {
		head, _ := sh.Cell(task.StrHeader-1, i)
		for _, p := range params {
			if p.FieldExcel == head.String() {
				mapHeader[p.FieldExcel] = Data{
					Col:     i,
					FieldDb: p.FieldDb,
				}
			}
		}
	}
	for i := task.StrHeader; i < sh.MaxRow; i++ {
		query := `insert into ` + task.TableDb
		var fields []string
		var values []string
		for _, val := range mapHeader {
			fields = append(fields, val.FieldDb)
			v, _ := sh.Cell(i, val.Col)
			vl := v.String()
			_, err := strconv.Atoi(vl)
			if err != nil {
				values = append(values, `'`+v.String()+`'`)
			} else {
				values = append(values, v.String())
			}

		}
		query += `(` + strings.Join(fields, ",") + `) values (` + strings.Join(values, ",") + `)`
		fmt.Println(query)
		db, _ := database.GetDb(user.ConnString)
		defer db.Close()
		_, err = db.Exec(query)
		if err != nil {
			return
			log.Println(err.Error())
		}

	}
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

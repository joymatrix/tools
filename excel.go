package utils

import (
	"fmt"
	"math/rand"

	"strconv"
	"time"

	"github.com/xuri/excelize/v2"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func InitJobData() ([]model.Job, error) {
	jobs := make([]model.Job, 0)
	f, err := excelize.OpenFile("./jobInfo.xlsx")
	if err != nil {
		fmt.Println("open excel file err:", err)
		return jobs, err
	}
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		fmt.Println("get rows err:%+v", err.Error())
		return jobs, err
	}
	for _, row := range rows[1:] {
		fmt.Println("job row:", row)
		id, err := strconv.ParseInt(row[0], 10, 64)
		if err != nil {
			fmt.Println("job err:", err.Error())
			return jobs, err
		}
		var parentId int64
		if len(row) == 3 {
			parentId, err = strconv.ParseInt(row[2], 10, 64)
			if err != nil {
				fmt.Println("job err:", err.Error())
				return jobs, err
			}
		}

		tempJob := model.Job{
			BaseModel: model.BaseModel{ID: id},
			JobName:   row[1],
			ParentId:  parentId,
		}
		jobs = append(jobs, tempJob)
	}
	return jobs, nil
}

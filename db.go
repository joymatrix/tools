package dao

import (
	"fmt"
	"sanicalc/configs"
	"sanicalc/internal/model"
	"sanicalc/internal/utils"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDb() {
	dbInfo := configs.Config.Postgres

	dsn := fmt.Sprintf("host=%s port=%d database=%s user=%s password=%s", dbInfo.Host, dbInfo.Port, dbInfo.DbName, dbInfo.Username, dbInfo.Password)
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		// 处理连接错误
		panic(err)
	}
	db = db.Debug()
	db.AutoMigrate(&model.User{}, &model.Coupon{}, &model.Invoice{}, &model.Order{}, &model.Project{}, &model.ProjectJob{}, &model.Job{}, &model.ProjectJobCraftVehicle{}, &model.ProjectJobParams{}, &model.ProjectVehicleRecycle{}, &model.JobCraftVehicle{})

	utils.GetLog().Info("conn db success")
}

func GetDb() *gorm.DB {
	return db
}

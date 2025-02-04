package config

import (
	"fmt"
	"go_graphql/internal/dto"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	dsn := GetDSN()
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	err = db.AutoMigrate(&dto.TenantResource{}, &dto.TenantMetadata{}, &dto.TNTRole{}, &dto.TNTPermission{}, &dto.TNTRolePermission{}, &dto.MstRole{}, &dto.MstPermission{}, &dto.MstRolePermission{})
	if err != nil {
		panic(err)
	}

	DB = db

	return db
}

func GetDSN() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%s&loc=%s",
		GetEnv("DB_USERNAME"), GetEnv("DB_PASSWORD"), GetEnv("DB_HOST"), GetEnv("DB_PORT"), GetEnv("DB_NAME"), DBCharset, DBParseTime, DBLoc,
	)
}

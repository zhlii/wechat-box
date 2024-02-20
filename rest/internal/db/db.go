package db

import (
	"fmt"

	"github.com/zhlii/wechat-box/rest/internal/config"
	"github.com/zhlii/wechat-box/rest/internal/db/tables"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var Db *gorm.DB

func Connect() error {
	cfg := config.Data.Database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:               dsn,
		DefaultStringSize: 256,
	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		return err
	}

	Db = db

	db.AutoMigrate(&tables.Message{})

	return nil
}

func Destory() error {
	if db, err := Db.DB(); db != nil {
		return db.Close()
	} else {
		return err
	}
}

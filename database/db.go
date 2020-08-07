package database

import (
	"database/sql"
	"fmt"

	"github.com/apulis/AIArtsBackend/configs"
	"github.com/apulis/AIArtsBackend/loggers"
	_ "github.com/apulis/AIArtsBackend/loggers"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var Db *gorm.DB
var logger = loggers.Log

func init() {

	dbConf := configs.Config.Db

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=utf8&parseTime=True&loc=Local",
		dbConf.Username, dbConf.Password, dbConf.Host, dbConf.Port))

	defer db.Close()
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS " + dbConf.Database)
	if err != nil {
		panic(err)
	}

	Db, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		dbConf.Username, dbConf.Password, dbConf.Host, dbConf.Port, dbConf.Database))
	if err != nil {
		panic(err)
	}

	logger.Info("DB connected success")
	Db.DB().SetMaxOpenConns(dbConf.MaxOpenConns)
	Db.DB().SetMaxIdleConns(dbConf.MaxIdleConns)
}

package database

import (
	"database/sql"
	"fmt"

	"github.com/apulis/AIArtsBackend/configs"
	"github.com/apulis/AIArtsBackend/loggers"
	_ "github.com/apulis/AIArtsBackend/loggers"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var Db *gorm.DB
var logger = loggers.Log

func init() {

	dbConf := configs.Config.Db

	if configs.Config.Database == "mysql"{

		db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=utf8&parseTime=True&loc=Local",
			dbConf.Mysql.Username, dbConf.Mysql.Password, dbConf.Mysql.Host, dbConf.Mysql.Port))

		defer db.Close()
		if err != nil {
			panic(err)
		}

		_, err = db.Exec("CREATE DATABASE IF NOT EXISTS " + dbConf.Mysql.Database)
		if err != nil {
			panic(err)
		}

		Db, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
			dbConf.Mysql.Username, dbConf.Mysql.Password, dbConf.Mysql.Host, dbConf.Mysql.Port, dbConf.Mysql.Database))
		if err != nil {
			panic(err)
		}

		logger.Info("DB connected success")
		Db.DB().SetMaxOpenConns(dbConf.Mysql.MaxOpenConns)
		Db.DB().SetMaxIdleConns(dbConf.Mysql.MaxIdleConns)
	}else{

		db, err := sql.Open("postgres",fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable",
			dbConf.PostgreSQL.Host,dbConf.PostgreSQL.Port, dbConf.PostgreSQL.Username, dbConf.PostgreSQL.Password))

		defer db.Close()
		if err != nil {
			panic(err)
		}

		exist := 0
		query := "SELECT count(1) FROM pg_database WHERE datname = $1"

		if err = db.QueryRow(query,dbConf.PostgreSQL.Database).Scan(
			&exist,
			); err != nil {
			logger.Info(err.Error())
			panic(err)
		}
		logger.Info(exist)

		if exist == 0 {
			logger.Info(dbConf.PostgreSQL.Database)
			_, err = db.Exec("CREATE DATABASE " + dbConf.PostgreSQL.Database)
			if err != nil {
				logger.Info(err.Error())
				panic(err)
			}
		}

		Db, err = gorm.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			dbConf.PostgreSQL.Host,dbConf.PostgreSQL.Port, dbConf.PostgreSQL.Username, dbConf.PostgreSQL.Password, dbConf.PostgreSQL.Database))
		if err != nil {
			logger.Info(err.Error())
			panic(err)
		}

		logger.Info("PostgreSQL connected success")
		Db.DB().SetMaxOpenConns(dbConf.PostgreSQL.MaxOpenConns)
		Db.DB().SetMaxIdleConns(dbConf.PostgreSQL.MaxIdleConns)
	}

}

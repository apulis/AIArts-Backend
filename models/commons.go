package models

import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/apulis/AIArtsBackend/database"
	"github.com/apulis/AIArtsBackend/loggers"
)

var db = database.Db
var logger = loggers.Log

type UnixTime struct {
	time.Time
}

func init() {
	createTableIfNotExists(Dataset{})
	createTableIfNotExists(Modelset{})
}

func createTableIfNotExists(modelType interface{}) {
	val := reflect.Indirect(reflect.ValueOf(modelType))
	modelName := val.Type().Name()

	hasTable := db.HasTable(modelType)
	if !hasTable {
		logger.Info(fmt.Sprintf("Table of %s not exists, create it.", modelName))
		db.CreateTable(modelType)
	} else {
		logger.Info(fmt.Sprintf("Table of %s already exists.", modelName))
	}
}

func (t UnixTime) MarshalJSON() ([]byte, error) {
	microSec := t.Unix() * 1000
	return []byte(strconv.FormatInt(microSec, 10)), nil
}

func (t UnixTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

func (t *UnixTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = UnixTime{Time: value}
		return nil
	}
	return fmt.Errorf("cannot convert %v to timestamp", v)
}

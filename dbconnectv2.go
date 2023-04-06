// Copyright 2021 Alexey Yanchenko <mail@yanchenko.me>
//
// This file is part of the Gufo library.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gufodao

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	//	"gorm.io/driver/sqlite"
	viper "github.com/spf13/viper"
)

// DB struct
//Uncomment after remove v1 Librry

type DBv2 struct {
	Conn *gorm.DB
}

// Connection instance
var DBConnectionv2 = &DBv2{}

func DBConnectv2() (*DBv2, error) {
	dbtype := viper.GetString("database.type")
	user := viper.GetString("database.user")
	pass := viper.GetString("database.password")
	pass = DecryptConfigPasswords(pass)
	dbname := viper.GetString("database.dbname")
	host := viper.GetString("database.host")
	port := viper.GetString("database.port")
	charset := viper.GetString("database.charset")
	sslmode := viper.GetString("database.sslmode")

	var request string

	var err error

	switch dbtype {
	case "mysql":
		request = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true", user, pass, host, port, dbname, charset)
		db, err := gorm.Open(mysql.Open(request), &gorm.Config{})
		DBConnectionv2.Conn = db
		return DBConnectionv2, err
	case "postgres":
		request = fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", host, port, user, dbname, pass, sslmode)
		//	dsn := "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"
		db, err := gorm.Open(postgres.Open(request), &gorm.Config{})
		DBConnectionv2.Conn = db
		return DBConnectionv2, err
	default:
		return nil, err
	}

}

func DBCheck() bool {
	_, err := DBConnectv2()
	//defer db.Close()
	if err != nil {
		SetErrorLog("dbconnectv2.go:77: " + err.Error())
		return false
	} else {
		return true
	}
}

// GetConnection - connect to DB
func ConnectDBv2() (*DBv2, error) {
	db, err := DBConnectv2()
	if err != nil {
		return db, err
	}

	sqlDB, err := db.Conn.DB()
	if err != nil {
		return db, err
	}

	dbcon := viper.GetInt("database.connectionssize")
	dbpool := viper.GetInt("database.poolsize")

	sqlDB.SetMaxIdleConns(dbcon)
	sqlDB.SetMaxOpenConns(dbpool)

	return db, nil
}

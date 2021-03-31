// Copyright 2020 Alexey Yanchenko <mail@yanchenko.me>
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

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	// import _ "github.com/jinzhu/gorm/dialects/sqlite"
	// import _ "github.com/jinzhu/gorm/dialects/mssql"
	viper "github.com/spf13/viper"
)

// DB struct
type DB struct {
	Conn *gorm.DB
}

// Connection instance
var DBConnection = &DB{}

func DBConnect() (*DB, error) {
	dbtype := viper.GetString("database.type")
	user := viper.GetString("database.user")
	pass := viper.GetString("database.password")
	pass = DecryptConfigPasswords(pass)
	dbname := viper.GetString("database.dbname")
	host := viper.GetString("database.host")
	port := viper.GetString("database.port")

	var request string

	switch dbtype {
	case "mysql":
		request = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, pass, host, port, dbname)
	case "postgres":
		request = fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s", host, port, user, dbname, pass)
	}
	db, err := gorm.Open(dbtype, request)
	DBConnection.Conn = db
	return DBConnection, err

}

func DBCheck() bool {
	_, err := DBConnect()
	//defer db.Close()
	if err != nil {
		SetErrorLog("dbconnect.go:65: " + err.Error())
		return false
	} else {
		return true
	}
}

// GetConnection - connect to DB
func ConnectDB() (*DB, error) {
	db, err := DBConnect()
	if err != nil {
		return db, err
	}
	err = db.Conn.DB().Ping()
	if err != nil {
		return db, err
	}
	return db, nil
}

// CloseConnection close connection db
func CloseConnection(db *DB) {
	/*
		SetErrorLog("I close db")

		err := db.Conn.Close()
		if err != nil {
			SetErrorLog("dbconnect.go:89: " + err.Error())
		}
	*/
}

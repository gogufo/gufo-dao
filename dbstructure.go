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
	"time"

	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

type Users struct {
	gorm.Model
	UID           string `gorm:"type:varchar(60);UNIQUE;NOT NULL;"`
	Name          string `gorm:"type:varchar(60);NOT NULL;DEFAULT '';UNIQUE"`
	Pass          string `gorm:"type:varchar(128);NOT NULL;DEFAULT ''"`
	Mail          string `gorm:"type:varchar(254);DEFAULT '';UNIQUE"`
	Mailsent      int    `gorm:"type:int(11);DEFAULT '0'"`
	Mailconfirmed int    `gorm:"type:int(11);DEFAULT '0'"`
	Created       int    `gorm:"type:int(11);DEFAULT '0'"`
	Access        int    `gorm:"type:int(11);DEFAULT '0'"`
	Login         int    `gorm:"type:int(11);DEFAULT '0'"`
	Status        int    `gorm:"type:double;DEFAULT 'false'"`
	Completed     int    `gorm:"type:double;DEFAULT 'false'"`
	IsAdmin       int    `gorm:"type:double;DEFAULT 'false'"`
}

/*
type Roles struct {
	gorm.Model
	UID   string `gorm:"type:varchar(60);UNIQUE;NOT NULL;"`
	Admin bool   `gorm:"type:double;DEFAULT 'false'"`
}


type Settings struct {
	gorm.Model
	Email_Confirmation bool `gorm:"type:double;DEFAULT 'false'"`
	Registration       bool `gorm:"type:double;DEFAULT 'false'"`
}
*/
/*
Timehash table structure:
Uid - users hash
email - users email
hash - 64 hash
param - Which function create this record. We need confirm email in signup and change current password
created - Where does record was created
livetime - hash life time
*/
type TimeHash struct {
	gorm.Model
	UID      string `gorm:"type:varchar(60);NOT NULL;"`
	Mail     string `gorm:"type:varchar(254);DEFAULT '';"`
	Hash     string `gorm:"type:varchar(254);DEFAULT '';"`
	Param    string `gorm:"type:varchar(254);DEFAULT '';"`
	Created  int    `gorm:"type:int(11);DEFAULT '0'"`
	Livetime int    `gorm:"type:int(11);DEFAULT '0'"`
}

func CheckDBStructure() {
	//Check DB and table config
	db, err := ConnectDB()
	if err != nil {
		SetErrorLog("dbstructure.go:81: " + err.Error())
		//return "error with db"
	}

	defer CloseConnection(db)
	/*
		if !db.Conn.HasTable(&Roles{}) {
			//Create roles table
			db.Conn.Set("gorm:table_options", "ENGINE=InnoDB;").CreateTable(&Roles{})
		}
	*/
	//Check if table users and roles exist
	if !db.Conn.HasTable(&Users{}) {
		SetErrorLog("dbstructure.go:94: " + "Table users do not exist. Create table Users")
		//db.Conn.Debug().AutoMigrate(&Users{})
		//Create users table
		db.Conn.Set("gorm:table_options", "ENGINE=InnoDB;").CreateTable(&Users{})

		//Add admin user
		//1. generate user hash
		userhash := Hashgen(8)
		//2. generate users Password
		userpass := RandomString(12)
		//2.1 generete pass passhash
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userpass), 8)
		if err != nil {

			SetErrorLog("dbstructure.go:108: " + err.Error())
		}

		//3. Admin User email
		useremail := viper.GetString("email.address")

		user := Users{
			UID:           userhash,
			Name:          "admin",
			Pass:          string(hashedPassword),
			Mail:          useremail,
			Mailsent:      int(time.Now().Unix()),
			Mailconfirmed: int(time.Now().Unix()),
			Created:       int(time.Now().Unix()),
			Status:        1,
			Completed:     1,
			IsAdmin:       1,
		}
		/*
			role := Roles{
				UID:   userhash,
				Admin: true,
			}
		*/
		db.Conn.Create(&user)
		//db.Conn.Create(&role)

		ans := fmt.Sprintf("Admin User created!\t\nname: admin\t\npass: %s\t\nemail: %s \t\n", userpass, useremail)
		fmt.Printf(ans)

	}

	//Create timehash table
	if !db.Conn.HasTable(&TimeHash{}) {
		db.Conn.Set("gorm:table_options", "ENGINE=InnoDB;").CreateTable(&TimeHash{})
	}
	/*
		if !db.Conn.HasTable(&Settings{}) {
			//Create settings table
			db.Conn.Set("gorm:table_options", "ENGINE=InnoDB;").CreateTable(&Settings{})
			setting := Settings{
				Email_Confirmation: false,
				Registration:       true,
			}
			db.Conn.Create(&setting)
		}
	*/
}

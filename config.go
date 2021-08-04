// Copyright 2019 Alexey Yanchenko <mail@yanchenko.me>
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
	"log"
	"os"
	"strings"

	viper "github.com/spf13/viper"
)

//CheckConfig() Check configuration file and stop app if config file mot foud or has errors
func CheckConfig() {
	viper.SetConfigName(configname) // name of config file (without extension)
	viper.AddConfigPath(Configpath)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			i := 0
			AskForConfigFile(i)
		} else {
			// Config file was found but another error was produced
			var m string = "Please Check config file. It is an error in it... \t"
			fmt.Printf(m)
			var ms string = "Server stop \t"
			fmt.Printf(ms)
			os.Exit(3)
		}
	}
	//Hash passwords
	HashConfigPasswords()
	SetLog("Gufo Starting. Config file OK")
}

// HashConfigPasswords - Change password in settings file to hash
func HashConfigPasswords() {
	/*
		viper.SetConfigName(configname) // name of config file (without extension)
		viper.AddConfigPath(configpath)
		err := viper.ReadInConfig() // Find and read the config file
		if err != nil {             // Handle errors reading the config file

			SetErrorLog("Please Check config file. It is an error in it...")
		}
	*/
	var pwd string = viper.GetString("database.password")
	//Hashed Password started with $2a##
	s := strings.Split(pwd, "a##")
	if s[0] != "$2" {
		// The password is not hashed

		// encrypt password
		encryptMsg, err := encrypt(pwd)
		if err != nil {
			SetErrorLog("config.go:71: " + err.Error())
		}
		//Create new passwrd recort
		var newpasswordrecord interface{} = "$2a##" + encryptMsg
		//s := make(interface{},  newpasswordrecord)
		viper.Set("database.password", newpasswordrecord)
		viper.WriteConfig()
	}

	var emailpwd string = viper.GetString("email.password")
	semail := strings.Split(emailpwd, "a##")
	if semail[0] != "$2" {
		// The password is not hashed

		// encrypt password
		encryptMsg, err := encrypt(emailpwd)
		if err != nil {
			SetErrorLog("config.go:88: " + err.Error())
		}
		//Create new passwrd recort
		var newpasswordrecord interface{} = "$2a##" + encryptMsg
		//s := make(interface{},  newpasswordrecord)
		viper.Set("email.password", newpasswordrecord)
		viper.WriteConfig()
	}

}

// DecryptConfigPasswords - Decrypt Hashed password in settings file to real password
func DecryptConfigPasswords(pwd string) string {
	s := strings.Split(pwd, "a##")
	if s[0] != "$2" {
		//The password is not hashed
		return pwd
	} else {
		//The password is hashed and need to be decrypted
		msg, _ := decrypt(s[1])
		return msg

	}
}

func EncryptConfigPassword(pwd string) (string, error) {
	encryptMsg, err := encrypt(pwd)
	if err != nil {
		return "", err
	}
	newpasswordrecord := "$2a##" + encryptMsg
	return newpasswordrecord, nil

}

func ConfigString(conf string) string {

	viper.SetConfigName(configname) // name of config file (without extension)
	viper.AddConfigPath(Configpath)
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		var errstring string = "Error to get " + conf + "from config file \t"
		SetErrorLog("config.go:120: " + errstring)
		return ""
	} else {
		var param string = viper.GetString(conf)
		return param
	}
}

func GetPass(conf string) string {
	p := viper.GetString(conf)
	pass := DecryptConfigPasswords(p)
	return pass
}

func AskForConfigFile(i int) {

	m := fmt.Sprintf("Config file was not found at %s\t\n\t\n", Configpath)
	fmt.Printf(m)
	fmt.Printf("Would you like to create it? [yes/no]: ")
	var ans string

	fmt.Scanln(&ans)
	switch ans {
	case "yes":
		CreateConfig()
	case "no":
		AnsNo()
	default:
		i = i + 1
		AnsDef(i)

	}

}

func AnsNo() {
	m := fmt.Sprintf("Please place \"%s\" file to the %s path and restart Gufo\t\n", configname, Configpath)
	fmt.Printf(m)
	fmt.Printf("Server stop \tn")
	os.Exit(3)
}

func AnsDef(i int) {
	if i == 3 {
		AnsNo()
	}
	AskForConfigFile(i)
}

func CreateConfig() {
	//fmt.Printf("Yes  \t\n Server stop \t\n")
	//os.Exit(3)
	fl := fmt.Sprintf("%s/settings.toml", Configpath)
	f, err := os.OpenFile(fl,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	viper.SetConfigName(configname)
	viper.AddConfigPath(Configpath)

	var ans string
	fmt.Printf("What's Gufo IP address? [127.0.0.1]: ")
	fmt.Scanln(&ans)
	if ans == "" {
		viper.Set("server.ip", "127.0.0.1")
	} else {
		viper.Set("server.ip", ans)
	}
	ans = ""

	fmt.Printf("Which port Gufo should listen? [8090]: ")
	fmt.Scanln(&ans)
	if ans == "" {
		viper.Set("server.port", "8090")
	} else {
		viper.Set("server.port", ans)
	}
	ans = ""

	fmt.Printf("Please set system directory (with slash)? [var/]: ")
	fmt.Scanln(&ans)
	sdir := "var/"
	if ans != "" {
		sdir = ans
	}
	viper.Set("server.sysdir", sdir)
	viper.Set("server.logdir", fmt.Sprintf("%slog/", sdir))
	viper.Set("server.plugindir", fmt.Sprintf("%slib/", sdir))
	viper.Set("server.langdir", fmt.Sprintf("%slang/", sdir))
	viper.Set("server.tempdir", fmt.Sprintf("%stemplates/", sdir))
	viper.Set("server.debug", true)
	ans = ""

	viper.Set("settings.registration", true)
	viper.Set("settings.email_confirmation", false)

	fmt.Printf("Please set system language? [english]: ")
	fmt.Scanln(&ans)
	if ans == "" {
		viper.Set("server.lang", "english")
	} else {
		viper.Set("server.lang", ans)
	}
	ans = ""

	fmt.Printf("What's DB type do you prefer? [mysql/postgres]: ")
	fmt.Scanln(&ans)
	viper.Set("database.type", ans)
	ans = ""

	fmt.Printf("What's DB protocol do you use? [tcp]: ")
	fmt.Scanln(&ans)
	if ans == "" {
		viper.Set("database.protocol", "tcp")
	} else {
		viper.Set("database.protocol", ans)
	}
	ans = ""

	fmt.Printf("What's DB host? [127.0.0.1]: ")
	fmt.Scanln(&ans)
	if ans == "" {
		viper.Set("database.host", "127.0.0.1")
	} else {
		viper.Set("database.host", ans)
	}
	ans = ""

	fmt.Printf("What's DB port? [3306/5432]: ")
	fmt.Scanln(&ans)
	viper.Set("database.port", ans)
	ans = ""

	fmt.Printf("DB username? [root]: ")
	fmt.Scanln(&ans)
	if ans == "" {
		viper.Set("database.user", "root")
	} else {
		viper.Set("database.user", ans)
	}
	ans = ""

	fmt.Printf("DB password? : ")
	fmt.Scanln(&ans)
	viper.Set("database.password", ans)
	ans = ""

	fmt.Printf("DB name? : ")
	fmt.Scanln(&ans)
	viper.Set("database.dbname", ans)
	ans = ""

	fmt.Printf("Redis settings? [redis://localhost]: ")
	fmt.Scanln(&ans)
	if ans == "" {
		viper.Set("redis.host", "redis://localhost")
	} else {
		viper.Set("redis.host", ans)
	}
	ans = ""

	fmt.Printf("Memcached host? [127.0.0.1]: ")
	fmt.Scanln(&ans)
	if ans == "" {
		viper.Set("memcached.host", "127.0.0.1")
	} else {
		viper.Set("memcached.host", ans)
	}
	ans = ""

	fmt.Printf("Memcached port? [11211]: ")
	fmt.Scanln(&ans)
	if ans == "" {
		viper.Set("memcached.port", "11211")
	} else {
		viper.Set("memcached.port", ans)
	}

	fmt.Printf("Would you like to setup email? [yes/no]: ")
	fmt.Scanln(&ans)
	if ans == "yes" {

		fmt.Printf("What's email adddress?: ")
		fmt.Scanln(&ans)
		viper.Set("email.address", ans)

		fmt.Printf("What's email host?: ")
		fmt.Scanln(&ans)
		viper.Set("email.host", ans)

		fmt.Printf("What's email port?: ")
		fmt.Scanln(&ans)
		viper.Set("email.port", ans)

		fmt.Printf("What's email username?: ")
		fmt.Scanln(&ans)
		viper.Set("email.user", ans)

		fmt.Printf("What's email password?: ")
		fmt.Scanln(&ans)
		viper.Set("email.password", ans)

		fmt.Printf("Reply-to: ")
		fmt.Scanln(&ans)
		viper.Set("email.reply", ans)

	} else {
		fmt.Printf("You can setup email later in config file \t\n")
	}

	fmt.Printf("Thank You! Config File created \t\n")
	CheckConfig()

}

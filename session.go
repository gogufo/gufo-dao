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
//

package gufodao

import (
	"fmt"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/getsentry/sentry-go"
	"github.com/gomodule/redigo/redis"
	"github.com/spf13/viper"
)

func SetSession(name string, isAdmin int, completed int, readonly int) (sessionToken string, exptime int, err error) {
	// Create a new random session token

	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	//exptime = time.Now().Add(time.Duration(viper.GetInt("token.expiretime"))).Unix()
	exptime = int(time.Now().Unix()) + viper.GetInt("token.expiretime")
	SetErrorLog(fmt.Sprintf("exptime: %v", exptime))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user":    name,
		"exipred": exptime,
	})

	// Sign and get the complete encoded token as a string using the secret
	sessionToken, err = token.SignedString([]byte(viper.GetString("token.secretKey")))

	// Set the token in the cache, along with the user whom it represents
	// The token has an expiry time of 120 seconds
	WriteTokenInRedis(sessionToken, name, isAdmin, completed, exptime, readonly)
	return sessionToken, exptime, err
}

func UpdateSession(sessionToken string) map[string]interface{} {

	tokenarray := strings.Split(sessionToken, " ")
	tokentype := tokenarray[0]
	token := tokenarray[1]

	//Get sesssion token
	ans := make(map[string]interface{})

	// Check Session in Redis
	n := ConfigString("redis.host")
	conn, err := redis.DialURL(n)
	if err != nil {
		SetErrorLog("session.go:59 " + err.Error())
	}

	response, err := redis.Values(conn.Do("HMGET", token, "expired", "uid", "isadmin", "completed", "readonly")) //commandName , ARG1, ARG2, ARG3
	if err != nil {
		// If there is an error in setting the cache, return an internal server error

		SetErrorLog("session.go:62: " + err.Error())
	}
	var exptime int
	var uid string
	var isadmin int
	var completed int
	var readonly int

	if _, err := redis.Scan(response, &exptime, &uid, &isadmin, &completed, &readonly); err != nil {
		// handle error
		SetErrorLog("session.go:70: " + err.Error())
	}

	if uid == "" {
		//Check Session in DB

		//Check DB and table config
		db, err := sf.ConnectDBv2()
		if err != nil {
			if viper.GetBool("server.sentry") {
				sentry.CaptureException(err)
			} else {
				sf.SetErrorLog(err.Error())
			}
			ans["httpcode"] = 500
			errormsg := []sf.ErrorMsg{}
			errorans := sf.ErrorMsg{
				Code:    "000027",
				Message: err.Error(),
			}
			errormsg = append(errormsg, errorans)
			return ans, errormsg, t
		}

		if tokentype == "APP" {
			tokentable := APITokens{}

			rows := db.Conn.Debug().Where(`token = ?`, token).First(&tokentable)

			if rows.RowsAffected == 0 || tokentable.Status == 0 || tokentable.Expiration < time.Now().Unix() {
				SetErrorLog("No uid")
				ans["error"] = "000011" // you are not authorised
				return ans
			}

			// Check Doues User is Admin in case of Token Admin Satatus
			if tokentable.IsAdmin == 1 {
				userExist := Users{}

				db.Conn.Debug().Where(`uid = ?`, tokentable.UID).First(&userExist)

				isadmin = userExist.IsAdmin
			}

			exptime = tokentable.Expiration
			uid = tokentable.UID
			completed = 1
			readonly = tokentable.Readonly

			//Write session into Redis
			WriteTokenInRedis(token, uid, isadmin, completed, exptime, readonly)

		} else {
			SetErrorLog("No uid")
			ans["error"] = "000011" // you are not authorised
			return ans
		}

	}

	//updates session
	newexptime := int(time.Now().Unix()) + viper.GetInt("token.expiretime")
	WriteTokenInRedis(token, uid, isadmin, completed, newexptime, readonly)

	ans["uid"] = uid
	ans["isadmin"] = isadmin
	ans["session_expired"] = newexptime
	ans["completed"] = completed
	ans["readonly"] = readonly
	ans["token"] = token
	ans["token_type"] = tokentype
	return ans

}

func DelSession(sessionToken string) {

	n := ConfigString("redis.host")
	conn, err := redis.DialURL(n)
	if err != nil {
		SetErrorLog("session.go:93: " + err.Error())
	}

	response, err := redis.Values(conn.Do("HMGET", sessionToken, "expired", "uid", "isadmin")) //commandName , ARG1, ARG2, ARG3
	if err != nil {
		// If there is an error in setting the cache, return an internal server error

		SetErrorLog("session.go:100: " + err.Error())
	}
	var exptime int64
	var uid string
	var isadmin int

	if _, err := redis.Scan(response, &exptime, &uid, &isadmin); err != nil {
		// handle error
		SetErrorLog("session.go:108: " + err.Error())
	}

	if uid == "" {

		return
	}

	_, err = cache.Do("DEL", sessionToken)
	if err != nil {

		return
	}
}

func WriteTokenInRedis(sessionToken string, uid string, isadmin int, completed int, exptime int, readonly int) {

	n := ConfigString("redis.host")
	conn, err := redis.DialURL(n)
	if err != nil {
		SetErrorLog("session.go:128: " + err.Error())
	}

	_, err = conn.Do("HMSET", sessionToken, "expired", exptime, "uid", uid, "isadmin", isadmin, "completed", completed, "readonly", readonly) //commandName , ARG1, ARG2, ARG3
	if err != nil {
		// If there is an error in setting the cache, return an internal server error

		SetErrorLog("session.go:137: " + err.Error())
	}

	_, err = conn.Do("EXPIRE", sessionToken, viper.GetInt("token.expiretime"))
	if err != nil {
		// If there is an error in setting the cache, return an internal server error

		SetErrorLog("session.go:146: " + err.Error())
	}

}

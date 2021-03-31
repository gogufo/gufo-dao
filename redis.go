// Copyright 2019 Alexey Yanchenko <mail@yanchenko.me>
//
// This file is part of the Neptune library.
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
	"github.com/gomodule/redigo/redis"
)

// Store the redis connection as a package level variable
var cache redis.Conn

func InitCache() {
	// Initialize the redis connection to a redis instance running on your local machine
	n := ConfigString("redis.host")
	conn, err := redis.DialURL(n)
	if err != nil {
		SetErrorLog("redis.go:31: " + err.Error())
	}
	// Assign the connection to the package level `cache` variable
	cache = conn
}

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
	"log"
	"os"
)

func GetLogDir() string {
	n := ConfigString("server.logdir")
	var logdir string = n
	return logdir
}

func SetLog(value string) {
	p := "gufo.log"
	WriteLog(value, p)
}

func SetErrorLog(value string) {
	p := "error.gufo.log"
	WriteLog(value, p)
}

func WriteLog(value string, p string) {
	n := GetLogDir()
	logdir := n + p
	f, err := os.OpenFile(logdir,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()

	logger := log.New(f, "", log.LstdFlags)
	logger.Println(value)

}

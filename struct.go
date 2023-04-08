// Copyright 2020 Alexey Yanchenko <mail@yanchenko.me>
//
// This file is part of the Gufo library.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package gufodao

// Request struct
type Request struct {
	Module     string
	Param      string
	ParamID    string
	Action     string
	Args       map[string]interface{} `json:"args"`
	Token      string
	TokenType  string
	TimeStamp  int    `json:"timestamp"`
	Language   string `json:"lang"`
	Dbversion  string
	UID        string
	IsAdmin    int
	SessionEnd int
	Completed  int
	Readonly   int
}

type ErrorMsg struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Error response struct
type ErrorResponse struct {
	Success   int                    `json:"success"`
	Error     []ErrorMsg             `json:"error"`
	Session   map[string]interface{} `json:"session"`
	TimeStamp int                    `json:"timestamp"`
	Language  string                 `json:"lang"`
	/*
		UID       string `json:"uid"`
		IsAdmin   string `json:"isadmin"`
		SesionExp int    `json:"sessionexp"`
	*/
}

// Succsess response struct
type SuccessResponse struct {
	Success   int                    `json:"success"`
	Data      map[string]interface{} `json:"data"`
	Session   map[string]interface{} `json:"session"`
	TimeStamp int                    `json:"timestamp"`
	Language  string                 `json:"lang"`
	/*
		UID       string `json:"uid"`
		IsAdmin   string `json:"isadmin"`
		SesionExp int    `json:"sessionexp"`
	*/
}

type Response struct {
	Data      map[string]interface{} `json:"data"`
	Session   map[string]interface{} `json:"session"`
	TimeStamp int                    `json:"timestamp"`
	Language  string                 `json:"lang"`
}

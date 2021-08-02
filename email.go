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
	"bytes"
	"crypto/tls"
	"html/template"
	"path"
	"strconv"
	"strings"

	"github.com/spf13/viper"
	"gopkg.in/gomail.v2"
)

//Request struct
type MailRequest struct {
	from    string
	to      []string
	subject string
	body    string
	file    []string
}

type MailSettings struct {
	Custom  bool
	Host    string
	Port    string
	User    string
	Pass    string
	Address string
	Reply   string
	Title   string
}

func NewRequest(to []string, subject, body string, attach []string) *MailRequest {
	return &MailRequest{
		to:      to,
		subject: subject,
		body:    body,
		file:    attach,
	}
}

func (r *MailRequest) SendEmail(ms *MailSettings) (bool, error) {

	host := viper.GetString("email.host")
	port := viper.GetString("email.port")
	user := viper.GetString("email.user")
	pass := GetPass("email.password")
	address := viper.GetString("email.address")
	reply := viper.GetString("email.reply")
	fromuser := viper.GetString("email.title")

	if ms.Custom {
		host = ms.Host
		port = ms.Port
		user = ms.User
		pass = ms.Pass
		address = ms.Address
		reply = ms.Reply
		fromuser = ms.Title
	}

	//xm := "X-Mailer: gufo; \r\n"
	//mime := "MIME-version: 1.0;\r\nContent-Type: text/html;charset=\"UTF-8\";\r\n\r\n"

	// creates  type Message
	mail := gomail.NewMessage(gomail.SetCharset("UTF-8"))

	// set the header and body parameters for a base message
	toString := strings.Join(r.to, ",")
	mail.SetAddressHeader("From", address, fromuser)
	mail.SetHeader("To", toString)
	mail.SetAddressHeader("Reply-To", reply, fromuser)
	mail.SetHeader("MIME-version", "1.0")
	mail.SetHeader("X-Mailer", "gufo")
	mail.SetHeader("Subject", r.subject)
	mail.SetBody("text/html", r.body)

	if len(r.file) > 0 {
		for i := 0; i < len(r.file); i++ {
			mail.Attach(r.file[i])
		}

	}

	portint, _ := strconv.Atoi(port)
	// Create an smtp dialer
	sendMail := gomail.NewDialer(host, portint, user, pass)

	// InsecureSkipVerify: true will ignore confirmation requests of using the source email in the smtp service.
	sendMail.TLSConfig = &tls.Config{InsecureSkipVerify: viper.GetBool("email.InsecureSkipVerify")}

	if err := sendMail.DialAndSend(mail); err != nil {
		SetErrorLog("email.go.81 Send email error: " + err.Error())
		return false, err
	} else {
		SetLog("Email to " + toString + " was sent")
	}

	return true, nil
}

func (r *MailRequest) ParseTemplate(templateFileName string, data interface{}) error {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return err

	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return err

	}
	r.body = buf.String()
	return nil
}

func SendHTMLEmail(to string, title string, link []string, subject string, templ string, attach []string, ms *MailSettings) {

	htmllink := []template.HTML{}

	for i := 0; i < len(link); i++ {

		htmllink = append(htmllink, template.HTML(link[i]))

	}

	templateData := struct {
		Title      string
		Paragraphs []template.HTML
	}{
		Title:      title,
		Paragraphs: htmllink,
	}

	templateDir := viper.GetString("server.tempdir")
	var emailtemplate = path.Join(templateDir, templ)

	r := NewRequest([]string{to}, subject, "", attach)

	if err := r.ParseTemplate(emailtemplate, templateData); err == nil {
		r.SendEmail(ms)
		//fmt.Println(ok)

	} else {
		SetErrorLog("email.go:143: " + err.Error())
	}
}

package models

import (
	"net/smtp"
)

const (
	username = "koteko.oficial@gmail.com"
	password = "Contra123"
	host     = "smtp.gmail.com"
	addr     = "smtp.gmail.com:587"
)

func SendEmail(to []string, msg []byte) (err error) {
	to = []string{"19030624@itcelaya.edu.mx"}
	msg = []byte("hola, no es spam")

	err = smtp.SendMail(addr, smtp.PlainAuth("", username, password, host), username, to, msg)
	if err != nil {
		return
	}

	return
}

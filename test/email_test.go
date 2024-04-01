package test

import (
	"github.com/jordan-wright/email"
	"net/smtp"
	"testing"
)

func TestSendEmail(t *testing.T) {
	e := email.NewEmail()
	e.From = "Crayon <2993373191@qq.com>"
	e.To = []string{"sudolsq@gmail.com"}
	e.Subject = "验证发送测试"

	e.HTML = []byte("<h1>Fancy HTML12311231 is supported, too!</h1>")
	err := e.Send("smtp.qq.com:587", smtp.PlainAuth("", "2993373191@qq.com", "vlkrwkqjayqedehc", "smtp.qq.com"))

	//返回EOF时，关闭SLL重试
	//err := e.SendWithTLS("smtp.qq.com:587",
	//	smtp.PlainAuth("", "2993373191@qq.com",
	//		"vlkrwkqjayqedehc", "smtp.qq.com"),
	//	&tls.Config{InsecureSkipVerify: true, ServerName: "smtp.qq.com"})

	if err != nil {
		t.Fatal(err)
	}
}

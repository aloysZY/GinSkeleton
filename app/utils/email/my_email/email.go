// Package my_email 主要是避免了variable的循环引用
package my_email

import (
	"crypto/tls"

	"gopkg.in/gomail.v2"
)

type Email struct {
	Host     string
	Port     int
	UserName string
	Password string
	From     string
	IsSSL    bool
}

// SendMail 发送邮件。就一个方法不放在一个单独的包里面了
func (e *Email) SendMail(to []string, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", e.From)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject) // 邮件主题
	m.SetBody("text/html", body)    // 邮件内容

	dialer := gomail.NewDialer(e.Host, e.Port, e.UserName, e.Password)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: e.IsSSL}
	return dialer.DialAndSend(m)
}

package pkg

import (
	"fmt"
	"gopkg.in/gomail.v2"
)

type Email struct {
	From         string
	To           []string
	SmtpHost     string
	SmtpPort     int
	SmtpUsername string
	SmtpPassword string
}

type Options func(*Email)

func WithSmtpHost(smtpHost string) Options {
	return func(e *Email) {
		e.SmtpHost = smtpHost
	}
}

func WithSmtpPort(smtpPort int) Options {
	return func(e *Email) {
		e.SmtpPort = smtpPort
	}
}

func WithSmtpUsername(smtpUsername string) Options {
	return func(e *Email) {
		e.SmtpUsername = smtpUsername
	}
}

func WithSmtpPassword(smtpPassword string) Options {
	return func(e *Email) {
		e.SmtpPassword = smtpPassword
	}
}

func WithTo(to []string) Options {
	return func(e *Email) {
		e.To = to
	}
}

func WithFrom(from string) Options {
	return func(e *Email) {
		e.From = from
	}
}

func NewEmailSMTP(option ...Options) *Email {
	em := Email{}
	for _, opt := range option {
		opt(&em)
	}
	return &em
}

func (e *Email) SendEmailSMTP(subject, message string) error {
	m := gomail.NewMessage()
	//发送人
	m.SetHeader("From", e.From)
	//接收人
	m.SetHeader("To", e.To...)
	//抄送人
	//m.SetAddressHeader("Cc", "xxx@qq.com", "xiaozhujiao")
	//主题
	m.SetHeader("Subject", subject)
	//内容
	m.SetBody("text/html", message)
	//附件
	//m.Attach("./myIpPic.png")

	//拿到token，并进行连接,第4个参数是填授权码
	d := gomail.NewDialer(e.SmtpHost, e.SmtpPort, e.SmtpUsername, e.SmtpPassword)

	// 发送邮件
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	fmt.Printf("send mail success\n")
	return nil
}

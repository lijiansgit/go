package libs

import (
	"fmt"
	"net"
	"net/mail"
	"net/smtp"
	"crypto/tls"
	"github.com/astaxie/beego"
)

type Mail struct {
	From		mail.Address
	To			mail.Address
	ServerName	string
	Password	string
	Subject		string
	Messages	string
}

func NewMail() *Mail {
	fromName := beego.AppConfig.String("FromName")
	fromAddr := beego.AppConfig.String("FromAddr")
	toName := beego.AppConfig.String("ToName")
	toAddr := beego.AppConfig.String("ToAddr")
	serverName := beego.AppConfig.String("ServerName")
	password := beego.AppConfig.String("Password")
	subject := beego.AppConfig.String("Subject")
	headers := make(map[string]string)
	from := mail.Address{fromName, fromAddr}
	to := mail.Address{toName, toAddr}
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = subject
	headers["Content-Type"] = "text/html; charset=UTF-8"
	messages := ""
	for k,v := range headers {
		messages += fmt.Sprintf("%s: %s\r\n", k, v)
	}

	return &Mail{
		From:       from,
		To:         to,
		ServerName: serverName,
		Password:   password,
		Subject:    subject,
		Messages:   messages,
	}
}

func (m *Mail) Send(body string) error  {
	messages := m.Messages + "\r\n" + body
	host, _, _ := net.SplitHostPort(m.ServerName)
	auth := smtp.PlainAuth("",m.From.Address, m.Password, host)

	// TLS config
	tlsConfig := &tls.Config {
		InsecureSkipVerify: true,
		ServerName: host,
	}

	conn, err := tls.Dial("tcp", m.ServerName, tlsConfig)
	if err != nil {
		return err
	}

	c, err := smtp.NewClient(conn, host)
	if err != nil {
		return err
	}

	// Auth
	if err = c.Auth(auth); err != nil {
		return err
	}

	// To && From
	if err = c.Mail(m.From.Address); err != nil {
		return err
	}

	if err = c.Rcpt(m.To.Address); err != nil {
		return err
	}

	// Data
	w, err := c.Data()
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(messages))
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	c.Quit()
	return nil
}










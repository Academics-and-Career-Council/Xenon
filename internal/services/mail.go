package services

import (
	"log"
	"strings"

	"github.com/spf13/viper"
	"gopkg.in/gomail.v2"
)

func SendMail(subject, body string, to []string) error {

	m := gomail.NewMessage()
	m.SetHeader("From", viper.GetString("smtp.mail"))
	m.SetHeader("To", strings.Join(to, ","))
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(viper.GetString("smtp.host"), viper.GetInt("smtp.port"), viper.GetString("smtp.user"), viper.GetString("smtp.pwd"))
	if err := d.DialAndSend(m); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

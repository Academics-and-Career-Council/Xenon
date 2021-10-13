package services

import (
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

func VerifyCaptcha(token string) error {
	endpoint := "https://hcaptcha.com/siteverify"
	data := url.Values{}
	data.Set("response", token)
	data.Set("secret", viper.GetString("hcaptcha.secret"))

	client := &http.Client{}
	r, err := http.NewRequest("POST", endpoint, strings.NewReader(data.Encode())) // URL-encoded payload
	if err != nil {
		log.Fatal(err)
	}
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	_, err = client.Do(r)
	return err
}

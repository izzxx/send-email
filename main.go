package main

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"gopkg.in/gomail.v2"
)

var (
	FromEmail      string
	RecipientEmail string
	Password       string
)

const (
	Host = "smtp.gmail.com"
	Port = 587
)

func new() {
	FromEmail = os.Getenv("FromEmail")
	RecipientEmail = os.Getenv("RecipientEmail")
	Password = os.Getenv("Password")
}

func Setup() error {
	data, err := ioutil.ReadFile(".env")
	if err != nil {
		return err
	}

	datas := strings.Split(string(data), "\n")
	for _, env := range datas {
		split := strings.Split(env, "=")
		if len(split) > 2 {
			return errors.New("invalid data")
		}

		err = os.Setenv(strings.TrimSpace(split[0]), strings.TrimSpace(split[1]))
		if err != nil {
			return err
		}
	}

	// Environtment variable
	new()

	return nil
}

func main() {
	// Setup environtment variable
	err := Setup()
	if err != nil {
		log.Fatal(err)
	}

	// Create new message
	message := gomail.NewMessage()

	// Sets a value to the given header field.
	message.SetHeader("From", FromEmail)
	message.SetHeader("To", RecipientEmail)
	message.SetHeader("Subject", "Send Email")
	message.SetAddressHeader("Cc", FromEmail, "Send-Email")

	//Body of the message
	message.SetBody("text/plain", "Message here")

	// Attaches the files to the email.
	message.Attach("./image.png", gomail.Rename("image"))

	dial := gomail.Dialer{
		Host:     Host,
		Port:     Port,
		Username: FromEmail,
		Password: Password,
	}

	log.Println("Start sending email...")

	err = dial.DialAndSend(message)
	if err != nil {
		log.Panic(err)
	}

	// Shorter way
	// if err = gomail.NewDialer(Host, Port, FromEmail, Password).DialAndSend(message); err != nil {
	// 	log.Fatal(err)
	// }

	log.Println("Email has been sent.")
}

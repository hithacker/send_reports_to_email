package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/mail"
	"net/smtp"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/scorredoira/email"
)

type config struct {
	FromEmailID string
	FromName    string
	Password    string
	ToEmailIds  []string
}

func main() {

	var conf config
	if _, err := toml.DecodeFile("config.toml", &conf); err != nil {
		panic(err)
	}

	fmt.Println(conf.FromEmailID)
	fmt.Println(conf.FromName)
	fmt.Println(conf.Password)
	fmt.Println(conf.ToEmailIds)

	todayDate := fmt.Sprintf("%d-%d-%d", time.Now().Day(), time.Now().Month(), time.Now().Year())

	m := email.NewMessage("Lab Report "+todayDate, "Lab Report")
	m.From = mail.Address{Name: conf.FromName, Address: conf.FromEmailID}
	m.To = conf.ToEmailIds

	files, err := ioutil.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		fileDate := fmt.Sprintf("%d-%d-%d", file.ModTime().Day(), file.ModTime().Month(), file.ModTime().Year())
		if filepath.Ext(file.Name()) == ".pdf" && fileDate == todayDate {
			err = m.Attach(file.Name())
			if err != nil {
				panic(err)
			}
		}
	}

	if err != nil {
		log.Println(err)
	}

	err = email.Send("smtp.gmail.com:587", smtp.PlainAuth("", conf.FromEmailID, conf.Password, "smtp.gmail.com"), m)
	log.Println(err)
}

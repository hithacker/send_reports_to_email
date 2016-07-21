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
	FileTypes   []string
}

func main() {

	var conf config
	if _, err := toml.DecodeFile("config.toml", &conf); err != nil {
		panic(err)
	}

	todayDate := fmt.Sprintf("%d-%d-%d", time.Now().Day(), time.Now().Month(), time.Now().Year())

	m := email.NewMessage("Lab Report "+todayDate, "Lab Report")
	m.From = mail.Address{Name: conf.FromName, Address: conf.FromEmailID}
	m.To = conf.ToEmailIds

	files, err := ioutil.ReadDir(".")
	if err != nil {
		panic(err)
	}

	numberOfAttachedFiles := 0

	for _, file := range files {
		fileDate := fmt.Sprintf("%d-%d-%d", file.ModTime().Day(), file.ModTime().Month(), file.ModTime().Year())

		validExtension := false
		for _, fileType := range conf.FileTypes {
			if fileType == filepath.Ext(file.Name()) {
				validExtension = true
				break
			}
		}

		if validExtension && fileDate == todayDate {
			err = m.Attach(file.Name())
			numberOfAttachedFiles++
			if err != nil {
				panic(err)
			}
		}
	}

	fmt.Printf("Going to send %d fies\n", numberOfAttachedFiles)

	err = email.Send("smtp.gmail.com:587", smtp.PlainAuth("", conf.FromEmailID, conf.Password, "smtp.gmail.com"), m)
	if err != nil {
		log.Println(err.Error())
	}
}

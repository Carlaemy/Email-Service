package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/mail"
	"net/smtp"
	"strconv"
	"time"

	"./Model"
)

func main() {

	go RecordFiles()
	RecordUsers()
}

func SendMail() {
	from := mail.Address{"Sent from Go", "fbfranco29@gmail.com"}
	to := mail.Address{"Bismarck Franco", "bfranco@info-arch.com"}
	subject := "Action on files:"

	recip := model.Recipient{Name: to.Address, FileName: "Auto.png", FileID: "1", Action: "Guardó"}

	users := GetUsers()

	for _, user := range users {

		to = mail.Address{user.Name + " " + user.LastName, user.Email}

		headers := make(map[string]string)
		headers["From"] = from.String()
		headers["To"] = to.String()
		headers["Subject"] = subject
		headers["Content-type"] = `text/html; charset="UTF-8"`

		message := ""

		for i, data := range headers {
			message += fmt.Sprintf("%s: %s\r\n", i, data)
		}

		t, err := template.ParseFiles("./Template/Template.html")
		Error(err)

		buf := new(bytes.Buffer)
		err = t.Execute(buf, recip)
		Error(err)

		message += buf.String()

		ServerName := "smtp.gmail.com:465"
		Host := "smtp.gmail.com"

		auth := smtp.PlainAuth("", "fbfranco29@gmail.com", "70287468", Host)

		tlsConfig := &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         Host,
		}

		conn, err := tls.Dial("tcp", ServerName, tlsConfig)
		Error(err)

		client, err := smtp.NewClient(conn, Host)
		Error(err)

		err = client.Auth(auth)
		Error(err)

		err = client.Mail(from.Address)
		Error(err)

		err = client.Rcpt(to.Address)
		Error(err)

		w, err := client.Data()
		Error(err)

		_, err = w.Write([]byte(message))
		Error(err)

		err = w.Close()
		Error(err)
		client.Quit()
		fmt.Println("Email sent " + user.Email)
	}
}

func Error(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func FilesCant() int {
	path := "files"
	url := fmt.Sprintf("http://localhost:3000/%s", path)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
	}
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
	}
	defer resp.Body.Close()

	var record []model.File
	if err := json.NewDecoder(resp.Body).Decode(&record); err != nil {
		log.Println(err)
	}
	cant := len(record)
	return cant
}

func UserCant() int {
	path := "users"
	url := fmt.Sprintf("http://localhost:3000/%s", path)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
	}
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
	}
	defer resp.Body.Close()

	var record []model.User
	if err := json.NewDecoder(resp.Body).Decode(&record); err != nil {
		log.Println(err)
	}
	cant := len(record)
	return cant
}

func RecordFiles() {
	fmt.Println(strconv.Itoa(FilesCant()) + " existing files.")
	y := 0

	for {
		x := FilesCant()
		time.Sleep(1 * time.Second)
		if x != y {
			y = x
			fmt.Println(strconv.Itoa(FilesCant()) + " existing files.")
			SendMail()
		}
	}
}

func RecordUsers() {
	fmt.Println(strconv.Itoa(UserCant()) + " existing users.")
	y := 0

	for {
		x := UserCant()
		time.Sleep(1 * time.Second)
		if x != y {
			y = x
			fmt.Println(strconv.Itoa(UserCant()) + " existing users.")
		}
	}
}

func GetUsers() []model.User {
	path := "users"
	url := fmt.Sprintf("http://localhost:3000/%s", path)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
	}
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
	}
	defer resp.Body.Close()

	var record []model.User
	if err := json.NewDecoder(resp.Body).Decode(&record); err != nil {
		log.Println(err)
	}
	return record
}

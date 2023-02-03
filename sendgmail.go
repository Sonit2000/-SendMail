package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
)

type Page struct {
	Title, Subject, Body string
}
type loginAuth struct {
	username, password string
}

func LoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}
func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte(a.username), nil
}
func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username":
			return []byte(a.username), nil
		case "Password":
			return []byte(a.password), nil
		}
	}
	return nil, nil
}
func SendToMail(addr string, auth smtp.Auth, from string, to []string, msg []byte) (err error) {
	c, err := smtp.Dial(addr)
	if err != nil {
		log.Println("Create smtp client error", err)
		return err
	}
	defer c.Close()
	if auth != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(auth); err != nil {
				log.Println("Error during AUTH", err)
				return err
			}
		}
	}
	if err = c.Mail(from); err != nil {
		return err
	}
	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}

	w, err := c.Data()
	if err != nil {
		return err
	}

	_, err = w.Write(msg)
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	return c.Quit()
}
func main() {
	tmpl, err := template.ParseFiles("template.html")
	if err != nil {
		log.Fatal("Failed to parse template: %s", err)
	}

	data := struct {
		Title   string
		Subject string
		Body    string
	}{
		Title:   "Bạn thân mến!",
		Subject: "Vậy là thời gian trôi nhanh quá, bây giờ đã là thời điểm cuối năm – đây cũng là thời gian để chúng ta cùng nhìn lại những kỉ niệm vui buồn của năm. Sắp bước sang năm mới, tôi chúc bạn cùng gia đình luôn tràn ngập tiếng cười, vui vẻ, đầm ấm, sum vầy. Chúng tôi hi vọng, qua dịp Tết này sẽ là thời điểm đánh dấu một bước ngoặt mới cho sự thành công trong công việc mà chúng ta đang phấn đấu.",
		Body:    "Chúc mừng năm mới và mọi điều tốt đẹp cho năm mới 2023.",
	}
	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	var msg bytes.Buffer
	subject := "Subject: Chúc Mừng Năm Mới\n%s\n\n"
	msg.Write([]byte(fmt.Sprintf(subject, mimeHeaders)))

	if err := tmpl.Execute(&msg, data); err != nil {
		log.Fatalf("Failed to render template: %s", err)
	}

	auth := smtp.PlainAuth("", "hoangson.it2000@gmail.com", "odbqgtscznsgmhxa", "smtp.gmail.com")
	to := []string{"nguyenhoangson221478681@gmail.com", "son.nguyen@alttekglobal.com"}
	/*msg := []byte("To:nguyenhoangson221478681@gmail.com\r\n" + "Subject: discount Gophers!\r\n" + "\r\n" + "This is the email body.\r\n")*/
	err = smtp.SendMail("smtp.gmail.com:587", auth, "hoangson.it2000@gmail.com", to, msg.Bytes())
	if err != nil {
		log.Fatal("Send Fail: ", err)
	}
	log.Println("Send Success!")
}

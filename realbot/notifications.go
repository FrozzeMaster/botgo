package realbot

import (
	"bytes"
	"fmt"
	"net/smtp"
	"text/template"
)

type Notifier struct {
	Msg       string
	MsgEmail  string
	MsgNumber int
}

func NewNotifier() *Notifier {
	return &Notifier{
		Msg:       "",
		MsgNumber: 0,
	}
}
func (not *Notifier) AddMsg(msg ...interface{}) {
	not.Msg += fmt.Sprintln(msg...)
	not.MsgEmail += fmt.Sprintln(msg...)
	not.MsgEmail += "<br>"
	not.MsgNumber++
}
func (not *Notifier) Clear() {
	not.Msg = ""
	not.MsgEmail = ""
	not.MsgNumber = 0
}
func (not *Notifier) SendEmail() {
	from := "botbocislawbot@gmail.com"
	password := "Travian1997"

	// Receiver email address.
	to := []string{
		"michalswnsk@gmail.com",
		"jasgrzegorek@gmail.com",
		"botbocislawbot@gmail.com",
	}

	// smtp server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	auth := smtp.PlainAuth("", from, password, smtpHost)

	t, err := template.ParseFiles("realbot/htmltemplates/mail.html")
	if err != nil {
		fmt.Println(err)
		return
	}

	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: Logi z bota \n%s\n\n", mimeHeaders)))

	t.Execute(&body, struct {
		Msg string
	}{
		Msg: not.MsgEmail,
	})

	// Sending email.
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, body.Bytes())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Wysłano powiadomienie mailem!")
}
func (not *Notifier) PrintAllOut() {
	fmt.Println(not.Msg)
}

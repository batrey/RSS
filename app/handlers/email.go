package handlers

import (
	"net/smtp"
	"os"
	"strings"
	"time"
)

func Email(addr string, Emailreceiver []string, emailSender string, content []byte) error {
	host := os.Getenv("MAIL_HOST")
	port := os.Getenv("MAIL_PORT")
	addr = host + ":" + port

	to := strings.Join(Emailreceiver, ",")
	pw := os.Getenv("MAIL_PASSWORD")
	auth := smtp.PlainAuth("", emailSender, pw, host)

	subject := "RSS Feeds for " + time.Now().Format("Jan 02, 2006")

	msg := strings.Builder{}
	msg.WriteString("From: \"Feed Update\" <" + emailSender + ">\n")
	msg.WriteString("To: " + to + "\n")
	msg.WriteString("Subject: " + subject + "\n")
	msg.WriteString("MIME-version: 1.0;\n")
	msg.WriteString("Content-Type: text/html;charset=\"UTF-8\";\n")
	content = append(content, []byte(msg.String())...)
	err := smtp.SendMail(addr, auth, emailSender, Emailreceiver, content)
	if err != nil {
		return err
	}
	return nil
}

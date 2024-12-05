package objects

import (
	"fmt"
	"log"
	"net/smtp"
	"strings"
)

type Filter interface {
	Apply(msg Message) *Message
}

type WordsFilter struct{}

func (f *WordsFilter) Apply(msg Message) *Message {
	stopWords := []string{"bird-watching", "ailurophobia", "mango"}
	if !containsStopWords(msg.Text, stopWords) {
		return &msg
	}
	return nil
}

func containsStopWords(text string, stopWords []string) bool {
	for _, word := range stopWords {
		if strings.Contains(text, word) {
			return true
		}
	}
	return false
}

type ScreamingFilter struct{}

func (f *ScreamingFilter) Apply(msg Message) *Message {
	msg.Text = strings.ToUpper(msg.Text)
	return &msg
}

type PublishFilter struct{}

func (f *PublishFilter) Apply(msg Message) *Message {
	sendEmail(msg)
	return &msg
}

func sendEmail(msg Message) {
	config := LoadConfig()
	auth := smtp.PlainAuth("", config.From, config.Password, config.SmtpHost)
	body := fmt.Sprintf("From user: %s\nMessage: %s", msg.Alias, msg.Text)
	subject := "Subject: SA Course\n"
	fromHeader := "From: " + config.From + "\n"
	toHeader := "To: " + strings.Join(config.Receivers, ", ") + "\n"

	message := []byte(subject + fromHeader + toHeader + "\n" + body)

	err := smtp.SendMail(config.SmtpHost+":"+config.SmtpPort, auth, config.From, config.Receivers, message)
	if err != nil {
		log.Printf("Failed to send email: %s", err)
		return
	}
	log.Println("Email sent successfully")
}

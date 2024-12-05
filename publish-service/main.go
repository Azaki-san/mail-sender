package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strings"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	RECEIVERS = []string{"leonardozakarus@yandex.ru"}
	rabbitURL string
	from      string
	password  string
	smtpHost  string
	smtpPort  string
)

func loadEnvVars() {
	rabbitURL = os.Getenv("RABBITMQ_URL")
	from = os.Getenv("EMAIL_FROM")
	password = os.Getenv("EMAIL_PASSWORD")
	smtpHost = os.Getenv("SMTP_HOST")
	smtpPort = os.Getenv("SMTP_PORT")

	if rabbitURL == "" || from == "" || password == "" || smtpHost == "" || smtpPort == "" {
		log.Fatal("Error: One or more required environment variables are missing")
	}
}

type Message struct {
	Text  string `json:"text"`
	Alias string `json:"alias"`
}

func main() {
	loadEnvVars()
	conn, err := amqp.Dial(rabbitURL)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	msgs, err := ch.Consume("screamed-messages", "", true, false, false, false, nil)
	failOnError(err, "Failed to register a consumer")

	log.Println("Publish Service running")

	for d := range msgs {
		var msg Message
		json.Unmarshal(d.Body, &msg)
		sendEmail(msg)
		log.Printf("Sent email for message: %s", msg.Text)
	}
}

func sendEmail(msg Message) {
	auth := smtp.PlainAuth("", from, password, smtpHost)
	body := fmt.Sprintf("From user: %s\nMessage: %s", msg.Alias, msg.Text)
	subject := "Subject: SA Course\n"
	fromHeader := "From: " + from + "\n"
	toHeader := "To: " + strings.Join(RECEIVERS, ", ") + "\n"

	message := []byte(subject + fromHeader + toHeader + "\n" + body)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, RECEIVERS, message)
	if err != nil {
		log.Printf("Failed to send email: %s", err)
		return
	}
	log.Println("Email sent successfully")
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

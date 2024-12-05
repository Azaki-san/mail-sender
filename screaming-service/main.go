package main

import (
	"encoding/json"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strings"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Message struct {
	Text  string `json:"text"`
	Alias string `json:"alias"`
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		if os.Getenv("RABBITMQ_URL") == "" {
			log.Fatal("empty environment")
		}
	}
	rabbitURL := os.Getenv("RABBITMQ_URL")
	conn, err := amqp.Dial(rabbitURL)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	msgs, err := ch.Consume("filtered-messages", "", true, false, false, false, nil)
	failOnError(err, "Failed to register a consumer")

	q, err := ch.QueueDeclare("screamed-messages", false, false, false, false, nil)
	failOnError(err, "Failed to declare a queue")

	log.Println("SCREAMING Service running")

	for d := range msgs {
		var msg Message
		json.Unmarshal(d.Body, &msg)
		msg.Text = strings.ToUpper(msg.Text)

		body, _ := json.Marshal(msg)
		ch.Publish("", q.Name, false, false, amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
		log.Printf("Processed message: %s", msg.Text)
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

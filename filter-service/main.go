package main

import (
	"encoding/json"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strings"

	amqp "github.com/rabbitmq/amqp091-go"
)

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

	msgs, err := ch.Consume("incoming-messages", "", true, false, false, false, nil)
	failOnError(err, "Failed to register a consumer")

	q, err := ch.QueueDeclare("filtered-messages", false, false, false, false, nil)
	failOnError(err, "Failed to declare a queue")

	stopWords := []string{"bird-watching", "ailurophobia", "mango"}
	log.Println("Filter Service running")

	for d := range msgs {
		var msg Message
		json.Unmarshal(d.Body, &msg)
		if containsStopWords(msg.Text, stopWords) {
			continue
		}
		body, _ := json.Marshal(msg)
		ch.Publish("", q.Name, false, false, amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	}
}

func containsStopWords(text string, stopWords []string) bool {
	for _, word := range stopWords {
		if strings.Contains(text, word) {
			return true
		}
	}
	return false
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

type Message struct {
	Text  string `json:"text"`
	Alias string `json:"alias"`
}

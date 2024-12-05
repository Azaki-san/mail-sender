package main

import (
	"encoding/json"
	"example.com/mod/objects"
	"example.com/mod/objects/pipeline"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"net/http"
	"os"
)

var (
	RABBIT_URL string
	RUN_MODE   string
	PORT       string
	PIPELINE   *pipeline.Pipeline

	from     string
	password string
	smtpHost string
	smtpPort string

	ch *amqp.Channel
	q  amqp.Queue
)

func loadEnvVars() {
	RABBIT_URL = os.Getenv("RABBITMQ_URL")
	RUN_MODE = os.Getenv("RUN_MODE")
	PORT = os.Getenv("APP_PORT")

	if RABBIT_URL == "" || RUN_MODE == "" || PORT == "" {
		log.Fatal("Error: One or more required environment variables are missing")
	}

	if RUN_MODE == "pipes" {
		from = os.Getenv("EMAIL_FROM")
		password = os.Getenv("EMAIL_PASSWORD")
		smtpHost = os.Getenv("SMTP_HOST")
		smtpPort = os.Getenv("SMTP_PORT")
		if from == "" || password == "" || smtpHost == "" || smtpPort == "" {
			log.Fatal("Error RUN_MODE pipes: One or more required environment variables are missing")
		}
		PIPELINE = initPipeline()
	} else {
		conn, err := amqp.Dial(RABBIT_URL)
		failOnError(err, "Failed to connect to RabbitMQ")

		ch, err = conn.Channel()
		failOnError(err, "Failed to open a channel")
		failOnError(err, "Failed to declare a queue")

		q, err = ch.QueueDeclare(
			"incoming-messages", false, false, false, false, nil,
		)
	}
}

func main() {
	loadEnvVars()

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(MetricsMiddleware)

	r.Post("/message", func(w http.ResponseWriter, r *http.Request) {
		var msg objects.Message
		err := json.NewDecoder(r.Body).Decode(&msg)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		switch RUN_MODE {
		case "pipes":
			PIPELINE.Process(msg)
		case "events":
			body, _ := json.Marshal(msg)
			err = ch.Publish("", q.Name, false, false, amqp.Publishing{
				ContentType: "application/json",
				Body:        body,
			})
		}

		failOnError(err, "Failed to publish a message")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "Message sent"})
	})

	r.Handle("/metrics", promhttp.Handler())

	defer ch.Close()
	log.Println(fmt.Sprintf("REST API Service running on :%s RUN_MODE %s", PORT, RUN_MODE))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", PORT), r))
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

package objects

import (
	"os"
)

type Config struct {
	From      string
	Password  string
	SmtpHost  string
	SmtpPort  string
	Receivers []string
}

func LoadConfig() *Config {
	return &Config{
		From:     os.Getenv("EMAIL_FROM"),
		Password: os.Getenv("EMAIL_PASSWORD"),
		SmtpHost: os.Getenv("SMTP_HOST"),
		SmtpPort: os.Getenv("SMTP_PORT"),
		//Receivers: []string{"a.alexeev@innopolis.university", "m.gladyshev@innopolis.university", "n.kuzmin@innopolis.university"},
		Receivers: []string{"leonardozakarus@yandex.ru"},
	}
}

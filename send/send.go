package main

import (
	"context"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/spf13/viper"
)

type Config struct {
	RABBITMQ_URL string `mapstructure:"RABBITMQ_URL"`
}

var C *Config

func loadConfig() {
	v := viper.New()
	v.AddConfigPath(".")
	v.SetConfigName(".env")
	v.SetConfigType("env")
	err := v.ReadInConfig()
	if err != nil {
		failOnError(err, "Failed to read config")
	}
	v.AutomaticEnv()

	err = v.Unmarshal(&C)
	if err != nil {
		failOnError(err, "Failed to read enivroment")
	}
}
func main() {
	loadConfig()
	conn, err := amqp.Dial(C.RABBITMQ_URL)
	failOnError(err, "Failed to connect to RabbitMQ")

	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare("hello", false, false, false, false, nil)
	failOnError(err, "Failed to declare a queue")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	body := "Hello World!"
	err = ch.PublishWithContext(
		ctx,
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)
	failOnError(err, "Failed to publish a message")
	log.Printf("[x] Sent %s\n", body)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

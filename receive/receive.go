package main

import (
	"log"

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
	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "failed to register a consumer")
	var forever chan struct{}
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()
	log.Printf("[*] Waiting for message. To exit press CTRL +C")
	<-forever
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

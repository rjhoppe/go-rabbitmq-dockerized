package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		fmt.Println("Error reading .env file: %v", err)
	}
	pwd := os.Getenv("D_COMPOSE_PWD")
	// RabbitMQ port = 5672
	addr := fmt.Sprintf("%s%s%s", "amqp://guest:", pwd, "@localhost:5672/")
	conn, err := amqp.Dial(addr)
	failOnError(err, "Failed to connect to RabbitMQ")

	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello",
		false,
		false,
		false,
		false,
		nil,
	)

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
			Body: []byte(body),
		},
	)
	failOnError(err, "Failed to publish a message")
	log.Printf("[x] Sent %s\n", body)
}

func failOnError(err error, msg string)  {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

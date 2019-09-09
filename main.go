package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"time"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	// set rand seed
	rand.Seed(int64(time.Now().UnixNano()))

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"task_queue", // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed to declare a queue")

	// generate some playload and send to rabbitmq continuously
	for {
		p := payload{
			Name: "John",
			Age:  30,
		}
		body, err := json.Marshal(p)
		failOnError(err, "fail to generate message body")

		err = ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,
			amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				ContentType:  "text/plain",
				Body:         body,
			})
		failOnError(err, "Failed to publish a message")
		log.Printf(" [x] Sent %s", body)

		time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
	}
}

type payload struct {
	Name string
	Age  int
}

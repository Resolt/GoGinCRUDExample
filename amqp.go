package main

import (
	"fmt"
	"os"

	"github.com/streadway/amqp"
)

type taskhandler struct {
	ac           *amqp.Connection
	queueName    string
	exchangeName string
	uri          string
}

func getTaskhandler() (th *taskhandler, err error) {
	host := os.Getenv("AMQP_HOST")
	user := os.Getenv("AMQP_USER")
	pass := os.Getenv("AMQP_PASS")
	port := os.Getenv("AMQP_PORT")
	vhost := os.Getenv("AMQP_VHOST")
	exchange := os.Getenv("AMQP_EXCHANGE")
	queue := os.Getenv("AMQP_QUEUE")

	uri := fmt.Sprintf("amqp://%s:%s@%s:%s/%s", user, pass, host, port, vhost)

	ac, err := amqp.Dial(uri)
	if err != nil {
		return
	}

	ch, err := ac.Channel()
	if err != nil {
		return
	}

	err = ch.ExchangeDeclare(
		exchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return
	}

	_, err = ch.QueueDeclare(
		queue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return
	}

	th = &taskhandler{
		ac:           ac,
		queueName:    queue,
		exchangeName: exchange,
		uri:          uri,
	}

	return
}

func (t *taskhandler) sendTask(routingKey string, body string) (err error) {
	ch, err := t.ac.Channel()
	if err != nil {
		return
	}
	defer ch.Close()

	err = ch.Publish(
		t.exchangeName,
		routingKey,
		true,
		true,
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "text/plain",
			ContentEncoding: "",
			Body:            []byte(body),
			DeliveryMode:    amqp.Persistent,
			Priority:        0,
		},
	)
	return
}

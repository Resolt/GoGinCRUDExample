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

	th = &taskhandler{}

	th.exchangeName = os.Getenv("AMQP_EXCHANGE")
	th.queueName = os.Getenv("AMQP_QUEUE")

	th.uri = fmt.Sprintf("amqp://%s:%s@%s:%s/%s", user, pass, host, port, vhost)

	th.ac, err = amqp.Dial(th.uri)
	if err != nil {
		return
	}

	ch, err := th.ac.Channel()
	if err != nil {
		return
	}
	defer func() { err = ch.Close() }()

	err = ch.ExchangeDeclare(
		th.exchangeName,
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

	return
}

func (th *taskhandler) sendTask(routingKey string, body string) (err error) {
	ch, err := th.ac.Channel()
	if err != nil {
		return
	}
	defer func() { err = ch.Close() }()

	_, err = ch.QueueDeclare(
		routingKey,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return
	}

	err = ch.Publish(
		"",
		routingKey,
		true,
		false,
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

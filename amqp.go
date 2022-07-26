package main

import (
	"errors"
	"fmt"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
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

	// Create taskhandler
	th = &taskhandler{
		exchangeName: os.Getenv("AMQP_EXCHANGE"),
		queueName:    os.Getenv("AMQP_QUEUE"),
		uri:          fmt.Sprintf("amqp://%s:%s@%s:%s/%s", user, pass, host, port, vhost),
	}

	// Dial AMQP server
	th.ac, err = amqp.Dial(th.uri)
	if err != nil {
		return
	}

	// Ensure a channel can be created
	ch, err := th.ac.Channel()
	if err != nil {
		return
	}
	defer func() { err = ch.Close() }()

	// Declare exchange
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

func (th *taskhandler) sendTask(routingKey string, body string, log *logrus.Logger) (err error) {
	// Create channel with a single redial attempt
	ch, err := th.ac.Channel()
	if err != nil {
		if errors.Is(err, amqp.ErrClosed) {
			log.Warn("Connection to AMQP server lost")
			th.ac, err = amqp.Dial(th.uri)
			if err != nil {
				return
			}
			ch, err = th.ac.Channel()
			if err != nil {
				return
			}
			log.Info("Succesfully reconnected to AMQP server")
		} else {
			return
		}
	}
	defer func() { err = ch.Close() }()

	// Declare queue
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

	// Send task
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

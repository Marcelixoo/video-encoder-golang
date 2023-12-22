package queue

import (
	"fmt"
	"log"
	"os"

	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	User              string
	Password          string
	Host              string
	Port              string
	Vhost             string
	ConsumerQueueName string
	ConsumerName      string
	AutoAck           bool
	Args              amqp.Table
	Channel           *amqp.Channel
}

func NewRabbitMQ() *RabbitMQ {
	rabbitMQArgs := amqp.Table{
		"x-dead-letter-exchange": os.Getenv("RABBITMQ_DLX"),
	}

	return &RabbitMQ{
		User:              os.Getenv("RABBITMQ_DEFAULT_USER"),
		Password:          os.Getenv("RABBITMQ_DEFAULT_PASSWORD"),
		Host:              os.Getenv("RABBITMQ_DEFAULT_HOST"),
		Port:              os.Getenv("RABBITMQ_DEFAULT_PORT"),
		Vhost:             os.Getenv("RABBITMQ_DEFAULT_VHOST"),
		ConsumerName:      os.Getenv("RABBITMQ_CONSUMER_NAME"),
		ConsumerQueueName: os.Getenv("RABBITMQ_CONSUMER_QUEUE_NAME"),
		AutoAck:           false,
		Args:              rabbitMQArgs,
	}
}

func (r *RabbitMQ) Connect() *amqp.Channel {
	dns := fmt.Sprintf("amqp://%s:%s@:%s:%s%s",
		r.User, r.Password, r.Host, r.Port, r.Vhost)

	conn, err := amqp.Dial(dns)
	failOnError(err, "failed to connect to RabbitMQ")

	r.Channel, err = conn.Channel()
	failOnError(err, "failed to open a channel to RabbitMQ connection")

	return r.Channel
}

// Consume reads message from a queue and
// publish them to a message channel (Go chan type)
func (r *RabbitMQ) Consume(messageChannel chan amqp.Delivery) {
	q, err := r.Channel.QueueDeclare(
		r.ConsumerQueueName,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		r.Args,
	)
	failOnError(err, "failed to declared queue")

	incomingMessage, err := r.Channel.Consume(
		q.Name,
		r.ConsumerName,
		r.AutoAck,
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	failOnError(err, "failed to register a consumer")

	go func() {
		for message := range incomingMessage {
			log.Println("got new message")
			messageChannel <- message
		}

		log.Println("will close RabbitMQ channel")
		close(messageChannel)
	}()
}

// Notify publishes a message to the specified
// exchange with the given routing key.
func (r *RabbitMQ) Notify(message, contentType, exchange, routingKey string) error {
	return r.Channel.Publish(
		exchange,
		routingKey,
		false, // mandatory
		false, //immediate
		amqp.Publishing{
			ContentType: contentType,
			Body:        []byte(message),
		},
	)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

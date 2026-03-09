package rabbitmq

import (
	"context"
	"encoding/json"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

type UserDeletedEvent struct {
	UserID    string    `json:"userId"`
	DeletedAt time.Time `json:"deletedAt"`
}

func NewRabbitMQ(connectionUrl string) (*RabbitMQ, error) {
	rabbit_conn, err := amqp.Dial(connectionUrl)
	if err != nil {
		return nil, err
	}
	// defer rabbit_conn.Close()

	rabbit_ch, err := rabbit_conn.Channel()
	if err != nil {
		rabbit_conn.Close()
		return nil, err
	}
	// defer rabbit_ch.Close()

	return &RabbitMQ{
		conn: rabbit_conn,
		ch:   rabbit_ch,
	}, nil
}

func (r *RabbitMQ) SendDeleteAccountMessage(ctx context.Context, userId string) (string, error) {

	exchangeName := "user_events"
	// Creating exchange with type "fanout"
	err := r.ch.ExchangeDeclare(
		exchangeName, // name
		"fanout",     // type
		false,        // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // args
	)
	if err != nil {
		return "", err
	}

	event := UserDeletedEvent{
		UserID:    userId,
		DeletedAt: time.Now(),
	}
	body, _ := json.Marshal(event)

	err = r.ch.Publish(
		exchangeName, // exchange
		"",           // routing key (for fanout it is ignored)
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		return "", err
	}

	time.Sleep(500 * time.Millisecond)
	return " [x] Published user.deleted:" + userId, nil
}

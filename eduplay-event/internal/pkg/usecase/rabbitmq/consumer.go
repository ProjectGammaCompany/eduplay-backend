package rabbitmq

import (
	"encoding/json"
	"log/slog"

	amqp "github.com/rabbitmq/amqp091-go"
)

type storage interface {
}

type RabbitMQ struct {
	conn    *amqp.Connection
	ch      *amqp.Channel
	log     *slog.Logger
	storage storage
}

type UserDeletedEvent struct {
	UserID    string `json:"userId"`
	DeletedAt string `json:"deletedAt"`
}

func NewRabbitMQ(connectionUrl string, log *slog.Logger, st storage) (*RabbitMQ, error) {
	rabbit_conn, err := amqp.Dial(connectionUrl)
	if err != nil {
		return nil, err
	}
	// defer rabbit_conn.Close()

	rabbit_ch, err := rabbit_conn.Channel()
	if err != nil {
		err := rabbit_conn.Close()
		if err != nil {
			return nil, err
		}
		return nil, err
	}
	// defer rabbit_ch.Close()

	return &RabbitMQ{
		conn:    rabbit_conn,
		ch:      rabbit_ch,
		log:     log,
		storage: st,
	}, nil
}

func (r *RabbitMQ) Close() error {
	if err := r.ch.Close(); err != nil {
		return err
	}
	if err := r.conn.Close(); err != nil {
		return err
	}
	return nil
}

func (r *RabbitMQ) ReceiveUserDeletedMessage() {
	exchangeName := "user_events"

	// Declared exchange
	err := r.ch.ExchangeDeclare(
		exchangeName, "fanout", false, false, false, false, nil,
	)
	if err != nil {
		r.log.Error("ExchangeDeclare", "error", err)
		return
	}

	// Создаём временную очередь с уникальным именем
	q, err := r.ch.QueueDeclare(
		"",    // пустое имя = сгенерируется
		false, // durable
		false, // autoDelete
		true,  // exclusive (умрёт, когда отключится consumer)
		false, // noWait
		nil,
	)
	if err != nil {
		r.log.Error("QueueDeclare", "error", err)
		return
	}

	// Привязываем очередь к exchange
	err = r.ch.QueueBind(
		q.Name,       // очередь
		"",           // routing key
		exchangeName, // exchange
		false,
		nil,
	)
	if err != nil {
		r.log.Error("QueueBind", "error", err)
		return
	}

	msgs, err := r.ch.Consume(
		q.Name, "", false, false, false, false, nil,
	)
	if err != nil {
		r.log.Error("Consume", "error", err)
		return
	}

	r.log.Info("Жду событий user-service...")

	// Бесконечно читаем сообщения
	for d := range msgs {
		var event UserDeletedEvent
		if err := json.Unmarshal(d.Body, &event); err != nil {
			r.log.Error("Ошибка парсинга JSON", "error", err)
			err = d.Ack(false) // или d.Nack(false, false)
			if err != nil {
				r.log.Error("Ошибка подтверждения", "error", err)
				return
			}
			continue
		}

		r.log.Info("Получено событие", "userID", event.UserID)

		// TODO: вызвать метод storage для удаления данных пользователя
		// Например: r.storage.DeleteUserData(event.UserID)
		// fmt.Println("OAAAAAAAAAAAAAAAAAAA USER DELETED")
		r.log.Info("User deleted", "userID", event.UserID)

		err = d.Ack(false)
		if err != nil {
			r.log.Error("Ошибка подтверждения", "error", err)
			return
		}
	}

	// Если цикл завершился (канал закрыт), логируем
	r.log.Info("Канал consumer'а закрыт, выходим из ReceiveUserDeletedMessage")
}

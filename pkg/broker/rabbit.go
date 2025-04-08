package broker

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitConfig struct {
	User     string `env:"RABBIT_USER" yaml:"user" env-required:"true"`
	Password string `env:"RABBIT_PASSWORD" yaml:"password" env-required:"true"`
	Host     string `env:"RABBIT_HOST" yaml:"host" env-required:"true"`
	Vhost    string `env:"RABBIT_VHOST" yaml:"vhost" env-required:"true"`
}

type RabbitClient struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func ConnectRabbitMQ(cfg RabbitConfig) (*amqp.Connection, error) {
	return amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s/%s", cfg.User, cfg.Password, cfg.Host, cfg.Vhost))
}

func NewRabbitClient(conn *amqp.Connection) (RabbitClient, error) {
	ch, err := conn.Channel()
	if err != nil {
		return RabbitClient{}, err
	}

	if err := ch.Confirm(false); err != nil {
		return RabbitClient{}, err
	}

	return RabbitClient{
		conn,
		ch,
	}, nil
}

func (rc RabbitClient) Close() error {
	return rc.ch.Close()
}

func (rc RabbitClient) CreateQueue(qname string, durable, autodelete bool) (amqp.Queue, error) {
	return rc.ch.QueueDeclare(qname, durable, autodelete, false, false, nil)
}

func (rc RabbitClient) CreateExchange(ename, kind string, durable, autodelete bool) error {
	return rc.ch.ExchangeDeclare(ename, kind, durable, autodelete, false, false, nil)
}

func (rc RabbitClient) CreateBinding(qname, key, exchange string) error {
	return rc.ch.QueueBind(qname, key, exchange, false, nil)
}

func (rc RabbitClient) Send(ctx context.Context, exchange, routingKey string, data amqp.Publishing) error {
	confirmation, err := rc.ch.PublishWithDeferredConfirmWithContext(ctx,
		exchange,
		routingKey,
		true,
		false,
		data,
	)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	confirmation.Wait()
	return nil
}

func (rc RabbitClient) Consume(queue, consumer string, autoAck bool) (<-chan amqp.Delivery, error) {
	return rc.ch.Consume(queue, consumer, autoAck, false, false, false, nil)
}

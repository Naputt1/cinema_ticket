package events

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/naputt/cinema-ticket-backend/internal/domain"
)

type RabbitMQPublisher struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	mu      sync.Mutex
}

func NewRabbitMQPublisher(url string) (*RabbitMQPublisher, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	err = ch.ExchangeDeclare(
		"cinema",
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, err
	}

	return &RabbitMQPublisher{
		conn:    conn,
		channel: ch,
	}, nil
}

func (p *RabbitMQPublisher) Publish(ctx context.Context, event domain.Event) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	routingKey := string(event.Type)
	if len(routingKey) > 5 && routingKey[:5] == "seat." {
		routingKey = "seat.*"
	} else if len(routingKey) > 7 && routingKey[:7] == "booking" {
		routingKey = "booking.*"
	}

	return p.channel.PublishWithContext(ctx,
		"cinema",
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

func (p *RabbitMQPublisher) Close() {
	p.channel.Close()
	p.conn.Close()
}

type ConsumerHandler func(domain.Event) error

type RabbitMQConsumer struct {
	conn      *amqp.Connection
	channel   *amqp.Channel
	queueName string
	bindings  []string
	handler   ConsumerHandler
	done      chan struct{}
}

func NewRabbitMQConsumer(url string, queueName string, bindings []string, handler ConsumerHandler) (*RabbitMQConsumer, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	err = ch.ExchangeDeclare(
		"cinema",
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, err
	}

	q, err := ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, err
	}

	for _, binding := range bindings {
		err = ch.QueueBind(
			q.Name,
			binding,
			"cinema",
			false,
			nil,
		)
		if err != nil {
			ch.Close()
			conn.Close()
			return nil, err
		}
	}

	return &RabbitMQConsumer{
		conn:      conn,
		channel:   ch,
		queueName: q.Name,
		bindings:  bindings,
		handler:   handler,
		done:      make(chan struct{}),
	}, nil
}

func (c *RabbitMQConsumer) Start() error {
	msgs, err := c.channel.Consume(
		c.queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			var event domain.Event
			if err := json.Unmarshal(msg.Body, &event); err != nil {
				log.Printf("failed to unmarshal event: %v", err)
				continue
			}
			if err := c.handler(event); err != nil {
				log.Printf("failed to handle event %s: %v", event.Type, err)
			}
		}
		close(c.done)
	}()

	return nil
}

func (c *RabbitMQConsumer) Stop() {
	c.channel.Close()
	c.conn.Close()
	<-c.done
}

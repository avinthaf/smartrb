package mq

import (
	"context"
	"log"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Client struct {
	conn  *amqp.Connection
	pubCh *amqp.Channel
	mu    sync.Mutex
}

func NewClient(_ context.Context, url string) (*Client, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, err
	}
	return &Client{conn: conn, pubCh: ch}, nil
}

func (c *Client) Close() error {
	if c.pubCh != nil {
		_ = c.pubCh.Close()
	}
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

var defaultClient *Client

func SetDefault(c *Client) { defaultClient = c }
func Default() *Client     { return defaultClient }

func Publish(topicName string, routingKey string, message string) {
	c := Default()
	if c == nil {
		log.Panic("mq default client not initialized")
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if err := c.pubCh.ExchangeDeclare(topicName, "topic", true, false, false, false, nil); err != nil {
		log.Panicf("Failed to declare an exchange: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := c.pubCh.PublishWithContext(ctx, topicName, routingKey, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(message),
	}); err != nil {
		log.Panicf("Failed to publish a message: %v", err)
	}
	log.Printf(" [x] Sent %s\n", message)
}

type MessageHandler func(topic string, message string)

func Subscribe(topicName string, routingKeys []string, handler MessageHandler) {
	c := Default()
	if c == nil {
		log.Panic("mq default client not initialized")
	}

	ch, err := c.conn.Channel()
	if err != nil {
		log.Panicf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	if err := ch.ExchangeDeclare(topicName, "topic", true, false, false, false, nil); err != nil {
		log.Panicf("Failed to declare an exchange: %v", err)
	}
	q, err := ch.QueueDeclare("", false, false, true, false, nil)
	if err != nil {
		log.Panicf("Failed to declare a queue: %v", err)
	}

	for _, key := range routingKeys {
		if err := ch.QueueBind(q.Name, key, topicName, false, nil); err != nil {
			log.Panicf("Failed to bind queue to exchange: %v", err)
		}
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Panicf("Failed to register a consumer: %v", err)
	}

	log.Printf(" [*] Waiting for messages on exchange '%s'. To exit press CTRL+C", topicName)
	for d := range msgs {
		log.Printf("Received a message: %s with routing key: %s", d.Body, d.RoutingKey)
		handler(d.RoutingKey, string(d.Body))
	}
}

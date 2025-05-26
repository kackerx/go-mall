package messaging

import (
	"context"
	"encoding/json"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"github.com/yourusername/go-d3shop/pkg/events"
)

// RabbitMQPublisher RabbitMQ事件发布器
type RabbitMQPublisher struct {
	connection *amqp091.Connection
	channel    *amqp091.Channel
	exchange   string
}

// NewRabbitMQPublisher 创建RabbitMQ发布器
func NewRabbitMQPublisher(url, exchange string) (*RabbitMQPublisher, error) {
	conn, err := amqp091.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	// 声明交换机
	err = ch.ExchangeDeclare(
		exchange,
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
		connection: conn,
		channel:    ch,
		exchange:   exchange,
	}, nil
}

// PublishAsync 发布集成事件
func (p *RabbitMQPublisher) PublishAsync(ctx context.Context, event events.IIntegrationEvent) error {
	// 创建事件信封
	envelope := events.IntegrationEventEnvelope{
		EventType: event.EventName(),
		Metadata: events.EventMetadata{
			EventID:    event.EventID(),
			OccurredAt: event.OccurredAt(),
		},
	}

	// 序列化事件数据
	eventData, err := json.Marshal(event)
	if err != nil {
		return err
	}
	envelope.EventData = eventData

	// 序列化信封
	body, err := json.Marshal(envelope)
	if err != nil {
		return err
	}

	// 发布消息
	return p.channel.PublishWithContext(
		ctx,
		p.exchange,
		event.EventName(), // routing key
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
			Timestamp:   time.Now(),
		},
	)
}

// Close 关闭连接
func (p *RabbitMQPublisher) Close() error {
	if p.channel != nil {
		p.channel.Close()
	}
	if p.connection != nil {
		return p.connection.Close()
	}
	return nil
}

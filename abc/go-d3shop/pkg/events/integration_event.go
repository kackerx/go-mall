package events

import (
	"context"
	"encoding/json"
	"time"
)

// IIntegrationEvent 集成事件接口
type IIntegrationEvent interface {
	EventName() string
	EventID() string
	OccurredAt() time.Time
}

// BaseIntegrationEvent 基础集成事件
type BaseIntegrationEvent struct {
	ID         string    `json:"id"`
	OccurredOn time.Time `json:"occurredAt"`
}

// EventID 获取事件ID
func (e BaseIntegrationEvent) EventID() string {
	return e.ID
}

// OccurredAt 获取事件发生时间
func (e BaseIntegrationEvent) OccurredAt() time.Time {
	return e.OccurredOn
}

// IIntegrationEventHandler 集成事件处理器接口
type IIntegrationEventHandler[T IIntegrationEvent] interface {
	HandleAsync(ctx context.Context, event T) error
}

// IIntegrationEventPublisher 集成事件发布器接口
type IIntegrationEventPublisher interface {
	PublishAsync(ctx context.Context, event IIntegrationEvent) error
}

// IIntegrationEventConverter 集成事件转换器接口
type IIntegrationEventConverter[TDomainEvent any, TIntegrationEvent IIntegrationEvent] interface {
	Convert(domainEvent TDomainEvent) TIntegrationEvent
}

// IntegrationEventEnvelope 集成事件信封
type IntegrationEventEnvelope struct {
	EventType string          `json:"eventType"`
	EventData json.RawMessage `json:"eventData"`
	Metadata  EventMetadata   `json:"metadata"`
}

// EventMetadata 事件元数据
type EventMetadata struct {
	EventID       string    `json:"eventId"`
	OccurredAt    time.Time `json:"occurredAt"`
	CorrelationID string    `json:"correlationId,omitempty"`
	CausationID   string    `json:"causationId,omitempty"`
	UserID        string    `json:"userId,omitempty"`
}

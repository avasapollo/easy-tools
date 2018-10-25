package broker

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
)

type Config struct {
	AppName string
}

type Message struct {
	ID        string
	EventName string
	Producer  string
	Data      json.RawMessage
	Version   int
}

func NewMessage(producer, eventName string, version int, data json.RawMessage) (*Message, error) {
	// validation
	if producer == "" {
		return nil, NewPublishMessageError(ErrorProducerEmpty)
	}

	if eventName == "" {
		return nil, NewPublishMessageError(ErrorEventNameEmpty)
	}

	if version == 0 {
		return nil, NewPublishMessageError(ErrorEventNameEmpty)
	}

	return &Message{
		ID:        uuid.New().String(),
		EventName: eventName,
		Producer:  producer,
		Data:      data,
		Version:   version,
	}, nil
}

type SubscribeRoute struct {
	EventName     string
	Version       int
	SubscribeFunc SubscribeHandleFunc
}

type SubscribeHandleFunc func(ctx context.Context, msg *Message) error

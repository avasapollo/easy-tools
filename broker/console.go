package broker

import (
	"context"
	"encoding/json"

	"github.com/sirupsen/logrus"
)

type ConsoleConfig struct {
	*Config
}

type ConsoleService struct {
	config *ConsoleConfig
	logger *logrus.Entry
	bus    chan *Message
	routes []*SubscribeRoute
}

func NewConsoleService(le *logrus.Entry, config *ConsoleConfig) (EventBusService, error) {
	return &ConsoleService{
		config: config,
		logger: le,
		bus:    make(chan *Message),
	}, nil
}

func (console *ConsoleService) Subscribe(routes []*SubscribeRoute) error {
	console.routes = routes
	go console.listening()
	return nil
}

func (console ConsoleService) listening() {
	logger := console.logger.WithField("operation", "listening")
	for {
		msg := <-console.bus
		for _, r := range console.routes {
			if r.EventName == msg.EventName && r.Version == msg.Version {
				ctx := context.WithValue(context.Background(), "service", "broker")
				if err := r.SubscribeFunc(ctx, msg); err != nil {
					logger.WithField("error", err.Error()).
						Error("couldn't possible handle the message, received error from the subscribe function")
				}
				logger.WithFields(logrus.Fields{
					"message_id":    msg.ID,
					"event_name":    msg.EventName,
					"event_version": msg.Version,
				}).Debug("the subscription for this message is not configured")
			}
		}
	}
}

func (console ConsoleService) HandleMessageInBus(msg *Message) {
	console.bus <- msg
}

func (console ConsoleService) PublishMessage(ctx context.Context, eventName string, version int,
	data interface{}) (*Message, error) {
	raw, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	msg, err := NewMessage(console.config.AppName, eventName, version, raw)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

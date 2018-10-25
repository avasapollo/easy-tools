package broker

import "context"

type EventBusService interface {
	Subscribe(routes []*SubscribeRoute) error
	PublishMessage(ctx context.Context, eventName string, version int, data interface{}) (*Message, error)
}

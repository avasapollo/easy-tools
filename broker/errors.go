package broker

import "fmt"

var (
	ErrorProducerEmpty  = fmt.Errorf("producer can not be empty")
	ErrorEventNameEmpty = fmt.Errorf("event name can not be empty")
	ErrorEventVersion   = fmt.Errorf("event verson can not be 0")
)

type PublishMessageError struct {
	text string
}

func NewPublishMessageError(errs ...error) error {
	text := "couldn't possible publish the message,"
	for _, err := range errs {
		text = fmt.Sprintf("%s %s", text, err.Error())
	}
	return &PublishMessageError{
		text: text,
	}
}

func (e PublishMessageError) Error() string {
	return e.text
}

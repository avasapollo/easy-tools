package broker

import (
	"context"
	"os"
	"testing"

	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var (
	le     *logrus.Entry
	broker EventBusService
)

func TestMain(m *testing.M) {
	le = logrus.New().WithField("service", "testing")

	broker, _ = NewConsoleService(le, &ConsoleConfig{
		Config: &Config{
			AppName: "testing-api",
		},
	})
	os.Exit(m.Run())
}

func TestConsoleService_PublishMessage(t *testing.T) {
	user := struct {
		Name    string
		Surname string
	}{
		Name:    "andrea",
		Surname: "vasapollo",
	}

	msg, err := broker.PublishMessage(context.Background(), "user_created", 1, &user)

	assert.Nil(t, err)
	assert.NotNil(t, msg)
	assert.Equal(t, msg.Producer, "testing-api")
	assert.Equal(t, msg.EventName, "user_created")
	assert.Equal(t, msg.Version, 1)
}

func TestConsoleService_Subscribe(t *testing.T) {
	msg, _ := NewMessage("testing-api", "user_created", 1, nil)

	broker.Subscribe([]*SubscribeRoute{
		{
			EventName:     "user_created",
			Version:       1,
			SubscribeFunc: handleMessage,
		},
	})

	broker.(*ConsoleService).HandleMessageInBus(msg)

	time.Sleep(3 * time.Second)
}

func handleMessage(ctx context.Context, msg *Message) error {
	fmt.Printf("message: %v \n", msg)
	return nil
}

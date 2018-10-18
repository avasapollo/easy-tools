package metrics

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestService_WrapDatabaseTransaction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dbConfig := DatabaseTransactionDetails{
		DatabaseType: MongoDatabaseType,
		DatabaseName: "rayhack-database",
		Collection:   "rayhack-collection",
		Operation:    "find_one_beer",
	}

	provider := NewMockProvider(ctrl)
	provider.EXPECT().GetDatabaseTransaction(context.Background(), dbConfig).Return(nil)

	svc, err := NewMetricService(le, provider)
	assert.Nil(t, err)

	myName := "my name is "

	svc.WrapDatabaseTransaction(context.Background(), dbConfig, func() {
		myName = myName + "rayhack"
	})
	assert.Nil(t, err)
	assert.Equal(t, "my name is rayhack", myName)
	t.Log(myName)
}

package metrics

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/newrelic/go-agent"
	"github.com/stretchr/testify/assert"
)

func TestNewrelic_WrapDatabaseTransaction(t *testing.T) {
	license := getEnvNewrelicLicenseKey()
	if license == "" {
		t.Skip("NEWRELIC_LICENSE_KEY is not in the env, the test will skipped")
	}

	// get newrelic app
	app, err := getNewrelicApp(license)
	assert.Nil(t, err, "newrelic is not configured")

	for i := 0; i < 20; i++ {
		// start newrelic transaction
		tnx := app.StartTransaction("database_transaction_test", nil, nil)

		dbConfig := DatabaseTransactionDetails{
			DatabaseType: MongoDatabaseType,
			DatabaseName: "rayhack-database",
			Collection:   "rayhack-collection",
			Operation:    "find",
		}
		ctx := context.WithValue(context.Background(), NerelicTransactionContextKeyName, tnx)

		provider := NewNewrelicProvider(le)

		svc, err := NewMetricService(le, provider)
		assert.Nil(t, err)

		myName := "my name is "
		svc.WrapDatabaseTransaction(ctx, dbConfig, func() {
			time.Sleep(2 * time.Second)
			myName = myName + "rayhack"
		})

		err = tnx.End()
		assert.Nil(t, err)

		assert.Equal(t, myName, "my name is rayhack")
		t.Log("my name is rayhack")
	}

	time.Sleep(2 * time.Minute)
}

func getEnvNewrelicLicenseKey() string {
	return os.Getenv("NEWRELIC_LICENSE_KEY")
}

func getNewrelicApp(license string) (newrelic.Application, error) {
	conf := newrelic.NewConfig("beers-api", license)
	conf.Logger = newrelic.NewDebugLogger(os.Stdout)
	conf.CrossApplicationTracer.Enabled = false
	conf.DistributedTracer.Enabled = true
	conf.Enabled = true
	err := conf.Validate()
	if err != nil {
		return nil, err
	}

	return newrelic.NewApplication(conf)
}

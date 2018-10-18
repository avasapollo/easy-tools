package metrics

import (
	"context"
	"fmt"

	"github.com/newrelic/go-agent"
	"github.com/sirupsen/logrus"
)

const (
	NerelicTransactionContextKeyName = "newrelic_transaction"
)

type newrelicProvider struct {
	logger *logrus.Entry
}

func NewNewrelicProvider(le *logrus.Entry) Provider {
	return &newrelicProvider{
		logger: le.WithField("metric_provider", "newrelic"),
	}
}

func (n newrelicProvider) GetDatabaseTransaction(ctx context.Context, config DatabaseTransactionDetails) Transaction {
	transaction, err := NewNewrelicDatabaseTransaction(ctx, config)
	if err != nil {
		n.logger.WithField("error", err.Error()).Error("newrelic database transaction error")
		return nil
	}
	return transaction
}

type newRelicDatabaseTransaction struct {
	transaction     newrelic.Transaction
	segment         *newrelic.DatastoreSegment
	config          DatabaseTransactionDetails
	transactionType TransactionType
	provider        MetricProvider
}

func NewNewrelicDatabaseTransaction(ctx context.Context, details DatabaseTransactionDetails) (Transaction, error) {
	txn := ctx.Value(NerelicTransactionContextKeyName).(newrelic.Transaction)
	if txn == nil {
		return nil, fmt.Errorf("couldn't possible recorder the transaction in newrelic, newrelic transaction is not in the context")
	}

	return &newRelicDatabaseTransaction{
		transaction: txn,
		segment: &newrelic.DatastoreSegment{
			Product:      parseNewrelicDatabaseType(details.DatabaseType),
			Collection:   details.Collection,
			Operation:    details.Operation,
			DatabaseName: details.DatabaseName,
		},
		config:          details,
		transactionType: DatabaseTransactionType,
		provider:        NewrelicProvider,
	}, nil
}

func (svc newRelicDatabaseTransaction) StartTransaction() error {
	svc.segment.StartTime = newrelic.StartSegmentNow(svc.transaction)
	return nil
}

func (svc newRelicDatabaseTransaction) EndTransaction() error {
	return svc.segment.End()
}

func (svc newRelicDatabaseTransaction) GetProvider() MetricProvider {
	return svc.provider
}

func (svc newRelicDatabaseTransaction) GetType() TransactionType {
	return svc.transactionType
}

func parseNewrelicDatabaseType(databaseType DatabaseType) newrelic.DatastoreProduct {
	switch databaseType {
	case MongoDatabaseType:
		return newrelic.DatastoreMongoDB
	default:
		return newrelic.DatastoreMongoDB
	}
}

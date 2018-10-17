package metrics

import "context"

type MetricService interface {
	WrapDatabaseTransaction(ctx context.Context, config DatabaseTransactionDetails, fn func())
}

type Provider interface {
	GetDatabaseTransaction(ctx context.Context, config DatabaseTransactionDetails) Transaction
}

type Transaction interface {
	GetProvider() MetricProvider
	GetType() TransactionType

	StartTransaction()
	EndTransaction()
}

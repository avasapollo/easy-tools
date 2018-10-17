package metrics

import (
	"context"

	"github.com/sirupsen/logrus"
)

type service struct {
	le        *logrus.Entry
	providers []Provider
}

func NewMetricService(le *logrus.Entry, providers ...Provider) (MetricService, error) {
	return &service{
		le:        le,
		providers: providers,
	}, nil
}

func (svc service) WrapDatabaseTransaction(ctx context.Context, config DatabaseTransactionDetails, fn func()) {
	var transactions SliceTransaction
	for _, p := range svc.providers {
		transaction := p.GetDatabaseTransaction(ctx, config)
		if transaction == nil {
			continue
		}
		transactions = append(transactions, transaction)
	}

	transactions.Start()
	defer transactions.End()

	fn()
}

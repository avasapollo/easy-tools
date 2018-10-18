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
	var transactions []Transaction
	for _, p := range svc.providers {
		transaction := p.GetDatabaseTransaction(ctx, config)
		if transaction == nil {
			continue
		}
		transactions = append(transactions, transaction)
	}

	// start transactions
	for _, t := range transactions {
		if err := t.StartTransaction(); err != nil {
			svc.le.WithFields(logrus.Fields{
				"error":            err.Error(),
				"provider":         t.GetProvider(),
				"transaction_type": t.GetType(),
			}).Error("couldn't possible start the transaction")
			continue
		}
	}

	fn()
	// end transactions
	for _, t := range transactions {
		if err := t.EndTransaction(); err != nil {
			svc.le.WithFields(logrus.Fields{
				"error":            err.Error(),
				"provider":         t.GetProvider(),
				"transaction_type": t.GetType(),
			}).Error("couldn't possible end the transaction")
			continue
		}
	}
}

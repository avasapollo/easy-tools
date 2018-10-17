package metrics

type DatabaseTransactionDetails struct {
	DatabaseType DatabaseType
	DatabaseName string
	Collection   string
	Operation    string
}

type DatabaseType string

var MongoDatabaseType DatabaseType

type TransactionType string

var DatabaseTransactionType TransactionType = "database"

type MetricProvider string

var NewrelicProvider MetricProvider = "newrelic"

type SliceTransaction []Transaction

func (slice SliceTransaction) Start() {
	for _, s := range slice {
		s.StartTransaction()
	}
}

func (slice SliceTransaction) End() {
	for _, s := range slice {
		s.EndTransaction()
	}
}

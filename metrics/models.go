package metrics

type DatabaseTransactionDetails struct {
	DatabaseType DatabaseType
	DatabaseName string
	Collection   string
	Operation    string
}

type DatabaseType string

var MongoDatabaseType DatabaseType = "MongoDB"

type TransactionType string

var DatabaseTransactionType TransactionType = "database"

type MetricProvider string

var NewrelicProvider MetricProvider = "newrelic"

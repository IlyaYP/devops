package storage

type MetricStorage interface {
	PutMetric(MetricType, MetricName, MetricValue string) error
	GetMetric(MetricType, MetricName string) (string, error)
	ReadMetrics() map[string]map[string]string
	Ping() error
}

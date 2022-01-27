package storage

import "context"

type MetricStorage interface {
	PutMetric(MetricType, MetricName, MetricValue string) error
	GetMetric(MetricType, MetricName string) (string, error)
	ReadMetrics() map[string]map[string]string
	Ping(ctx context.Context) error
}

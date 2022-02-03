package storage

import "context"

type MetricStorage interface {
	PutMetric(ctx context.Context, MetricType, MetricName, MetricValue string) error
	GetMetric(ctx context.Context, MetricType, MetricName string) (string, error)
	ReadMetrics(ctx context.Context) map[string]map[string]string
	Ping(ctx context.Context) error
}

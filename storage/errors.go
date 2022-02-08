package storage

import (
	"fmt"
)

type TypeError struct {
	Type string
}

func (e *TypeError) Error() string {
	return fmt.Sprintf("No such Type: %v", e.Type)
}

type MetricError struct {
	Metric string
}

func (e *MetricError) Error() string {
	return fmt.Sprintf("No such Metric: %v", e.Metric)
}

func NewTypeError(text string) error {
	return &TypeError{text}
}

func NewMetricError(text string) error {
	return &MetricError{text}
}

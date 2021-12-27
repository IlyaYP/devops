package inmemory

import (
	"context"
	"fmt"
	"github.com/IlyaYP/devops/storage"
	"sync"
)

var _ storage.MetricStorage = (*Storage)(nil) // Q: Вот это для чего? я ещё не изучил (

type Storage struct {
	*sync.RWMutex
	mtr map[string]map[string]string
} //{mtr: make(map[string]map[string]string)}

func NewStorage() *Storage {
	s := Storage{mtr: make(map[string]map[string]string)}
	s.mtr["counter"] = make(map[string]string)
	s.mtr["gauge"] = make(map[string]string)
	return &s
}

func (s Storage) PutMetric(ctx context.Context, MetricType, MetricName, MetricValue string) error {
	// To write to the storage, take the write lock:
	s.Lock()
	defer s.Unlock()
	t, ok := s.mtr[MetricType]
	if !ok {
		return fmt.Errorf("wrong type")
	}
	t[MetricName] = MetricValue

	return nil
}

func (s Storage) GetMetric(ctx context.Context, MetricType, MetricName string) (string, error) {
	// To read from the storage, take the read lock:
	s.RLock()
	defer s.RUnlock()
	t, ok := s.mtr[MetricType]
	if !ok {
		return "", fmt.Errorf("wrong type")
	}

	n, ok := t[MetricName]
	if !ok {
		return "", fmt.Errorf("no such metric")
	}

	return n, nil
}
func (s Storage) ReadMetrics(ctx context.Context) map[string]map[string]string {
	s.RLock()
	defer s.RUnlock()
	ret := make(map[string]map[string]string)
	for k, v := range s.mtr {
		ret[k] = make(map[string]string)
		for kk, vv := range v {
			ret[k][kk] = vv
		}
	}

	return ret
}

package inmemory

import (
	"context"
	"errors"
	"github.com/IlyaYP/devops/storage"
	"log"
	"strconv"
	"sync"
)

var _ storage.MetricStorage = (*MemStorage)(nil)

type MemStorage struct {
	sync.RWMutex
	Mtr map[string]map[string]string
} //{mtr: make(map[string]map[string]string)}

func NewMemStorage() *MemStorage {
	s := MemStorage{Mtr: make(map[string]map[string]string)}
	s.Mtr["counter"] = make(map[string]string)
	s.Mtr["gauge"] = make(map[string]string)
	return &s
}

func (s *MemStorage) Ping(ctx context.Context) error {
	if s.Mtr != nil {
		return nil
	}
	return errors.New("MemStorage.Ping")
}

func (s *MemStorage) PutMetric(ctx context.Context, MetricType, MetricName, MetricValue string) error {
	//fmt.Println("Put:", MetricType, MetricName, MetricValue)
	// To write to the storage, take the write lock:
	s.Lock()
	defer s.Unlock()
	t, ok := s.Mtr[MetricType]
	if !ok {
		log.Println("Error:", MetricType, MetricName, MetricValue)
		return storage.NewTypeError(MetricType)
	}
	if MetricType == "gauge" {
		if _, err := strconv.ParseFloat(MetricValue, 64); err != nil {
			return err
		}
	} else if MetricType == "counter" {
		v, err := strconv.ParseInt(MetricValue, 10, 64)
		if err != nil {
			return err
		}
		tv, ok := t[MetricName]
		if !ok {
			tv = "0"
		}
		vv, err := strconv.ParseInt(tv, 10, 64)
		if err != nil {
			return err
		}
		MetricValue = strconv.FormatInt(v+vv, 10)
	}

	t[MetricName] = MetricValue

	return nil
}

func (s *MemStorage) GetMetric(ctx context.Context, MetricType, MetricName string) (string, error) {
	// To read from the storage, take the read lock:
	s.RLock()
	defer s.RUnlock()
	t, ok := s.Mtr[MetricType]
	if !ok {
		return "", storage.NewTypeError(MetricType)
	}

	n, ok := t[MetricName]
	if !ok {
		return "", storage.NewMetricError(MetricName)
	}

	return n, nil
}

func (s *MemStorage) ReadMetrics(ctx context.Context) map[string]map[string]string {
	s.RLock()
	defer s.RUnlock()
	ret := make(map[string]map[string]string)
	for k, v := range s.Mtr {
		ret[k] = make(map[string]string)
		for kk, vv := range v {
			ret[k][kk] = vv
		}
	}

	return ret
}

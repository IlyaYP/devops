package inmemory

import (
	"fmt"
	"github.com/IlyaYP/devops/storage"
	"strconv"
	"sync"
)

var _ storage.MetricStorage = (*MemStorage)(nil) // Q: Вот это для чего? я ещё не изучил (

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

func (s *MemStorage) PutMetric(MetricType, MetricName, MetricValue string) error {
	//fmt.Println("Put:", MetricType, MetricName, MetricValue)
	// To write to the storage, take the write lock:
	s.Lock()
	defer s.Unlock()
	t, ok := s.Mtr[MetricType]
	if !ok {
		fmt.Println("Error:", MetricType, MetricName, MetricValue)
		return fmt.Errorf("wrong type")
	}
	if MetricType == "gauge" {
		if _, err := strconv.ParseFloat(MetricValue, 64); err != nil {
			return fmt.Errorf("wrong value")
		}
	} else if MetricType == "counter" {
		v, err := strconv.ParseInt(MetricValue, 10, 64)
		if err != nil {
			return fmt.Errorf("wrong value")
		}
		tv, ok := t[MetricName]
		if !ok {
			tv = "0"
		}
		vv, err := strconv.ParseInt(tv, 10, 64)
		if err != nil {
			return fmt.Errorf("strconv.ParseInt error")
		}
		MetricValue = strconv.FormatInt(v+vv, 10)
	}

	t[MetricName] = MetricValue

	return nil
}

func (s *MemStorage) GetMetric(MetricType, MetricName string) (string, error) {
	// To read from the storage, take the read lock:
	s.RLock()
	defer s.RUnlock()
	t, ok := s.Mtr[MetricType]
	if !ok {
		return "", fmt.Errorf("wrong type")
	}

	n, ok := t[MetricName]
	if !ok {
		return "", fmt.Errorf("no such metric")
	}

	return n, nil
}
func (s *MemStorage) ReadMetrics() map[string]map[string]string {
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

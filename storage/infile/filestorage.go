package infile

import (
	"encoding/json"
	"github.com/IlyaYP/devops/storage"
	"github.com/IlyaYP/devops/storage/inmemory"
	"os"
)

var _ storage.MetricStorage = (*FileStorage)(nil)

type FileStorage struct {
	inmemory.MemStorage
	file    *os.File
	encoder *json.Encoder
	decoder *json.Decoder
}

func NewFileStorage(filename string) (*FileStorage, error) {

	s := new(FileStorage)
	s.Mtr = make(map[string]map[string]string)
	s.Mtr["counter"] = make(map[string]string)
	s.Mtr["gauge"] = make(map[string]string)

	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return nil, err
	}
	s.encoder = json.NewEncoder(file)
	s.decoder = json.NewDecoder(file)
	return s, nil
}

func (c *FileStorage) Close() error {
	return c.file.Close()
}

func (c *FileStorage) PutMetric(MetricType, MetricName, MetricValue string) error {
	return c.MemStorage.PutMetric(MetricType, MetricName, MetricValue)
}

func (c *FileStorage) GetMetric(MetricType, MetricName string) (string, error) {
	return c.MemStorage.GetMetric(MetricType, MetricName)
}

func (c *FileStorage) ReadMetrics() map[string]map[string]string {
	return c.MemStorage.ReadMetrics()
}

package infile

import (
	"encoding/json"
	"fmt"
	"github.com/IlyaYP/devops/cmd/server/config"
	"github.com/IlyaYP/devops/internal"
	"github.com/IlyaYP/devops/storage"
	"github.com/IlyaYP/devops/storage/inmemory"
	"log"
	"os"
	"strconv"
)

var _ storage.MetricStorage = (*FileStorage)(nil)

type FileStorage struct {
	inmemory.MemStorage
	file    *os.File
	encoder *json.Encoder
	decoder *json.Decoder
	cfg     *config.Config
}

func NewFileStorage(cfg *config.Config) (*FileStorage, error) {

	s := new(FileStorage)
	s.cfg = cfg
	s.Mtr = make(map[string]map[string]string)
	s.Mtr["counter"] = make(map[string]string)
	s.Mtr["gauge"] = make(map[string]string)

	//file, err := os.OpenFile(cfg.StoreFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	file, err := os.OpenFile(cfg.StoreFile, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	s.file = file
	s.encoder = json.NewEncoder(file)
	s.decoder = json.NewDecoder(file)

	if cfg.Restore {
		if err := s.Restore(); err != nil {
			log.Println(err.Error())
		}
	}
	return s, nil
}

func (c *FileStorage) Close() error {
	_ = c.Save()
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

func (c *FileStorage) Restore() error {
	for c.decoder.More() {
		var m internal.Metrics
		var MetricValue string
		// decode an array value (Message)
		err := c.decoder.Decode(&m)
		if err != nil {
			return err
		}
		//log.Println(m)
		if m.MType == "gauge" {
			MetricValue = fmt.Sprintf("%v", *m.Value)
		} else if m.MType == "counter" {
			MetricValue = fmt.Sprintf("%v", *m.Delta)
		} else {
			return fmt.Errorf("wrong type")
		}

		if err := c.MemStorage.PutMetric(m.MType, m.ID, MetricValue); err != nil {
			return err
		}
	}

	return nil
}

func (c *FileStorage) Save() error {

	if _, err := c.file.Seek(0, 0); err != nil {
		log.Println(err.Error())
		log.Println("bla bla bla")
	}
	if err := c.file.Truncate(0); err != nil {
		log.Println(err)
		log.Println("bla bla bla")
	}

	m := c.MemStorage.ReadMetrics()
	for k, v := range m {
		for kk, vv := range v {
			if k == "gauge" {
				value, err := strconv.ParseFloat(vv, 64)
				if err != nil {
					log.Println(err)
				}
				mt := internal.Metrics{ID: kk, MType: k, Value: &value}
				c.encoder.Encode(mt)
			} else if k == "counter" {
				delta, err := strconv.ParseInt(vv, 10, 64)
				if err != nil {
					log.Println(err)
				}
				mt := internal.Metrics{ID: kk, MType: k, Delta: &delta}
				c.encoder.Encode(mt)
			}
		}
	}

	return nil
}

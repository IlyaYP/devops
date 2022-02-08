package infile

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IlyaYP/devops/cmd/server/config"
	"github.com/IlyaYP/devops/internal"
	"github.com/IlyaYP/devops/storage"
	"github.com/IlyaYP/devops/storage/inmemory"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

var _ storage.MetricStorage = (*FileStorage)(nil)

type FileStorage struct {
	inmemory.MemStorage
	file          *os.File
	encoder       *json.Encoder
	decoder       *json.Decoder
	cfg           *config.Config
	lastWrite     time.Time
	StoreInterval time.Duration
	fm            sync.RWMutex
	dirty         bool
}

func NewFileStorage(ctx context.Context, cfg *config.Config) (*FileStorage, error) {

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
	//s.StoreInterval = time.Duration(cfg.StoreInterval) * time.Second
	s.StoreInterval = cfg.StoreInterval
	s.dirty = false

	if cfg.Restore {
		if err := s.Restore(ctx); err != nil {
			log.Println(err.Error())
		}
	}
	return s, nil
}

func (c *FileStorage) Close(ctx context.Context) error {
	_ = c.Save(ctx)
	return c.file.Close()
}

func (c *FileStorage) PutMetric(ctx context.Context, MetricType, MetricName, MetricValue string) error {
	if err := c.MemStorage.PutMetric(ctx, MetricType, MetricName, MetricValue); err != nil {
		return err
	}
	c.dirty = true
	if time.Since(c.lastWrite) >= c.StoreInterval {
		time.AfterFunc(time.Duration(1)*time.Second, func() { c.DelayedSave(ctx) })
		c.lastWrite = time.Now().Add(time.Duration(1) * time.Second)
	}

	//log.Println(MetricType, MetricName, MetricValue, c.dirty)
	c.dirty = true
	return nil
}

func (c *FileStorage) DelayedSave(ctx context.Context) {
	if err := c.Save(ctx); err != nil {
		log.Println(err)
		return
	}
	time.AfterFunc(c.StoreInterval, func() { c.CheckState(ctx) })
	c.lastWrite = time.Now()
}

func (c *FileStorage) CheckState(ctx context.Context) {
	if c.dirty {
		c.DelayedSave(ctx)
	}
}

func (c *FileStorage) Restore(ctx context.Context) error {
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
			return storage.NewTypeError(m.MType)
		}

		if err := c.MemStorage.PutMetric(ctx, m.MType, m.ID, MetricValue); err != nil {
			return err
		}
	}

	return nil
}

func (c *FileStorage) Save(ctx context.Context) error {
	c.fm.Lock()
	defer c.fm.Unlock()

	log.Println("Writing data to", c.cfg.StoreFile)
	if _, err := c.file.Seek(0, 0); err != nil {
		log.Println(err.Error())
	}
	if err := c.file.Truncate(0); err != nil {
		log.Println(err)
	}

	m := c.MemStorage.ReadMetrics(ctx)
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
	c.dirty = false
	return nil
}

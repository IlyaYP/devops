package inmemory

import (
	"context"
	"reflect"
	"sync"
	"testing"
)

func TestNewStorage(t *testing.T) {
	tests := []struct {
		name string
		want *Storage
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewStorage(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewStorage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorage_GetMetric(t *testing.T) {
	type fields struct {
		RWMutex *sync.RWMutex
		mtr     map[string]map[string]string
	}
	type args struct {
		ctx        context.Context
		MetricType string
		MetricName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Storage{
				RWMutex: tt.fields.RWMutex,
				mtr:     tt.fields.mtr,
			}
			got, err := s.GetMetric(tt.args.ctx, tt.args.MetricType, tt.args.MetricName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMetric() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetMetric() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorage_PutMetric(t *testing.T) {
	type fields struct {
		RWMutex *sync.RWMutex
		mtr     map[string]map[string]string
	}
	type args struct {
		ctx         context.Context
		MetricType  string
		MetricName  string
		MetricValue string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Storage{
				RWMutex: tt.fields.RWMutex,
				mtr:     tt.fields.mtr,
			}
			if err := s.PutMetric(tt.args.ctx, tt.args.MetricType, tt.args.MetricName, tt.args.MetricValue); (err != nil) != tt.wantErr {
				t.Errorf("PutMetric() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStorage_ReadMetrics(t *testing.T) {
	type fields struct {
		RWMutex *sync.RWMutex
		mtr     map[string]map[string]string
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string]map[string]string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Storage{
				RWMutex: tt.fields.RWMutex,
				mtr:     tt.fields.mtr,
			}
			if got := s.ReadMetrics(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadMetrics() = %v, want %v", got, tt.want)
			}
		})
	}
}

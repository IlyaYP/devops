package postgres

import (
	"context"
	"fmt"
	"github.com/IlyaYP/devops/storage"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"strconv"
	"time"
)

var _ storage.MetricStorage = (*Postgres)(nil)

type Postgres struct {
	DBDsn   string
	pool    *pgxpool.Pool
	timeout time.Duration
}

func NewPostgres(ctx context.Context, DBDsn string) (*Postgres, error) {
	s := new(Postgres)
	s.DBDsn = DBDsn
	s.timeout = time.Duration(1) * time.Second
	pool, err := pgxpool.Connect(ctx, s.DBDsn)
	if err != nil {
		log.Printf("Unable to connect to database: %v\n", err)
		return nil, err
	}
	s.pool = pool

	// create DB
	//dbName := "devops"
	//_, err = pool.Exec(ctx, "create database "+dbName)
	//if err != nil {
	//	log.Printf("Unable to create database: %v\n", err)
	//	return nil, err
	//}

	// creating tables
	_, err = pool.Exec(ctx, "CREATE TABLE IF NOT EXISTS gauges ( id varchar(40) primary key, value double precision);"+
		"CREATE TABLE IF NOT EXISTS counters ( id varchar(40) primary key, delta bigint);")
	if err != nil {
		log.Printf("Unable to create table: %v\n", err)
		return nil, err
	}

	return s, nil
}

func (c *Postgres) Ping(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	return c.pool.Ping(ctx)
}

func (c *Postgres) Close() {
	c.pool.Close()
}

func (c *Postgres) PutMetric(MetricType, MetricName, MetricValue string) error {
	switch MetricType {
	case "gauge":
		value, err := strconv.ParseFloat(MetricValue, 64)
		if err != nil {
			return err
		}
		return c.PutGauge(MetricName, value)
	case "counter":
		delta, err := strconv.ParseInt(MetricValue, 10, 64)
		if err != nil {
			return err
		}
		return c.PutCounter(MetricName, delta)
	default:
		return fmt.Errorf("wrong type %s", MetricType)
	}
}

func (c *Postgres) PutGauge(MetricName string, value float64) error {
	_, err := c.pool.Exec(context.Background(), `insert into gauges(id, value) values ($1, $2)
	on conflict (id) do update set value=excluded.value`, MetricName, value)
	return err
}

func (c *Postgres) PutCounter(MetricName string, delta int64) error {
	_, err := c.pool.Exec(context.Background(), `insert into counters(id, delta) values ($1, $2)
	on conflict (id) do update set delta= counters.delta + excluded.delta`, MetricName, delta)
	return err
}

func (c *Postgres) GetMetric(MetricType, MetricName string) (string, error) {
	if MetricType == "gauge" { // TODO: split into small functions
		var value float64
		err := c.pool.QueryRow(context.Background(), "select value from gauges where id=$1", MetricName).Scan(&value)
		switch err {
		case nil:
			return fmt.Sprintf("%v", value), nil
		case pgx.ErrNoRows:
			return "", fmt.Errorf("no such metric %s in DB", MetricName)
		default:
			return "", err
		}
	} else if MetricType == "counter" {
		var delta int64
		err := c.pool.QueryRow(context.Background(), "select delta from counters where id=$1", MetricName).Scan(&delta)
		switch err {
		case nil:
			return fmt.Sprintf("%v", delta), nil
		case pgx.ErrNoRows:
			return "", fmt.Errorf("no such metric %s in DB", MetricName)
		default:
			return "", err
		}
	} else {
		return "", fmt.Errorf("wrong type %s", MetricType)
	}
}
func (c *Postgres) ReadMetrics() map[string]map[string]string {
	ret := make(map[string]map[string]string)
	ret["counters"] = make(map[string]string)
	counters, _ := c.pool.Query(context.Background(), "select * from counters")
	defer counters.Close()

	for counters.Next() {
		var id string
		var delta int64
		err := counters.Scan(&id, &delta)
		if err != nil {
			log.Println(err)
		}
		ret["counters"][id] = fmt.Sprintf("%v", delta)
	}

	ret["gauges"] = make(map[string]string)
	gauges, _ := c.pool.Query(context.Background(), "select * from gauges")
	defer gauges.Close()

	for gauges.Next() {
		var id string
		var value float64
		err := gauges.Scan(&id, &value)
		if err != nil {
			log.Println(err)
		}
		ret["gauges"][id] = fmt.Sprintf("%v", value)
	}

	return ret

	//return map[string]map[string]string{"counter": {"PollCount": "1"}}
}

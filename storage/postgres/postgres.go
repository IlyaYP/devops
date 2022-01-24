package postgres

import (
	"context"
	"fmt"
	"github.com/IlyaYP/devops/storage"
	"github.com/jackc/pgx/v4"
	"log"
	"os"
)

var _ storage.MetricStorage = (*Postgres)(nil)

type Postgres struct {
	DBDsn string
	conn  *pgx.Conn
}

func NewPostgres(DBDsn string) (*Postgres, error) {
	s := new(Postgres)
	s.DBDsn = DBDsn
	conn, err := pgx.Connect(context.Background(), s.DBDsn)
	if err != nil {
		log.Printf("Unable to connect to database: %v\n", err)
		return nil, err
	}
	s.conn = conn
	return s, nil
}

func (c *Postgres) Ping() error {
	return c.conn.Ping(context.Background())
}

func (c *Postgres) Close() {
	c.conn.Close(context.Background())
}

func (c *Postgres) PutMetric(MetricType, MetricName, MetricValue string) error {

	return nil
}

func (c *Postgres) GetMetric(MetricType, MetricName string) (string, error) {

	return "777", nil
}
func (c *Postgres) ReadMetrics() map[string]map[string]string {
	return map[string]map[string]string{"counter": {"PollCount": "1"}}
}

func test() {
	// urlExample := "postgres://username:password@localhost:5432/database_name"
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	var name string
	var weight int64
	err = conn.QueryRow(context.Background(), "select name, weight from widgets where id=$1", 42).Scan(&name, &weight)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(name, weight)
}

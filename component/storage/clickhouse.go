package storage

import (
	"context"
	"database/sql"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/pkg/errors"
)

var clickhouseDatabase *sql.DB

var clickhouseConn *sql.Conn

func ClickhouseDatabase() *sql.DB {
	return clickhouseDatabase
}

func ClickhouseConn() *sql.Conn {
	return clickhouseConn
}

func InitClickhouseDatabase(ctx context.Context, dsn string) error {
	opt, err := clickhouse.ParseDSN(dsn)
	if err != nil {
		return errors.Wrap(err, "clickhouse parse dsn error")
	}
	db := clickhouse.OpenDB(opt)
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(10)
	db.SetConnMaxLifetime(time.Hour)

	clickhouseDatabase = db
	conn, err := db.Conn(ctx)
	if err != nil {
		return errors.Wrap(err, "clickhouse connect error")
	}
	clickhouseConn = conn
	return nil
}

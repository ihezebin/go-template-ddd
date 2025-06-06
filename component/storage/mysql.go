package storage

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	driverMysql "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/ihezebin/go-template-ddd/config"
)

var mysqlDatabase *gorm.DB

func MySQLDatabase() *gorm.DB {
	return mysqlDatabase
}

// InitMySQLClient init mysql storage client
// dsn: "user:password@tcp(host:port)/dbname?charset=utf8&parseTime=True&loc=Local"
func InitMySQLClient(ctx context.Context, dsn string) error {
	mysqlDsn, err := parseMysqlDSN(dsn)
	if err != nil {
		return errors.Wrap(err, "mysql parse dsn error")
	}
	if err = runMigration(ctx, mysqlDsn); err != nil {
		return errors.Wrap(err, "mysql run migration error")
	}

	db, err := gorm.Open(mysql.Open(mysqlDsn.FormatDSN()), &gorm.Config{})
	if err != nil {
		return errors.Wrap(err, "mysql connect error")
	}
	// https://gorm.io/docs/generic_interface.html#Connection-Pool
	sqlDB, err := db.DB()
	if err != nil {
		return errors.Wrap(err, "mysql get sql db error")
	}

	err = sqlDB.Ping()
	if err != nil {
		return errors.Wrap(err, "mysql ping error")
	}

	sqlDB.SetConnMaxIdleTime(time.Minute * 30)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(50)

	mysqlDatabase = db
	return nil
}

func runMigration(ctx context.Context, mysqlDsn *driverMysql.Config) error {
	dbName := mysqlDsn.DBName
	mysqlDsn.DBName = ""
	db, err := sql.Open("mysql", mysqlDsn.FormatDSN())
	if err != nil {
		return errors.Wrap(err, "connect error")
	}
	defer db.Close()

	// create database
	_, err = db.ExecContext(ctx, fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", dbName))
	if err != nil {
		return errors.Wrap(err, "create database error")
	}

	// run migration
	mysqlDsn.DBName = dbName
	migrationDir := "file://" + filepath.Join(config.GetConfig().Pwd, "migration")
	dbDsn := "mysql://" + mysqlDsn.FormatDSN()

	m, err := migrate.New(migrationDir, dbDsn)
	if err != nil {
		return errors.Wrap(err, "migrate new error")
	}

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return errors.Wrap(err, "migrate up error")
	}

	defer m.Close()

	return nil
}

func parseMysqlDSN(dsn string) (*driverMysql.Config, error) {
	u, err := url.Parse(dsn)
	if err != nil {
		return nil, errors.Wrap(err, "dsn parse error")
	}
	dbName := strings.TrimPrefix(u.Path, "/")
	if dbName == "" {
		return nil, errors.New("db name is empty")
	}

	password, _ := u.User.Password()
	dsn = fmt.Sprintf("%s:%s@tcp(%s)%s?%s", u.User.Username(), password, u.Host, u.Path, u.Query().Encode())
	mysqlDsn, err := driverMysql.ParseDSN(dsn)
	if err != nil {
		return nil, errors.Wrap(err, "mysql dsn parse error")
	}

	return mysqlDsn, nil
}

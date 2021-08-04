package database

import (
	"context"
	"database/sql"
	logger "log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

//Driver list
const (
	DriverPostgres = "postgres"
	DriverMysql    = "mysql"
)

type (
	Database struct {
		*sqlx.DB
	}

	DB struct {
		DBConnection   *Database
		DSN            string
		Driver         string
		MaxIdleConn    int
		MaxConn        int
		IdleTimeoutSec time.Duration
	}

	// Store is used to persist master and slave DB connection
	Store struct {
		Write *Database
		Read  *Database
	}

	DatabaseConfig struct {
		WriteDSN           string
		ReadDSN            string
		MaxIdleConn        int
		MaxConn            int
		IdleTimeoutSeconds time.Duration
	}
)

func (s *Store) GetWrite() *Database {
	return s.Write
}

func (s *Store) GetRead() *Database {
	return s.Read
}

func NewSqlDriverWithSqlxMock(sqlx *sqlx.DB) *Store {
	db := &DB{}
	db.DBConnection = &Database{DB: sqlx}

	return &Store{
		Write: db.DBConnection,
		Read:  db.DBConnection,
	}
}

func NewSqlDriverMock(sql *sql.DB) *Store {
	sqlxDB := sqlx.NewDb(sql, "sqlmock")

	db := &DB{}
	db.DBConnection = &Database{DB: sqlxDB}

	return &Store{
		Write: db.DBConnection,
		Read:  db.DBConnection,
	}
}

func NewSqlDriver(cfg DatabaseConfig, dbDriver string) *Store {
	write := &DB{
		DSN:            cfg.WriteDSN,
		Driver:         dbDriver,
		MaxIdleConn:    cfg.MaxIdleConn,
		MaxConn:        cfg.MaxConn,
		IdleTimeoutSec: cfg.IdleTimeoutSeconds,
	}
	write.Connect()

	read := &DB{
		DSN:            cfg.ReadDSN,
		Driver:         dbDriver,
		MaxIdleConn:    cfg.MaxIdleConn,
		MaxConn:        cfg.MaxConn,
		IdleTimeoutSec: cfg.IdleTimeoutSeconds,
	}
	read.Connect()

	return &Store{
		Read:  read.DBConnection,
		Write: write.DBConnection,
	}
}

func (d *DB) Connect() {
	db, err := sqlx.Open(d.Driver, d.DSN)
	if err != nil {
		logger.Fatalln("Error when open write database connection: ", err.Error())
	}

	db.SetMaxOpenConns(d.MaxConn)
	db.SetMaxIdleConns(d.MaxIdleConn)
	db.SetConnMaxLifetime(d.IdleTimeoutSec * time.Second)

	err = db.Ping()
	if err != nil {
		logger.Fatalln("Error when ping write db: ", err)
	}

	d.DBConnection = &Database{DB: db}
}

func (db *Database) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "[db][ExecContext]")
	defer span.Finish()

	ext.DBStatement.Set(span, query)
	ext.DBInstance.Set(span, db.DriverName())
	ext.DBType.Set(span, "sql")
	span.SetTag("db.values", args)

	return db.DB.ExecContext(ctx, query, args...)
}

func (db *Database) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "[db][QueryContext]")
	defer span.Finish()

	ext.DBStatement.Set(span, query)
	ext.DBInstance.Set(span, db.DriverName())
	ext.DBType.Set(span, "sql")
	span.SetTag("db.values", args)

	return db.DB.QueryContext(ctx, query, args...)
}

func (db *Database) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	span, ctx := opentracing.StartSpanFromContext(ctx, "[db][QueryRowContext]")
	defer span.Finish()

	ext.DBStatement.Set(span, query)
	ext.DBInstance.Set(span, db.DriverName())
	ext.DBType.Set(span, "sql")
	span.SetTag("db.values", args)

	return db.DB.QueryRowContext(ctx, query, args...)
}

func (db *Database) QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row {
	span, ctx := opentracing.StartSpanFromContext(ctx, "[db][QueryRowxContext]")
	defer span.Finish()

	ext.DBStatement.Set(span, query)
	ext.DBInstance.Set(span, db.DriverName())
	ext.DBType.Set(span, "sql")
	span.SetTag("db.values", args)

	return db.DB.QueryRowxContext(ctx, query, args...)
}

func (db *Database) QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "[db][QueryxContext]")
	defer span.Finish()

	ext.DBStatement.Set(span, query)
	ext.DBInstance.Set(span, db.DriverName())
	ext.DBType.Set(span, "sql")
	span.SetTag("db.values", args)

	return db.DB.QueryxContext(ctx, query, args...)
}

func (db *Database) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "[db][SelectContext]")
	defer span.Finish()

	ext.DBStatement.Set(span, query)
	ext.DBInstance.Set(span, db.DriverName())
	ext.DBType.Set(span, "sql")
	span.SetTag("db.values", args)

	return db.DB.SelectContext(ctx, dest, query, args...)
}

func (db *Database) NamedQueryContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "[db][NamedQueryContext]")
	defer span.Finish()

	ext.DBStatement.Set(span, query)
	ext.DBInstance.Set(span, db.DriverName())
	ext.DBType.Set(span, "sql")
	span.SetTag("db.values", args)

	return db.DB.NamedQueryContext(ctx, query, args)
}

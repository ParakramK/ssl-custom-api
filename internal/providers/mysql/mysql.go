package mysql

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"net/url"
	"time"

	"github.com/go-sql-driver/mysql"
	gormmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Provider struct {
	db *gorm.DB
}

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

// readonlyConnector wraps the MySQL driver connector and forces every new
// pooled connection into read-only mode. Writes then fail fast instead of
// silently succeeding. This is a guardrail, not a substitute for a
// SELECT-only database user, which is still recommended.
type readonlyConnector struct {
	driver.Connector
}

func (c readonlyConnector) Connect(ctx context.Context) (driver.Conn, error) {
	conn, err := c.Connector.Connect(ctx)
	if err != nil {
		return nil, err
	}

	execer, ok := conn.(driver.ExecerContext)
	if !ok {
		_ = conn.Close()
		return nil, fmt.Errorf("mysql: driver connection does not support exec")
	}

	if _, err := execer.ExecContext(ctx, "SET SESSION TRANSACTION READ ONLY", nil); err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("mysql: enforce read-only connection: %w", err)
	}

	return conn, nil
}

func NewProvider(cfg Config) (*Provider, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s",
		url.QueryEscape(cfg.User),
		url.QueryEscape(cfg.Password),
		cfg.Host,
		cfg.Port,
		cfg.Database,
	)

	parsed, err := mysql.ParseDSN(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse MySQL DSN: %w", err)
	}

	connector, err := mysql.NewConnector(parsed)
	if err != nil {
		return nil, fmt.Errorf("create MySQL connector: %w", err)
	}

	sqldb := sql.OpenDB(readonlyConnector{Connector: connector})

	sqldb.SetMaxOpenConns(20)
	sqldb.SetMaxIdleConns(5)
	sqldb.SetConnMaxLifetime(30 * time.Minute)
	sqldb.SetConnMaxIdleTime(5 * time.Minute)

	if err := sqldb.Ping(); err != nil {
		sqldb.Close()
		return nil, fmt.Errorf("ping MySQL: %w", err)
	}

	gormDB, err := gorm.Open(gormmysql.New(gormmysql.Config{Conn: sqldb}), &gorm.Config{
		SkipDefaultTransaction: true,
		// The repositories run a small set of fixed-shape queries, so
		// cache their prepared statements instead of re-preparing.
		PrepareStmt: true,
	})
	if err != nil {
		sqldb.Close()
		return nil, fmt.Errorf("open GORM MySQL: %w", err)
	}

	return &Provider{db: gormDB}, nil
}

// DB exposes the GORM database handle for repositories. All connections
// behind it are read-only (see readonlyConnector).
func (p *Provider) DB() *gorm.DB {
	return p.db
}

func (p *Provider) Close() error {
	sqldb, err := p.db.DB()
	if err != nil {
		return err
	}
	return sqldb.Close()
}

func (p *Provider) Ping() error {
	sqldb, err := p.db.DB()
	if err != nil {
		return err
	}
	if err := sqldb.Ping(); err != nil {
		return fmt.Errorf("ping MySQL: %w", err)
	}

	return nil
}

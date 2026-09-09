package mysql

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"log"
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
// pooled connection into read-only mode.
// GUARDRAIL
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

		return nil, fmt.Errorf(
			"mysql: driver connection does not support ExecContext",
		)
	}

	if _, err := execer.ExecContext(
		ctx,
		"SET SESSION TRANSACTION READ ONLY",
		nil,
	); err != nil {
		_ = conn.Close()

		return nil, fmt.Errorf(
			"mysql: enforce read-only connection: %w",
			err,
		)
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

	sqldb := sql.OpenDB(
		readonlyConnector{
			Connector: connector,
		},
	)

	// Connection pool configuration.
	sqldb.SetMaxOpenConns(20)
	sqldb.SetMaxIdleConns(5)
	sqldb.SetConnMaxLifetime(30 * time.Minute)
	sqldb.SetConnMaxIdleTime(5 * time.Minute)

	pingCtx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	if err := sqldb.PingContext(pingCtx); err != nil {
		_ = sqldb.Close()

		return nil, fmt.Errorf("ping MySQL: %w", err)
	}

	warmCtx, cancel := context.WithTimeout(
		context.Background(),
		30*time.Second,
	)
	defer cancel()

	if err := warmConnections(warmCtx, sqldb, 2); err != nil {
		_ = sqldb.Close()

		return nil, fmt.Errorf(
			"warm MySQL connections: %w",
			err,
		)
	}

	stats := sqldb.Stats()

	log.Printf(
		"MySQL pool warmed: open=%d idle=%d inuse=%d",
		stats.OpenConnections,
		stats.Idle,
		stats.InUse,
	)

	gormDB, err := gorm.Open(
		gormmysql.New(
			gormmysql.Config{
				Conn: sqldb,
			},
		),
		&gorm.Config{
			SkipDefaultTransaction: true,
			PrepareStmt:            false,
		},
	)
	if err != nil {
		_ = sqldb.Close()

		return nil, fmt.Errorf(
			"open GORM MySQL: %w",
			err,
		)
	}

	return &Provider{
		db: gormDB,
	}, nil
}

func warmConnections(
	ctx context.Context,
	db *sql.DB,
	count int,
) error {
	if count <= 0 {
		return nil
	}

	conns := make([]*sql.Conn, 0, count)

	for i := 0; i < count; i++ {
		conn, err := db.Conn(ctx)
		if err != nil {
			for _, c := range conns {
				_ = c.Close()
			}

			return fmt.Errorf(
				"warm connection %d/%d: %w",
				i+1,
				count,
				err,
			)
		}

		conns = append(conns, conn)
	}

	for i, conn := range conns {
		if err := conn.Close(); err != nil {
			for j := i + 1; j < len(conns); j++ {
				_ = conns[j].Close()
			}

			return fmt.Errorf(
				"release warmed connection %d/%d: %w",
				i+1,
				count,
				err,
			)
		}
	}

	return nil
}

func (p *Provider) DB() *gorm.DB {
	return p.db
}

func (p *Provider) Close() error {
	if p == nil || p.db == nil {
		return nil
	}

	sqldb, err := p.db.DB()
	if err != nil {
		return err
	}

	return sqldb.Close()
}

func (p *Provider) Ping() error {
	if p == nil || p.db == nil {
		return fmt.Errorf("MySQL provider is not initialized")
	}

	sqldb, err := p.db.DB()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	if err := sqldb.PingContext(ctx); err != nil {
		return fmt.Errorf("ping MySQL: %w", err)
	}

	return nil
}

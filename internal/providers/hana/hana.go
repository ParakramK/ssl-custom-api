package hana

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"time"

	_ "github.com/SAP/go-hdb/driver"
)

type Provider struct {
	db *sql.DB
}

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
}

func NewProvider(cfg Config) (*Provider, error) {
	dsn := fmt.Sprintf(
		"hdb://%s:%s@%s:%d",
		url.QueryEscape(cfg.User),
		url.QueryEscape(cfg.Password),
		cfg.Host,
		cfg.Port,
	)

	db, err := sql.Open("hdb", dsn)
	if err != nil {
		return nil, fmt.Errorf("open HANA connection: %w", err)
	}

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping HANA: %w", err)
	}

	return &Provider{
		db: db,
	}, nil
}

func (p *Provider) Ping(ctx context.Context) error {
	if err := p.db.PingContext(ctx); err != nil {
		return fmt.Errorf("ping HANA: %w", err)
	}

	return nil
}
func (p *Provider) WithSchema(
	ctx context.Context,
	schema string,
	fn func(*sql.Conn) error,
) error {
	if schema == "" {
		return fmt.Errorf("HANA schema is required")
	}

	conn, err := p.db.Conn(ctx)
	if err != nil {
		return fmt.Errorf("get HANA connection: %w", err)
	}
	defer conn.Close()

	_, err = conn.ExecContext(
		ctx,
		fmt.Sprintf(`SET SCHEMA "%s"`, schema),
	)
	if err != nil {
		return fmt.Errorf("set HANA schema: %w", err)
	}

	return fn(conn)
}

func (p *Provider) Close() error {
	return p.db.Close()
}

func (p *Provider) Query(
	ctx context.Context,
	schema string,
	query string,
	args ...any,
) (*sql.Rows, error) {

	if schema == "" {
		return nil, fmt.Errorf("HANA schema is required")
	}

	conn, err := p.db.Conn(ctx)
	if err != nil {
		return nil, fmt.Errorf("get HANA connection: %w", err)
	}

	_, err = conn.ExecContext(
		ctx,
		fmt.Sprintf(`SET SCHEMA "%s"`, schema),
	)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("set HANA schema: %w", err)
	}

	rows, err := conn.QueryContext(ctx, query, args...)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("execute HANA query: %w", err)
	}

	return rows, nil
}

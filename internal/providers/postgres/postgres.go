package postgres

import (
	"fmt"
	"time"

	gormpostgres "gorm.io/driver/postgres"
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
	SSLMode  string
}

func NewProvider(cfg Config) (*Provider, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.Database,
		cfg.SSLMode,
	)

	db, err := gorm.Open(
		gormpostgres.Open(dsn),
		&gorm.Config{
			SkipDefaultTransaction: true,
			PrepareStmt:            true,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("open PostgreSQL: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get PostgreSQL connection: %w", err)
	}

	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)

	if err := sqlDB.Ping(); err != nil {
		sqlDB.Close()
		return nil, fmt.Errorf("ping PostgreSQL: %w", err)
	}

	return &Provider{db: db}, nil
}

func (p *Provider) DB() *gorm.DB {
	return p.db
}

func (p *Provider) Close() error {
	sqlDB, err := p.db.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}

func (p *Provider) Ping() error {
	sqlDB, err := p.db.DB()
	if err != nil {
		return err
	}

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("ping PostgreSQL: %w", err)
	}

	return nil
}

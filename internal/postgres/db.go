package postgres

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

type Config struct {
	Connection      string        `yaml:"postgresql"`
	MaxOpenConn     int           `yaml:"max_open_conn"`
	MaxIdleConn     int           `yaml:"max_idle_conn"`
	MaxConnLifetime time.Duration `yaml:"max_conn_lifetime"`
	opts            options
}

type options struct {
	Wrapper func(driver.Connector) driver.Connector
}

type AuthConfig struct {
	Host     string
	Port     string
	Database string
	Username string
	Password string
	SSLMode  string

	MaxOpenConn     int           `yaml:"max_open_conn"`
	MaxIdleConn     int           `yaml:"max_idle_conn"`
	MaxConnLifetime time.Duration `yaml:"max_conn_lifetime"`
}

func NewAuth(postrgresConfig *AuthConfig) (*DB, error) {
	dbConnectionString, err := GenerateConnectionString(postrgresConfig)
	if err != nil {
		return nil, fmt.Errorf("can't generate db connection string: %w", err)
	}

	postgresDB, err := NewPostgresDB(&Config{
		Connection:      dbConnectionString,
		MaxOpenConn:     postrgresConfig.MaxOpenConn,
		MaxIdleConn:     postrgresConfig.MaxIdleConn,
		MaxConnLifetime: postrgresConfig.MaxConnLifetime,
	})
	if err != nil {
		return nil, fmt.Errorf("can't create postgres db connect: %w", err)
	}

	return postgresDB, nil
}

func GenerateConnectionString(cfg *AuthConfig) (string, error) {
	if cfg == nil {
		return "", errors.New("config is nil")
	}

	u := url.URL{
		// TODO: scheme в config
		Scheme: "postgres",
		User:   url.UserPassword(cfg.Username, cfg.Password),
		Host:   cfg.Host + ":" + cfg.Port,
		Path:   cfg.Database,
	}

	q := u.Query()
	q.Set("sslmode", cfg.SSLMode)
	u.RawQuery = q.Encode()

	return u.String(), nil
}

func NewPostgresDB(config *Config) (*DB, error) {
	db, err := config.CreateDBX()
	if err != nil {
		return nil, fmt.Errorf("can't create db: %w", err)
	}

	if errPing := db.Ping(); errPing != nil {
		return nil, fmt.Errorf("can't ping database: %w", errPing)
	}

	return &DB{
		DBConn: db,
	}, nil
}

func WithWrapper(cfg Config, wrapper func(drv driver.Connector) driver.Connector) Config {
	if cfg.opts.Wrapper == nil {
		cfg.opts.Wrapper = wrapper
	} else {
		cfg.opts.Wrapper = func(drv driver.Connector) driver.Connector {
			return wrapper(cfg.opts.Wrapper(drv))
		}
	}
	return cfg
}

func New(cfg Config) (*sql.DB, error) {
	return cfg.CreateDB()
}

func NewDBX(cfg Config) (*sqlx.DB, error) {
	return cfg.CreateDBX()
}

// CreateDB creates sql.DB for postgres
func (cfg Config) CreateDB() (*sql.DB, error) {
	var (
		ctor driver.Connector
	)

	drv := stdlib.GetDefaultDriver().(*stdlib.Driver)

	ctor, err := drv.OpenConnector(cfg.Connection)
	if err != nil {
		return nil, err
	}

	if cfg.opts.Wrapper != nil {
		ctor = cfg.opts.Wrapper(ctor)
	}

	db := sql.OpenDB(ctor)

	if cfg.MaxConnLifetime != 0 {
		db.SetConnMaxLifetime(cfg.MaxConnLifetime)
	}

	if cfg.MaxIdleConn != 0 {
		db.SetMaxIdleConns(cfg.MaxIdleConn)
	}

	if cfg.MaxOpenConn != 0 {
		db.SetMaxOpenConns(cfg.MaxOpenConn)
	}

	return db, nil
}

// CreateDBX creates sqlx.DB for postgres
func (cfg Config) CreateDBX() (*sqlx.DB, error) {
	db, err := cfg.CreateDB()
	if err != nil {
		return nil, err
	}
	dbx := sqlx.NewDb(db, "pgx")
	return dbx, nil
}

package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

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

type DB struct {
	DBConn *sqlx.DB
	tx     *sql.Tx
}

// "user:password@/dbname"
func NewConnect(config *AuthConfig) (*DB, error) {
	dbConnectionString, err := GenerateConnectionString(config)
	if err != nil {
		return nil, err
	}
	// Открытие соединения с базой данных
	db, err := sqlx.Open("mysql", dbConnectionString)
	if err != nil {
		return nil, err
	}

	return &DB{
		DBConn: db,
	}, nil
}

// jdbc:mysql://localhost:3306
func GenerateConnectionString(cfg *AuthConfig) (string, error) {
	if cfg == nil {
		return "", errors.New("config is nil")
	}

	query := fmt.Sprintf(
		"%s:%s@(%s:%s)/%s",
		cfg.Username, cfg.Password,
		cfg.Host, cfg.Port,
		cfg.Database,
	)

	return query, nil
}

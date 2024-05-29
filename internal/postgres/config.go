package postgres

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type DB struct {
	DBConn *sqlx.DB
	tx     *sql.Tx
}

func (db *DB) GetDatabases() ([]string, error) {
	rows, err := db.DBConn.Query("SELECT datname FROM pg_database WHERE datistemplate = false;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var databases []string
	for rows.Next() {
		var dbName string
		err := rows.Scan(&dbName)
		if err != nil {
			return nil, err
		}
		databases = append(databases, dbName)
	}

	return databases, nil
}

func (db *DB) GetSchemas() ([]string, error) {
	rows, err := db.DBConn.Query("SELECT schema_name FROM information_schema.schemata")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	schemas := []string{}

	var schemaName string
	for rows.Next() {
		if err := rows.Scan(&schemaName); err != nil {
			return nil, err
		}

		schemas = append(schemas, schemaName)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return schemas, nil
}

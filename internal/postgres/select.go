package postgres

import (
	"database/sql"
	"datasupervision/internal/service"
	"fmt"
)

func (db *DB) Select(query string, args ...any) (*service.TableData, error) {
	var err error
	var rows *sql.Rows
	if db.tx != nil {
		rows, err = db.tx.Query(query, args...)
	} else {
		rows, err = db.DBConn.Query(query, args...)
	}
	if err != nil {
		db.tx = nil
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var data [][]interface{}
	for rows.Next() {
		values := make([]interface{}, len(columns))
		for i := range values {
			values[i] = new(interface{})
		}
		if err := rows.Scan(values...); err != nil {
			return nil, err
		}
		data = append(data, values)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	tableData := &service.TableData{
		Columns: columns,
		Rows:    data,
	}

	return tableData, nil
}

func (db *DB) SelectTableData(schemaName, tableName string) (*service.TableData, error) {
	query := "SELECT * FROM " + schemaName + "." + tableName

	tableData, err := db.Select(query)
	if err != nil {
		return nil, err
	}

	return tableData, nil
}

func (db *DB) SelectWithFilter(schemaName, tableName, columnName, filterValue string) (*service.TableData, error) {
	query := fmt.Sprintf("SELECT * FROM %s.%s WHERE %s = $1", schemaName, tableName, columnName)

	tableData, err := db.Select(query, filterValue)
	if err != nil {
		return nil, err
	}

	return tableData, err
}

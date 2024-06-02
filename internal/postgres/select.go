package postgres

import (
	"database/sql"
	"datasupervision/internal/service"
	"fmt"
	"strconv"
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
	query := "SELECT * FROM " + tableName

	tableData, err := db.Select(query)
	if err != nil {
		return nil, err
	}

	return tableData, nil
}

func (db *DB) SelectWithFilter(schemaName, tableName string, filter service.Filters) (*service.TableData, error) {
	query := fmt.Sprintf("SELECT * FROM %s.%s", schemaName, tableName)
	args := []interface{}{}

	if filter.FilterValue != "" {
		query = fmt.Sprintf("%s WHERE %s =", query, filter.ColumnName)
		query += " $" + strconv.Itoa(len(args)+1)
		args = append(args, filter.FilterValue)
	}

	if filter.Limit > 0 {
		query += " LIMIT $" + strconv.Itoa(len(args)+1)
		args = append(args, filter.Limit)
	}

	if filter.Offset > 0 {
		query += " OFFSET $" + strconv.Itoa(len(args)+1)
		args = append(args, filter.Offset)
	}

	tableData, err := db.Select(query, args...)
	if err != nil {
		return nil, err
	}

	return tableData, nil
}

package mysql

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

	colTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}

	stringTypes := map[string]string{
		"VARCHAR":   "VARCHAR",
		"TIMESTAMP": "TIMESTAMP",
		"CHAR":      "CHAR",
		"ENUM":      "ENUM",
		"TEXT":      "TEXT",
	}

	var data [][]interface{}
	for rows.Next() {
		values := make([]interface{}, len(columns))
		for i := range values {
			_, ok := stringTypes[colTypes[i].DatabaseTypeName()]
			if ok {
				values[i] = toPoint(sql.NullString{})
				continue
			}
			values[i] = new(interface{})
		}

		if err := rows.Scan(values...); err != nil {
			return nil, err
		}
		for i := range len(columns) {
			vvv, ok := values[i].(*sql.NullString)
			if ok {
				if vvv.Valid {
					values[i] = vvv.String
				} else {
					values[i] = "null"
				}
			}
		}

		// Добавление мапы в слайс
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
	query := fmt.Sprintf("SELECT * FROM %s.%s", schemaName, tableName)
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
		query = fmt.Sprintf("%s WHERE %s = ?", query, filter.ColumnName)
		args = append(args, filter.FilterValue)
	}

	if filter.Limit > 0 {
		query += " LIMIT ?"
		args = append(args, filter.Limit)
	}

	if filter.Offset > 0 {
		query += " OFFSET ?"
		args = append(args, filter.Offset)
	}

	tableData, err := db.Select(query, args...)
	if err != nil {
		return nil, err
	}

	return tableData, nil
}

func toPoint[T any](t T) *T {
	return &t
}

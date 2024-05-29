package postgres

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"

	"datasupervision/internal/service"
)

func (db *DB) InsertRow(schemaName, tableName string, row map[string]interface{}) error {
	columns, values := buildInsertRow(row)
	// Подготовка SQL-запроса на вставку данных
	query := fmt.Sprintf("INSERT INTO %s.%s (%s) VALUES (",
		schemaName,
		tableName,
		columns,
	)
	for i := range values {
		query = fmt.Sprintf("%s $%d, ", query, i+1)
	}

	query = query[:len(query)-2]
	query = fmt.Sprintf("%s )", query)

	log.Info(query)

	// Выполнение SQL-запроса
	var err error
	if db.tx != nil {
		_, err = db.tx.Exec(query, values...)
		if err != nil {
			db.tx = nil
			return fmt.Errorf("posgres, exex insert row: %w", err)
		}
		return nil
	}

	_, err = db.DBConn.Exec(query, values...)
	if err != nil {
		return fmt.Errorf("posgres, exex insert row: %w", err)
	}

	return nil
}

// Построение списка столбцов для SQL-запроса
func buildInsertRow(row map[string]interface{}) (string, []interface{}) {
	columns := ""
	values := []interface{}{}
	for column, value := range row {
		columns += column + ","
		values = append(values, value)
	}
	return columns[:len(columns)-1], values // Удаляем последнюю запятую
}

func (db *DB) UpdateRow(schemaName, tableName string, row service.UpdateRow) error {
	// Подготовка SQL-запроса на изменение данных
	query := fmt.Sprintf("UPDATE %s.%s SET", schemaName, tableName)
	for i, column := range row.Сolumns {
		query = fmt.Sprintf("%s %s='%s',", query, column, row.NewRow[i])
	}

	query = query[:len(query)-1] // Удаление последней запятой
	query += " WHERE"
	for i, column := range row.Сolumns {
		query = fmt.Sprintf("%s %s='%s' AND", query, column, row.OldRow[i])
	}

	query = query[:len(query)-3] // Удаление последнего AND
	log.Info(query)

	// Выполнение SQL-запроса
	var err error
	if db.tx != nil {
		_, err = db.tx.Exec(query)
		if err != nil {
			db.tx = nil
			return fmt.Errorf("posgres, exex update row: %w", err)
		}
		return nil
	}

	_, err = db.DBConn.Exec(query)
	if err != nil {
		return fmt.Errorf("posgres, exex update row: %w", err)
	}

	return nil
}

func (db *DB) DeleteRow(schemaName, tableName string, row service.Row) error {
	// Подготовка SQL-запроса на вставку данных
	query := fmt.Sprintf("DELETE FROM %s.%s WHERE", schemaName, tableName)
	for column, value := range row.Row {
		query = fmt.Sprintf("%s %s='%s' AND", query, column, value)
	}

	query = query[:len(query)-3]

	log.Info(query)

	// Выполнение SQL-запроса
	var err error
	if db.tx != nil {
		_, err = db.tx.Exec(query)
		if err != nil {
			db.tx = nil
			return fmt.Errorf("posgres, exex delete row: %w", err)
		}
		return nil
	}

	_, err = db.DBConn.Exec(query)
	if err != nil {
		return fmt.Errorf("posgres, exex delete row: %w", err)
	}

	return nil
}

func (db *DB) InsertRows(schemaName, tableName string, columns []string, rows [][]interface{}) error {
	if len(rows) == 0 {
		return nil
	}

	colNames := strings.Join(columns, ", ")

	var valuePlaceholders []string
	var allValues []interface{}

	for i, row := range rows {
		var placeholders []string
		for j := range row {
			placeholders = append(placeholders, fmt.Sprintf("$%d", i*len(columns)+j+1))
		}
		valuePlaceholders = append(valuePlaceholders, fmt.Sprintf("(%s)", strings.Join(placeholders, ", ")))
		allValues = append(allValues, row...)
	}

	query := fmt.Sprintf("INSERT INTO %s.%s (%s) VALUES %s", schemaName, tableName, colNames, strings.Join(valuePlaceholders, ", "))

	_, err := db.DBConn.Exec(query, allValues...)
	if err != nil {
		return fmt.Errorf("execute statement: %v", err)
	}

	return nil
}

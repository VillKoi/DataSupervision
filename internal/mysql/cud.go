package mysql

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
	for range values {
		query = fmt.Sprintf("%s ?, ", query)
	}

	query = query[:len(query)-2]
	query = fmt.Sprintf("%s )", query)

	log.Info(query)

	// Выполнение SQL-запроса
	var err error
	if db.tx != nil {
		_, err = db.tx.Exec(query, values...)
	} else {
		_, err = db.DBConn.Exec(query, values...)
	}
	if err != nil {
		db.tx = nil
		return err
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
	query := fmt.Sprintf("UPDATE %s.%s SET", schemaName, tableName)
	args := []interface{}{}
	for i, column := range row.Сolumns {
		query = fmt.Sprintf("%s %s=?,", query, column)
		args = append(args, row.NewRow[i])
	}

	query = query[:len(query)-1] // Удаление последней запятой
	query += " WHERE"
	for i, column := range row.Сolumns {
		query = fmt.Sprintf("%s %s=? AND", query, column)
		args = append(args, row.OldRow[i])
	}

	query = query[:len(query)-3] // Удаление последнего AND
	log.Info(query)

	// Выполнение SQL-запроса
	var err error
	if db.tx != nil {
		_, err = db.tx.Exec(query, args...)
	} else {
		_, err = db.DBConn.Exec(query, args...)
	}
	if err != nil {
		db.tx = nil
		return err
	}

	return nil
}

func (db *DB) DeleteRow(schemaName, tableName string, row service.Row) error {
	query := fmt.Sprintf("DELETE FROM %s.%s WHERE", schemaName, tableName)
	args := []interface{}{}
	for column, value := range row.Row {
		query = fmt.Sprintf("%s %s=? AND", query, column)
		args = append(args, value)
	}

	query = query[:len(query)-3]

	log.Info(query)

	var err error
	if db.tx != nil {
		_, err = db.tx.Exec(query, args...)
	} else {
		_, err = db.DBConn.Exec(query, args...)
	}
	if err != nil {
		db.tx = nil
		return err
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

	for _, row := range rows {
		var placeholders []string
		for range row {
			placeholders = append(placeholders, "?")
		}
		valuePlaceholders = append(valuePlaceholders, fmt.Sprintf("(%s)", strings.Join(placeholders, ", ")))
		allValues = append(allValues, row...)
	}

	query := fmt.Sprintf("INSERT INTO `%s`.`%s` (%s) VALUES %s", schemaName, tableName, colNames, strings.Join(valuePlaceholders, ", "))

	stmt, err := db.DBConn.Prepare(query)
	if err != nil {
		return fmt.Errorf("prepare statement: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(allValues...)
	if err != nil {
		return fmt.Errorf("execute statement: %v", err)
	}

	return nil
}

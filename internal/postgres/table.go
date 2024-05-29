package postgres

import (
	"datasupervision/internal/service"
	"fmt"

	log "github.com/sirupsen/logrus"
)

func (db *DB) GetTables(schemaName string) ([]string, error) {
	rows, err := db.DBConn.Query("SELECT table_name FROM information_schema.tables WHERE table_schema = $1", schemaName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			return nil, err
		}
		tables = append(tables, table)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tables, nil
}

func (db *DB) GetColumns(schemaName, tableName string) ([]service.TableColumn, error) {
	// Получение информации о колонках таблицы
	columnQuery := `SELECT column_name, data_type, character_maximum_length, is_nullable
		FROM information_schema.columns
		WHERE table_schema = $1 AND table_name = $2`

	rows, err := db.DBConn.Query(columnQuery, schemaName, tableName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []service.TableColumn
	for rows.Next() {
		var column service.TableColumn
		err := rows.Scan(&column.Name, &column.DataType, &column.CharacterMaximumLength, &column.IsNullable)
		if err != nil {
			return nil, err
		}
		columns = append(columns, column)
	}

	return columns, nil
}

func (db *DB) GetTablesAndColumns(schemaName string) (map[string][]string, error) {
	query := `
		SELECT table_name, column_name
		FROM information_schema.columns
		WHERE table_schema = $1
		ORDER BY table_name, ordinal_position;
	`

	rows, err := db.DBConn.Query(query, schemaName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tableColumns := make(map[string][]string)
	for rows.Next() {
		var tableName, columnName string
		if err := rows.Scan(&tableName, &columnName); err != nil {
			return nil, err
		}
		tableColumns[tableName] = append(tableColumns[tableName], columnName)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tableColumns, nil
}

func (db *DB) CreateTable(schemaName string, tableData service.CreatedTableData) error {
	query := generateCreateTableSQL(schemaName, tableData)

	log.WithField("func", "CreateTable").Info(query)

	_, err := db.DBConn.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) DeleteTable(schemaName, tableName string) error {
	query := fmt.Sprintf("DROP TABLE IF EXISTS %s.%s", schemaName, tableName)

	_, err := db.DBConn.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

var typeMapping = map[string]string{
	"UUID":    "UUID",
	"INT":     "INTEGER",
	"VARCHAR": "VARCHAR",
	"TEXT":    "TEXT",
	"DATE":    "DATE",
	"BOOLEAN": "BOOLEAN",
	"FLOAT":   "FLOAT",
}

func (db *DB) GetColumnTypes() []string {
	return []string{"UUID", "INT", "VARCHAR", "TEXT", "DATE", "BOOLEAN", "FLOAT"}
}

func generateCreateTableSQL(schemaName string, tableData service.CreatedTableData) string {
	query := fmt.Sprintf("CREATE TABLE %s.%s (", schemaName, tableData.TableName)

	primaryKeys := []string{}
	foreignKeys := []string{}

	for i, column := range tableData.Columns {
		if i > 0 {
			query += ", "
		}

		columnType := typeMapping[column.Type]

		query += fmt.Sprintf("%s %s", column.Name, columnType)
		if column.NotNull {
			query += " NOT NULL"
		}
		if column.Primary {
			primaryKeys = append(primaryKeys, column.Name)
		}
		if column.ForeignKey.ForeignTable != "" && column.ForeignKey.ForeignColumn != "" {
			foreignKey := fmt.Sprintf("FOREIGN KEY (%s) REFERENCES %s(%s)", column.Name, column.ForeignKey.ForeignTable, column.ForeignKey.ForeignColumn)
			foreignKeys = append(foreignKeys, foreignKey)
		}
	}

	if len(primaryKeys) > 0 {
		query += fmt.Sprintf(", PRIMARY KEY (%s)", join(primaryKeys, ", "))
	}
	if len(foreignKeys) > 0 {
		query += ", " + join(foreignKeys, ", ")
	}

	query += ");"

	return query
}

func join(elements []string, delimiter string) string {
	result := ""
	for i, element := range elements {
		if i > 0 {
			result += delimiter
		}
		result += element
	}
	return result
}

package mysql

import (
	"datasupervision/internal/service"
)

func (db *DB) GetKeys(schemaName, tableName string) ([]service.TableKey, error) {
	// Получение информации о ключах таблицы
	keyQuery := `SELECT
		tc.CONSTRAINT_NAME,
		tc.CONSTRAINT_TYPE,
		kcu.COLUMN_NAME,
		kcu.REFERENCED_TABLE_SCHEMA AS foreign_table_schema,
		kcu.REFERENCED_TABLE_NAME AS foreign_table_name,
		kcu.REFERENCED_COLUMN_NAME AS foreign_column_name
	FROM
		information_schema.TABLE_CONSTRAINTS AS tc
		JOIN information_schema.KEY_COLUMN_USAGE AS kcu
		  ON tc.CONSTRAINT_NAME = kcu.CONSTRAINT_NAME
		  AND tc.TABLE_SCHEMA = kcu.TABLE_SCHEMA
	WHERE
		tc.CONSTRAINT_TYPE IN ('PRIMARY KEY', 'FOREIGN KEY') AND
		tc.TABLE_SCHEMA = ? AND
		tc.TABLE_NAME = ?`

	rows, err := db.DBConn.Query(keyQuery, schemaName, tableName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var keys []service.TableKey
	for rows.Next() {
		var key service.TableKey
		err := rows.Scan(&key.ConstraintName, &key.ConstraintType, &key.ColumnName,
			&key.ForeignTableSchema, &key.ForeignTableName, &key.ForeignColumnName)
		if err != nil {
			return nil, err
		}
		keys = append(keys, key)
	}

	return keys, nil
}

func (db *DB) GetPrimaryKeys(schemaName, tableName string) ([]string, error) {
	query := `
        SELECT COLUMN_NAME
        FROM information_schema.COLUMNS
        WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? AND COLUMN_KEY = 'PRI';
    `

	rows, err := db.DBConn.Query(query, schemaName, tableName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var primaryKeyColumns []string
	for rows.Next() {
		var columnName string
		if err := rows.Scan(&columnName); err != nil {
			return nil, err
		}
		primaryKeyColumns = append(primaryKeyColumns, columnName)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return primaryKeyColumns, nil
}

func (db *DB) GetForeignKeys(schemaName, tableName string) ([]service.ForeignKey, error) {
	query := `
        SELECT
            kcu.COLUMN_NAME,
            kcu.REFERENCED_TABLE_NAME AS foreign_table_name,
            kcu.REFERENCED_COLUMN_NAME AS foreign_column_name
        FROM
            information_schema.KEY_COLUMN_USAGE AS kcu
        WHERE
            kcu.TABLE_SCHEMA = ? AND
            kcu.TABLE_NAME = ? AND
            kcu.REFERENCED_TABLE_NAME IS NOT NULL;
    `
	rows, err := db.DBConn.Query(query, schemaName, tableName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var foreignKeys []service.ForeignKey
	for rows.Next() {
		var foreignKey service.ForeignKey
		if err := rows.Scan(&foreignKey.Column, &foreignKey.ReferencedTable, &foreignKey.ReferencedColumn); err != nil {
			return nil, err
		}
		foreignKeys = append(foreignKeys, foreignKey)
	}

	return foreignKeys, nil
}

package postgres

import (
	"datasupervision/internal/service"
)

func (db *DB) GetKeys(schemaName, tableName string) ([]service.TableKey, error) {
	// Получение информации о ключах таблицы
	keyQuery := `SELECT
		tc.constraint_name,
		tc.constraint_type,
		kcu.column_name,
		ccu.table_schema AS foreign_table_schema,
		ccu.table_name AS foreign_table_name,
		ccu.column_name AS foreign_column_name
	FROM
		information_schema.table_constraints AS tc
		JOIN information_schema.key_column_usage AS kcu
		  ON tc.constraint_name = kcu.constraint_name
		JOIN information_schema.constraint_column_usage AS ccu
		  ON ccu.constraint_name = tc.constraint_name
	WHERE
		tc.constraint_type IN ('PRIMARY KEY', 'FOREIGN KEY') AND
		tc.table_schema = 'public' AND
		tc.table_name = $1`

	rows, err := db.DBConn.Query(keyQuery, tableName)
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

func (db *DB) GetPrimaryKeys(_, tableName string) ([]string, error) {
	rows, err := db.DBConn.Query(`
		SELECT a.attname
		FROM   pg_index i
		JOIN   pg_attribute a ON a.attrelid = i.indrelid
		                    AND a.attnum = ANY(i.indkey)
		WHERE  i.indrelid = $1::regclass
		AND    i.indisprimary`, tableName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var primaryKeys []string
	for rows.Next() {
		var primaryKey string
		if err := rows.Scan(&primaryKey); err != nil {
			return nil, err
		}
		primaryKeys = append(primaryKeys, primaryKey)
	}

	return primaryKeys, nil
}

func (db *DB) GetForeignKeys(_, tableName string) ([]service.ForeignKey, error) {
	rows, err := db.DBConn.Query(`
		SELECT
			kcu.column_name,
			ccu.table_name AS foreign_table_name,
			ccu.column_name AS foreign_column_name
		FROM
			information_schema.table_constraints AS tc
			JOIN information_schema.key_column_usage AS kcu
			  ON tc.constraint_name = kcu.constraint_name
			JOIN information_schema.constraint_column_usage AS ccu
			  ON ccu.constraint_name = tc.constraint_name
		WHERE constraint_type = 'FOREIGN KEY' AND tc.table_name=$1`, tableName)
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

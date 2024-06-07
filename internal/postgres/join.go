package postgres

import (
	"datasupervision/internal/service"
	"fmt"
)

func (db *DB) Join(schemaName, table1, column1, table2, column2 string) (*service.TableData, error) {
	query := fmt.Sprintf(
		"SELECT * FROM %s.%s as t1 INNER JOIN %s as t2 ON t1.%s = t2.%s",
		schemaName,
		table1, table2,
		column1, column2,
	)

	return db.Select(query)
}

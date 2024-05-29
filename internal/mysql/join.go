package mysql

import (
	"datasupervision/internal/service"
)

func (db *DB) Join(table1, column1, table2, column2 string) (*service.TableData, error) {
	query := "SELECT * FROM ? as t1 INNER JOIN ? as t2 ON t1.? = t2.?"

	return db.Select(query, table1, table2, column1, column2)
}

package mysql

func (db *DB) GetDatabases() ([]string, error) {
	return []string{}, nil
}

func (db *DB) GetSchemas() ([]string, error) {
	rows, err := db.DBConn.Query("SELECT schema_name FROM information_schema.schemata")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	schemas := []string{}

	var schemaName string
	for rows.Next() {
		if err := rows.Scan(&schemaName); err != nil {
			return nil, err
		}

		schemas = append(schemas, schemaName)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return schemas, nil
}

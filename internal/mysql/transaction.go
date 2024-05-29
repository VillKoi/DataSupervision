package mysql

func (db *DB) BeginTransaction() error {
	tx, err := db.DBConn.Begin()
	if err != nil {
		return err
	}

	db.tx = tx
	return nil
}

func (db *DB) Rollback() error {
	err := db.tx.Rollback()
	db.tx = nil
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) Commit() error {
	err := db.tx.Commit()
	db.tx = nil
	if err != nil {
		return err
	}

	return nil
}

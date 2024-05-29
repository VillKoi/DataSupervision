package postgres

func (db *DB) BeginTransaction() error {
	tx, err := db.DBConn.Begin()
	if err != nil {
		return err
	}

	db.tx = tx
	return nil
}

func (db *DB) Rollback() error {
	if db.tx != nil {
		err := db.tx.Rollback()
		db.tx = nil
		if err != nil {
			return err
		}
	}

	return nil
}

func (db *DB) Commit() error {
	if db.tx != nil {
		err := db.tx.Commit()
		db.tx = nil
		if err != nil {
			return err
		}
	}

	return nil
}

package repositories

import (
	"database/sql"
	"fmt"
)

type TxFn func(*sql.Tx) error

type BaseRepository struct {
	db *sql.DB
}

func (r *BaseRepository) WithTransaction(fn TxFn) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("error beginning transaction: %v", err)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // re-throw panic after rollback
		}
	}()

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// Helper methods for common operations
func (r *BaseRepository) execTx(tx *sql.Tx, query string, args ...interface{}) error {
	_, err := tx.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("error executing query: %v", err)
	}
	return nil
}

func (r *BaseRepository) queryRowTx(tx *sql.Tx, query string, args ...interface{}) *sql.Row {
	return tx.QueryRow(query, args...)
}

func (r *BaseRepository) queryTx(tx *sql.Tx, query string, args ...interface{}) (*sql.Rows, error) {
	return tx.Query(query, args...)
}

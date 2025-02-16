package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"smart_campus/internal/domain/repositories"

	"github.com/go-sql-driver/mysql"
)

// BaseRepository provides common functionality for all repositories
type BaseRepository struct {
	conn *Connection
}

// NewBaseRepository creates a new base repository
func NewBaseRepository(conn *Connection) BaseRepository {
	return BaseRepository{conn: conn}
}

// execContext executes a query that doesn't return rows
func (r *BaseRepository) execContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	// Add timeout to context if not already set
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
		defer cancel()
	}

	result, err := r.conn.DB().ExecContext(ctx, query, args...)
	return result, r.handleError(err)
}

// queryRowContext executes a query that returns a single row
func (r *BaseRepository) queryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	if ctx == nil {
		ctx = context.Background()
	}

	// Add timeout to context if not already set
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
		defer cancel()
	}

	return r.conn.DB().QueryRowContext(ctx, query, args...)
}

// queryContext executes a query that returns multiple rows
func (r *BaseRepository) queryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	// Add timeout to context if not already set
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
		defer cancel()
	}

	rows, err := r.conn.DB().QueryContext(ctx, query, args...)
	return rows, r.handleError(err)
}

// beginTx starts a new transaction
func (r *BaseRepository) beginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	return r.conn.BeginTx(ctx, opts)
}

// withTransaction executes a function within a transaction
func (r *BaseRepository) withTransaction(ctx context.Context, fn func(*sql.Tx) error) error {
	if ctx == nil {
		ctx = context.Background()
	}

	tx, err := r.beginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p) // re-throw panic after rollback
		}
	}()

	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx rollback error: %v (original error: %w)", rbErr, err)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}

// handleError converts database errors to domain errors
func (r *BaseRepository) handleError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		return repositories.ErrNotFound
	}

	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		switch mysqlErr.Number {
		case 1062:
			return repositories.ErrDuplicate
		case 1452:
			return repositories.ErrInvalidInput
		case 1451:
			return fmt.Errorf("cannot delete: record is referenced by other records")
		case 1213:
			return fmt.Errorf("deadlock detected")
		case 1205:
			return fmt.Errorf("lock wait timeout exceeded")
		}
	}

	if errors.Is(err, context.DeadlineExceeded) {
		return fmt.Errorf("operation timed out")
	}

	return fmt.Errorf("database error: %w", err)
}

// scanRows scans multiple rows into a slice of interfaces
func (r *BaseRepository) scanRows(rows *sql.Rows, dest interface{}) error {
	if rows == nil {
		return fmt.Errorf("nil rows provided")
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(dest); err != nil {
			return fmt.Errorf("error scanning row: %w", err)
		}
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating rows: %w", err)
	}

	return nil
}

// buildWhereClause builds a WHERE clause from a map of filters
func (r *BaseRepository) buildWhereClause(filters map[string]interface{}) (string, []interface{}) {
	if len(filters) == 0 {
		return "", nil
	}

	var (
		where string
		args  []interface{}
		i     int
	)

	for key, value := range filters {
		if i > 0 {
			where += " AND "
		}
		where += fmt.Sprintf("%s = ?", key)
		args = append(args, value)
		i++
	}

	return "WHERE " + where, args
}

// buildOrderByClause builds an ORDER BY clause
func (r *BaseRepository) buildOrderByClause(orderBy string, orderDir string) string {
	if orderBy == "" {
		return ""
	}

	if orderDir == "" {
		orderDir = "ASC"
	}

	return fmt.Sprintf("ORDER BY %s %s", orderBy, orderDir)
}

// buildLimitClause builds a LIMIT clause
func (r *BaseRepository) buildLimitClause(offset, limit int) string {
	if limit <= 0 {
		return ""
	}
	return fmt.Sprintf("LIMIT %d OFFSET %d", limit, offset)
}

// buildQuery builds a complete query from components
func (r *BaseRepository) buildQuery(base string, filters map[string]interface{}, orderBy, orderDir string, offset, limit int) (string, []interface{}) {
	where, args := r.buildWhereClause(filters)
	order := r.buildOrderByClause(orderBy, orderDir)
	limitClause := r.buildLimitClause(offset, limit)

	query := base
	if where != "" {
		query += " " + where
	}
	if order != "" {
		query += " " + order
	}
	if limitClause != "" {
		query += " " + limitClause
	}

	return query, args
}

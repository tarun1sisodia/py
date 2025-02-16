package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"smart_campus/internal/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

// Config holds the MySQL database configuration
type Config struct {
	Host         string
	Port         int
	User         string
	Password     string
	DatabaseName string
	MaxOpenConns int
	MaxIdleConns int
	MaxLifetime  time.Duration
}

// Connection represents a MySQL database connection
type Connection struct {
	db     *sql.DB
	config *config.Config
	logger *logrus.Logger
}

// NewConnection creates a new database connection
func NewConnection(cfg *config.Config, logger *logrus.Logger) (*Connection, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local&charset=utf8mb4&collation=utf8mb4_unicode_ci&timeout=30s&readTimeout=30s&writeTimeout=30s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Database,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	db.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.Database.MaxLifetime)
	db.SetConnMaxIdleTime(cfg.Database.MaxLifetime / 2)

	// Create connection instance
	conn := &Connection{
		db:     db,
		config: cfg,
		logger: logger,
	}

	// Test connection with context timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := conn.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	// Start health check routine
	conn.StartHealthCheck(30 * time.Second)

	return conn, nil
}

// DB returns the underlying database connection
func (c *Connection) DB() *sql.DB {
	return c.db
}

// Close closes the database connection
func (c *Connection) Close() error {
	return c.db.Close()
}

// BeginTx starts a new transaction with context and options
func (c *Connection) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return c.db.BeginTx(ctx, opts)
}

// PingContext checks the database connection health with context
func (c *Connection) PingContext(ctx context.Context) error {
	return c.db.PingContext(ctx)
}

// Stats returns the database connection statistics
func (c *Connection) Stats() sql.DBStats {
	return c.db.Stats()
}

// LogStats logs the current database connection statistics
func (c *Connection) LogStats() {
	stats := c.db.Stats()
	c.logger.WithFields(logrus.Fields{
		"max_open_connections": c.config.Database.MaxOpenConns,
		"max_idle_connections": c.config.Database.MaxIdleConns,
		"max_lifetime":         c.config.Database.MaxLifetime,
		"open_connections":     stats.OpenConnections,
		"in_use":               stats.InUse,
		"idle":                 stats.Idle,
		"wait_count":           stats.WaitCount,
		"wait_duration":        stats.WaitDuration,
		"max_idle_closed":      stats.MaxIdleClosed,
		"max_lifetime_closed":  stats.MaxLifetimeClosed,
	}).Info("Database connection stats")
}

// StartHealthCheck starts a periodic health check
func (c *Connection) StartHealthCheck(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			if err := c.PingContext(ctx); err != nil {
				c.logger.WithError(err).Error("Database health check failed")
			}
			cancel()
			c.LogStats()
		}
	}()
}

// IsConnected checks if the database is connected
func (c *Connection) IsConnected() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	return c.PingContext(ctx) == nil
}

// GetContext returns a context with timeout for database operations
func (c *Connection) GetContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 30*time.Second)
}

// DefaultConfig returns the default database configuration
func DefaultConfig() *Config {
	return &Config{
		Host:         "localhost",
		Port:         3306,
		MaxOpenConns: 25,
		MaxIdleConns: 5,
		MaxLifetime:  time.Minute * 5,
	}
}

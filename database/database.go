package database

import (
	"database/sql"
	"fmt"
	"net/url"

	_ "github.com/lib/pq"
	"go.nhat.io/otelsql"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

// Config is the required properties to use the database.
type Config struct {
	User     string
	Password string
	Host     string
	Name     string
	// MaxIdleConns int
	// MaxOpenConns int
	DisableTLS bool
}

// Open knows how to open a database connection based on the configuration.
func Open(cfg Config) (*sql.DB, error) {
	sslMode := "require"
	if cfg.DisableTLS {
		sslMode = "disable"
	}

	q := make(url.Values)
	q.Set("sslmode", sslMode)
	q.Set("timezone", "utc")

	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(cfg.User, cfg.Password),
		Host:     cfg.Host,
		Path:     cfg.Name,
		RawQuery: q.Encode(),
	}

	driveName, err := otelsql.Register("postgres",
		otelsql.AllowRoot(),
		otelsql.TraceQueryWithArgs(),
		otelsql.TraceRowsClose(),
		otelsql.WithDatabaseName(cfg.Name),
		otelsql.WithSystem(semconv.DBSystemPostgreSQL),
	)
	if err != nil {
		return nil, fmt.Errorf("registering otelsql driver: %w", err)
	}

	db, err := sql.Open(driveName, u.String())
	if err != nil {
		return nil, fmt.Errorf("opening database connection: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("pinging database: %w", err)
	}

	return db, nil
}

package persistence

import (
    "database/sql"
    "log"
    "sync"

    _ "github.com/go-sql-driver/mysql" // MySQL driver
)

type Store struct {
    db     *sql.DB
    once   sync.Once
    closed bool
}

// NewStore initializes a new Store instance with a database connection.
func NewStore(dataSourceName string) (*Store, error) {
    store := &Store{}
    var err error
    store.once.Do(func() {
        store.db, err = sql.Open("mysql", dataSourceName)
        if err != nil {
            log.Fatalf("Error opening database: %v", err)
        }
    })
    return store, err
}

// Close closes the database connection.
func (s *Store) Close() error {
    if s.closed {
        return nil
    }
    s.closed = true
    return s.db.Close()
}

// BeginTransaction starts a new database transaction.
func (s *Store) BeginTransaction() (*sql.Tx, error) {
    return s.db.Begin()
}

// Execute executes a query without returning any rows.
func (s *Store) Execute(query string, args ...interface{}) (sql.Result, error) {
    return s.db.Exec(query, args...)
}

// QueryRow executes a query that is expected to return a single row.
func (s *Store) QueryRow(query string, args ...interface{}) *sql.Row {
    return s.db.QueryRow(query, args...)
}

// Query executes a query that returns multiple rows.
func (s *Store) Query(query string, args ...interface{}) (*sql.Rows, error) {
    return s.db.Query(query, args...)
}
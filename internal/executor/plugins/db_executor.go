package plugins

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "log"

    _ "github.com/go-sql-driver/mysql" // MySQL driver
    "gotcc/internal/executor"
)

type DBExecutor struct {
    db *sql.DB
}

func NewDBExecutor(dataSourceName string) (*DBExecutor, error) {
    db, err := sql.Open("mysql", dataSourceName)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to database: %w", err)
    }

    return &DBExecutor{db: db}, nil
}

func (e *DBExecutor) Execute(task executor.Task) error {
    var query string
    if err := json.Unmarshal(task.InputData, &query); err != nil {
        return fmt.Errorf("failed to unmarshal input data: %w", err)
    }

    _, err := e.db.Exec(query)
    if err != nil {
        log.Printf("failed to execute query: %s, error: %v", query, err)
        return err
    }

    return nil
}

func (e *DBExecutor) Close() error {
    return e.db.Close()
}
package persistence

import (
    "database/sql"
    "errors"
    "time"

    _ "github.com/go-sql-driver/mysql" // MySQL driver
)

type Repository struct {
    db *sql.DB
}

func NewRepository(dataSourceName string) (*Repository, error) {
    db, err := sql.Open("mysql", dataSourceName)
    if err != nil {
        return nil, err
    }

    // Set connection pool settings
    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(25)
    db.SetConnMaxLifetime(time.Minute * 5)

    return &Repository{db: db}, nil
}

func (r *Repository) Close() error {
    return r.db.Close()
}

// Example methods for task and instance data manipulation

func (r *Repository) CreateTask(taskID string, groupID string, name string, taskType string) error {
    query := "INSERT INTO dist_task (id, group_id, name, type) VALUES (?, ?, ?, ?)"
    _, err := r.db.Exec(query, taskID, groupID, name, taskType)
    return err
}

func (r *Repository) GetTask(taskID string) (string, error) {
    var name string
    query := "SELECT name FROM dist_task WHERE id = ?"
    err := r.db.QueryRow(query, taskID).Scan(&name)
    if err != nil {
        if err == sql.ErrNoRows {
            return "", errors.New("task not found")
        }
        return "", err
    }
    return name, nil
}

func (r *Repository) UpdateTaskStatus(taskID string, status string) error {
    query := "UPDATE dist_task SET status = ? WHERE id = ?"
    _, err := r.db.Exec(query, status, taskID)
    return err
}

// Additional methods for managing task groups and instances can be added here.
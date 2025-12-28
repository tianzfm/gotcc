package dao

import (
    "database/sql"
    "encoding/json"
    "errors"
    "fmt"

    _ "github.com/go-sql-driver/mysql"
)

type TaskDAO struct {
    db *sql.DB
}

func NewTaskDAO(db *sql.DB) *TaskDAO {
    return &TaskDAO{db: db}
}

type DistTask struct {
    ID              string          `json:"id"`
    GroupID         string          `json:"group_id"`
    Name            string          `json:"name"`
    Type            string          `json:"type"`
    Subtype         string          `json:"subtype"`
    Status          string          `json:"status"`
    Priority        int             `json:"priority"`
    Config          json.RawMessage `json:"config"`
    ExecutionContext json.RawMessage `json:"execution_context"`
    InputData       json.RawMessage `json:"input_data"`
    OutputData      json.RawMessage `json:"output_data"`
    StartedAt       sql.NullTime    `json:"started_at"`
    CompletedAt     sql.NullTime    `json:"completed_at"`
    ErrorMessage    sql.NullString   `json:"error_message"`
    ErrorStack      sql.NullString   `json:"error_stack"`
}

func (dao *TaskDAO) CreateTask(task *DistTask) error {
    query := `INSERT INTO dist_task (id, group_id, name, type, subtype, status, priority, config, execution_context, input_data, output_data, started_at, completed_at, error_message, error_stack) 
              VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

    _, err := dao.db.Exec(query, task.ID, task.GroupID, task.Name, task.Type, task.Subtype, task.Status, task.Priority, task.Config, task.ExecutionContext, task.InputData, task.OutputData, task.StartedAt, task.CompletedAt, task.ErrorMessage, task.ErrorStack)
    if err != nil {
        return fmt.Errorf("failed to create task: %w", err)
    }
    return nil
}

func (dao *TaskDAO) GetTaskByID(id string) (*DistTask, error) {
    query := `SELECT id, group_id, name, type, subtype, status, priority, config, execution_context, input_data, output_data, started_at, completed_at, error_message, error_stack 
              FROM dist_task WHERE id = ?`

    row := dao.db.QueryRow(query, id)
    task := &DistTask{}
    err := row.Scan(&task.ID, &task.GroupID, &task.Name, &task.Type, &task.Subtype, &task.Status, &task.Priority, &task.Config, &task.ExecutionContext, &task.InputData, &task.OutputData, &task.StartedAt, &task.CompletedAt, &task.ErrorMessage, &task.ErrorStack)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, errors.New("task not found")
        }
        return nil, fmt.Errorf("failed to get task: %w", err)
    }
    return task, nil
}

func (dao *TaskDAO) UpdateTaskStatus(id string, status string) error {
    query := `UPDATE dist_task SET status = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
    _, err := dao.db.Exec(query, status, id)
    if err != nil {
        return fmt.Errorf("failed to update task status: %w", err)
    }
    return nil
}

func (dao *TaskDAO) DeleteTask(id string) error {
    query := `DELETE FROM dist_task WHERE id = ?`
    _, err := dao.db.Exec(query, id)
    if err != nil {
        return fmt.Errorf("failed to delete task: %w", err)
    }
    return nil
}
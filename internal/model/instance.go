package model

import (
    "time"
)

// TaskGroupInstance represents a transaction instance in the system.
type TaskGroupInstance struct {
    ID          string    `json:"id"`
    FlowID      string    `json:"flow_id"`
    FlowType    string    `json:"flow_type"`
    Status      string    `json:"status"` // pending, running, success, failed, cancelled, rolling_back
    TaskType    string    `json:"task_type"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
    CompletedAt *time.Time `json:"completed_at,omitempty"`
}

// DistTask represents a distributed task associated with a transaction instance.
type DistTask struct {
    ID              string    `json:"id"`
    GroupID         string    `json:"group_id"`
    Name            string    `json:"name"`
    Type            string    `json:"type"` // rpc, local, mq, http, db, file
    Subtype         string    `json:"subtype,omitempty"`
    Status          string    `json:"status"` // pending, running, success, failed, cancelled, rollback_success, rollback_failed
    Priority        int       `json:"priority"`
    Config          interface{} `json:"config"` // JSON configuration
    ExecutionContext interface{} `json:"execution_context"` // Execution context
    InputData       interface{} `json:"input_data"` // Input parameters
    OutputData      interface{} `json:"output_data"` // Output results
    StartedAt       *time.Time `json:"started_at,omitempty"`
    CompletedAt     *time.Time `json:"completed_at,omitempty"`
    ErrorMessage    string     `json:"error_message,omitempty"`
    ErrorStack      string     `json:"error_stack,omitempty"`
}

// ExceptionRecord represents an exception that occurred during task execution.
type ExceptionRecord struct {
    ID             int64     `json:"id"`
    GroupID        string    `json:"group_id"`
    GroupName      string    `json:"group_name"`
    TaskID         string    `json:"task_id"`
    TaskName       string    `json:"task_name"`
    ErrorType      int       `json:"error_type"` // 1: http_error, 2: rpc_error, etc.
    ErrorCode      string    `json:"error_code,omitempty"`
    ErrorMessage   string    `json:"error_message,omitempty"`
    StackTrace     string    `json:"stack_trace,omitempty"`
    OccurredAt     time.Time `json:"occurred_at"`
    Handled        bool      `json:"handled"`
    RetryTimes     int       `json:"retry_times"`
    LastRetryAt    *time.Time `json:"last_retry_at,omitempty"`
}

// ExecutionLog represents a log entry for actions taken on tasks.
type ExecutionLog struct {
    ID        int64     `json:"id"`
    TaskID    string    `json:"task_id"`
    GroupID   string    `json:"group_id"`
    Action    string    `json:"action"` // execute, retry, rollback, cancel, timeout
    OldStatus string    `json:"old_status,omitempty"`
    NewStatus string    `json:"new_status,omitempty"`
    Message   string    `json:"message,omitempty"`
    Details   interface{} `json:"details,omitempty"` // Additional details in JSON format
    CreatedAt time.Time `json:"created_at"`
}
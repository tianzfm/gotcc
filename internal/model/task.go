package model

import (
    "time"
)

type Task struct {
    ID               string          `json:"id"`
    GroupID          string          `json:"group_id"`
    Name             string          `json:"name"`
    Type             string          `json:"type"`
    Subtype          string          `json:"subtype,omitempty"`
    Status           string          `json:"status"`
    Priority         int             `json:"priority"`
    Config           map[string]interface{} `json:"config,omitempty"`
    ExecutionContext map[string]interface{} `json:"execution_context,omitempty"`
    InputData        map[string]interface{} `json:"input_data,omitempty"`
    OutputData       map[string]interface{} `json:"output_data,omitempty"`
    StartedAt        *time.Time      `json:"started_at,omitempty"`
    CompletedAt      *time.Time      `json:"completed_at,omitempty"`
    ErrorMessage     string          `json:"error_message,omitempty"`
    ErrorStack       string          `json:"error_stack,omitempty"`
}
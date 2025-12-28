package model

import (
    "time"
)

// FlowDefinition represents the structure of a transaction flow definition.
type FlowDefinition struct {
    ID          string          `json:"id"`
    Name        string          `json:"name"`
    Description string          `json:"description,omitempty"`
    FlowType    string          `json:"flow_type"`
    Version     int             `json:"version"`
    Definition  map[string]interface{} `json:"definition"`
    IsActive    bool            `json:"is_active"`
    CreatedAt   time.Time       `json:"created_at"`
    UpdatedAt   time.Time       `json:"updated_at"`
    CreateUser  string          `json:"create_user"`
    UpdatedUser string          `json:"updated_user"`
}

// FlowInstance represents the structure of a transaction flow instance.
type FlowInstance struct {
    ID        string    `json:"id"`
    FlowID    string    `json:"flow_id"`
    FlowType  string    `json:"flow_type"`
    Status    string    `json:"status"`
    TaskType  string    `json:"task_type"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    CompletedAt *time.Time `json:"completed_at,omitempty"`
}

// NewFlowDefinition creates a new FlowDefinition instance.
func NewFlowDefinition(id, name, flowType, createUser string, definition map[string]interface{}) *FlowDefinition {
    return &FlowDefinition{
        ID:          id,
        Name:        name,
        FlowType:    flowType,
        Version:     1,
        Definition:  definition,
        IsActive:    true,
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
        CreateUser:  createUser,
        UpdatedUser: createUser,
    }
}

// NewFlowInstance creates a new FlowInstance instance.
func NewFlowInstance(id, flowID, flowType, taskType string) *FlowInstance {
    return &FlowInstance{
        ID:        id,
        FlowID:    flowID,
        FlowType:  flowType,
        Status:    "pending",
        TaskType:  taskType,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
}
package executor

import (
    "context"
    "fmt"
)

// ActionExecutor defines the interface for executing actions in the transaction flow.
type ActionExecutor interface {
    Execute(ctx context.Context, taskID string, config interface{}) (interface{}, error)
}

// Executor manages the execution of tasks using different action executors.
type Executor struct {
    executors map[string]ActionExecutor
}

// NewExecutor creates a new Executor with the provided action executors.
func NewExecutor(executors map[string]ActionExecutor) *Executor {
    return &Executor{
        executors: executors,
    }
}

// ExecuteTask executes a task based on its type and configuration.
func (e *Executor) ExecuteTask(ctx context.Context, taskType string, taskID string, config interface{}) (interface{}, error) {
    executor, exists := e.executors[taskType]
    if !exists {
        return nil, fmt.Errorf("no executor found for task type: %s", taskType)
    }
    return executor.Execute(ctx, taskID, config)
}
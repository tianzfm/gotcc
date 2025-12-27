// models.go

import (
	"encoding/json"
	"time"
)

// 任务状态
type TaskStatus string

const (
	StatusPending   TaskStatus = "PENDING"
	StatusExecuting TaskStatus = "EXECUTING"
	StatusSuccess   TaskStatus = "SUCCESS"
	StatusFailed    TaskStatus = "FAILED"
	StatusRetrying  TaskStatus = "RETRYING"
)

// 重试策略
type RetryPolicy struct {
	MaxAttempts       int     `yaml:"max_attempts" json:"max_attempts"`
	BackoffAlgorithm  string  `yaml:"backoff_algorithm" json:"backoff_algorithm"`
	InitialIntervalMs int     `yaml:"initial_interval_ms" json:"initial_interval_ms"`
	Multiplier        float64 `yaml:"multiplier" json:"multiplier"`
	MaxIntervalMs     int     `yaml:"max_interval_ms" json:"max_interval_ms"`
	IntervalMs        int     `yaml:"interval_ms" json:"interval_ms"`
}

// Action定义
type ActionConfig struct {
	ID              string                 `yaml:"id" json:"id"`
	Type            string                 `yaml:"type" json:"type"`
	Name            string                 `yaml:"name" json:"name"`
	Config          map[string]interface{} `yaml:"config" json:"config"`
	RequestTemplate string                 `yaml:"request_template" json:"request_template"`
	RetryPolicy     string                 `yaml:"retry_policy" json:"retry_policy"`
	Compensation    *ActionConfig          `yaml:"compensation" json:"compensation"`
}

// 事务流程定义
type TransactionFlow struct {
	ID            string                 `yaml:"id" json:"id"`
	Name          string                 `yaml:"name" json:"name"`
	Actions       []ActionConfig         `yaml:"actions" json:"actions"`
	RetryPolicies map[string]RetryPolicy `yaml:"retry_policies" json:"retry_policies"`
}

// 任务实例
type TaskInstance struct {
	ID            string     `json:"id"`
	TransactionID string     `json:"transaction_id"`
	ActionID      string     `json:"action_id"`
	ActionName    string     `json:"action_name"`
	Status        TaskStatus `json:"status"`
	TryTimes      int        `json:"try_times"`
	MaxTryTimes   int        `json:"max_try_times"`
	InputParams   string     `json:"input_params"`
	OutputResult  string     `json:"output_result"`
	ErrorMsg      string     `json:"error_msg"`
	NextRetryTime time.Time  `json:"next_retry_time"`
	CreatedTime   time.Time  `json:"created_time"`
	ModifiedTime  time.Time  `json:"modified_time"`
}
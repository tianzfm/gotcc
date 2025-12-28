package config

type TaskGroupFlow struct {
    ID          string          `json:"id" yaml:"id"`
    Name        string          `json:"name" yaml:"name"`
    Description *string         `json:"description,omitempty" yaml:"description,omitempty"`
    FlowType    string          `json:"flow_type" yaml:"flow_type"`
    Version     int             `json:"version" yaml:"version"`
    Definition  map[string]interface{} `json:"definition" yaml:"definition"`
    IsActive    bool            `json:"is_active" yaml:"is_active"`
    CreatedAt   string          `json:"created_at" yaml:"created_at"`
    UpdatedAt   string          `json:"updated_at" yaml:"updated_at"`
    CreateUser  string          `json:"create_user" yaml:"create_user"`
    UpdatedUser string          `json:"updated_user" yaml:"updated_user"`
}

type TaskGroupInstance struct {
    ID        string `json:"id" yaml:"id"`
    FlowID    string `json:"flow_id" yaml:"flow_id"`
    FlowType  string `json:"flow_type" yaml:"flow_type"`
    Status    string `json:"status" yaml:"status"`
    TaskType  string `json:"task_type" yaml:"task_type"`
    CreatedAt string `json:"created_at" yaml:"created_at"`
    UpdatedAt string `json:"updated_at" yaml:"updated_at"`
    CompletedAt *string `json:"completed_at,omitempty" yaml:"completed_at,omitempty"`
}

type DistTask struct {
    ID             string `json:"id" yaml:"id"`
    GroupID        string `json:"group_id" yaml:"group_id"`
    Name           string `json:"name" yaml:"name"`
    Type           string `json:"type" yaml:"type"`
    Subtype        *string `json:"subtype,omitempty" yaml:"subtype,omitempty"`
    Status         string `json:"status" yaml:"status"`
    Priority       int    `json:"priority" yaml:"priority"`
    Config         map[string]interface{} `json:"config" yaml:"config"`
    ExecutionContext map[string]interface{} `json:"execution_context" yaml:"execution_context"`
    InputData      map[string]interface{} `json:"input_data" yaml:"input_data"`
    OutputData     map[string]interface{} `json:"output_data" yaml:"output_data"`
    StartedAt      *string `json:"started_at,omitempty" yaml:"started_at,omitempty"`
    CompletedAt    *string `json:"completed_at,omitempty" yaml:"completed_at,omitempty"`
    ErrorMessage   *string `json:"error_message,omitempty" yaml:"error_message,omitempty"`
    ErrorStack     *string `json:"error_stack,omitempty" yaml:"error_stack,omitempty"`
}

type RetryPolicyConfig struct {
    ID               int     `json:"id" yaml:"id"`
    Name             string  `json:"name" yaml:"name"`
    PolicyType       string  `json:"policy_type" yaml:"policy_type"`
    BaseInterval     int     `json:"base_interval" yaml:"base_interval"`
    MaxInterval      int     `json:"max_interval" yaml:"max_interval"`
    MaxAttempts      int     `json:"max_attempts" yaml:"max_attempts"`
    Multiplier       float64 `json:"multiplier" yaml:"multiplier"`
    RandomizationFactor float64 `json:"randomization_factor" yaml:"randomization_factor"`
    Enabled          bool    `json:"enabled" yaml:"enabled"`
    Config           map[string]interface{} `json:"config" yaml:"config"`
}

type ExceptionRecord struct {
    ID            int64   `json:"id" yaml:"id"`
    GroupID       string  `json:"group_id" yaml:"group_id"`
    GroupName     string  `json:"group_name" yaml:"group_name"`
    TaskID        string  `json:"task_id" yaml:"task_id"`
    TaskName      string  `json:"task_name" yaml:"task_name"`
    ErrorType     int     `json:"error_type" yaml:"error_type"`
    ErrorCode     *string `json:"error_code,omitempty" yaml:"error_code,omitempty"`
    ErrorMessage  *string `json:"error_message,omitempty" yaml:"error_message,omitempty"`
    StackTrace    *string `json:"stack_trace,omitempty" yaml:"stack_trace,omitempty"`
    OccurredAt    string  `json:"occurred_at" yaml:"occurred_at"`
    Handled       bool    `json:"handled" yaml:"handled"`
    RetryTimes    int     `json:"retry_times" yaml:"retry_times"`
    LastRetryAt   *string `json:"last_retry_at,omitempty" yaml:"last_retry_at,omitempty"`
}

type ExecutionLog struct {
    ID        int64   `json:"id" yaml:"id"`
    TaskID    string  `json:"task_id" yaml:"task_id"`
    GroupID   string  `json:"group_id" yaml:"group_id"`
    Action    string  `json:"action" yaml:"action"`
    OldStatus *string `json:"old_status,omitempty" yaml:"old_status,omitempty"`
    NewStatus *string `json:"new_status,omitempty" yaml:"new_status,omitempty"`
    Message   *string `json:"message,omitempty" yaml:"message,omitempty"`
    Details   map[string]interface{} `json:"details" yaml:"details"`
    CreatedAt string  `json:"created_at" yaml:"created_at"`
}
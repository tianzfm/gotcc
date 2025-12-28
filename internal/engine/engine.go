package engine

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"gopkg.in/yaml.v2"

	"gotcc/common"
	"register"
)

// Minimal flow / policy / task models used by engine.
type RetryPolicy struct {
	MaxAttempts       int     `yaml:"maxAttempts" json:"maxAttempts"`
	BackoffAlgorithm  string  `yaml:"backoff" json:"backoff"`
	InitialIntervalMs int     `yaml:"initialIntervalMs" json:"initialIntervalMs"`
	Multiplier        float64 `yaml:"multiplier" json:"multiplier"`
	MaxIntervalMs     int     `yaml:"maxIntervalMs" json:"maxIntervalMs"`
	IntervalMs        int     `yaml:"intervalMs" json:"intervalMs"`
}

type ActionConfig struct {
	ID              string                 `yaml:"id" json:"id"`
	Name            string                 `yaml:"name" json:"name"`
	Type            string                 `yaml:"type" json:"type"`
	ComponentID     string                 `yaml:"component" json:"component"`
	Config          map[string]interface{} `yaml:"config" json:"config"`
	RetryPolicy     string                 `yaml:"retryPolicy" json:"retryPolicy"`
	RequestTemplate string                 `yaml:"requestTemplate" json:"requestTemplate"`
}

type TransactionFlow struct {
	Name          string                 `yaml:"name" json:"name"`
	Actions       []ActionConfig         `yaml:"actions" json:"actions"`
	RetryPolicies map[string]RetryPolicy `yaml:"retryPolicies" json:"retryPolicies"`
}

type TaskInstance struct {
	ID         string
	ActionID   string
	ActionName string
	Component  register.TaskComponent
	Input      json.RawMessage
	Output     json.RawMessage
	TryTimes   int
	Status     string
	ErrorMsg   string
}

// FlowLoader loads a TransactionFlow by ID (e.g., from DB).
type FlowLoader func(ctx context.Context, flowID string) (*TransactionFlow, error)

// Engine is a higher-level wrapper around executing flows.
type Engine struct {
	loader FlowLoader
}

// NewEngine creates an Engine using the provided FlowLoader.
func NewEngine(loader FlowLoader) (*Engine, error) {
	if loader == nil {
		return nil, errors.New("flow loader required")
	}
	return &Engine{loader: loader}, nil
}

// ExecuteTransaction loads flow by id, builds components and runs Try/Confirm/Cancel.
func (e *Engine) ExecuteTransaction(ctx context.Context, flowID string, params map[string]interface{}) (string, error) {
	// load flow
	flow, err := e.loader(ctx, flowID)
	if err != nil {
		return "", fmt.Errorf("load flow failed: %w", err)
	}

	// build task instances
	payload, _ := json.Marshal(params)
	var tasks []*TaskInstance
	for _, a := range flow.Actions {
		cfgBytes, _ := json.Marshal(a.Config)
		comp, err := register.BuildComponent(a.ComponentID, cfgBytes)
		if err != nil {
			return "", fmt.Errorf("build component %s failed: %w", a.ComponentID, err)
		}
		if err := comp.Prepare(cfgBytes); err != nil {
			return "", fmt.Errorf("component prepare failed: %w", err)
		}
		if err := comp.Validate(); err != nil {
			return "", fmt.Errorf("component validate failed: %w", err)
		}
		t := &TaskInstance{ID: common.GenerateTransactionID(), ActionID: a.ID, ActionName: a.Name, Component: comp, Input: payload, Status: "pending"}
		tasks = append(tasks, t)
	}

	txID := common.GenerateTransactionID()
	log.Printf("Transaction %s started flow=%s tasks=%d", txID, flow.Name, len(tasks))

	tried := make([]*TaskInstance, 0, len(tasks))
	for _, t := range tasks {
		if err := e.tryWithRetry(ctx, flow, txID, t); err != nil {
			log.Printf("task Try failed: %s, triggering compensation", t.ActionName)
			e.compensate(ctx, txID, reverseTasks(tried))
			return txID, err
		}
		tried = append(tried, t)
	}

	for _, t := range tasks {
		if _, err := t.Component.Confirm(ctx, txID); err != nil {
			log.Printf("Confirm failed for %s: %v", t.ActionName, err)
			e.compensate(ctx, txID, reverseTasks(tasks))
			return txID, err
		}
	}
	log.Printf("Transaction %s completed successfully", txID)
	return txID, nil
}

func (e *Engine) tryWithRetry(ctx context.Context, flow *TransactionFlow, txID string, t *TaskInstance) error {
	policy := e.resolveRetryPolicyForAction(flow, t.ActionID)
	max := 1
	if policy != nil && policy.MaxAttempts > 0 {
		max = policy.MaxAttempts
	}
	var lastErr error
	for attempt := 0; attempt < max; attempt++ {
		t.TryTimes = attempt + 1
		req := &register.TCCReq{TxID: txID, Payload: t.Input}
		resp, err := t.Component.Try(ctx, req)
		if err == nil && resp != nil && resp.Success {
			if resp.Result != nil {
				t.Output = resp.Result
			}
			t.Status = "success"
			return nil
		}
		if err == nil && resp != nil {
			lastErr = errors.New(resp.Error)
		} else {
			lastErr = err
		}
		t.ErrorMsg = lastErr.Error()
		t.Status = "retrying"
		if attempt == max-1 {
			t.Status = "failed"
			return fmt.Errorf("task %s try failed after %d attempts: %w", t.ActionName, attempt+1, lastErr)
		}
		delay := time.Duration(1000) * time.Millisecond
		if policy != nil {
			delay = calculateDelay(*policy, attempt)
		}
		time.Sleep(delay)
	}
	return lastErr
}

func (e *Engine) compensate(ctx context.Context, txID string, tasks []*TaskInstance) {
	for _, t := range tasks {
		if !t.Component.SupportsCompensation() {
			log.Printf("component %s does not support compensation, task %s requires manual intervention", t.Component.ID(), t.ActionName)
			continue
		}
		if _, err := t.Component.Cancel(ctx, txID); err != nil {
			log.Printf("compensation failed task=%s err=%v", t.ActionName, err)
		} else {
			log.Printf("compensation success task=%s", t.ActionName)
		}
	}
}

func reverseTasks(in []*TaskInstance) []*TaskInstance {
	out := make([]*TaskInstance, 0, len(in))
	for i := len(in) - 1; i >= 0; i-- {
		out = append(out, in[i])
	}
	return out
}

func (e *Engine) resolveRetryPolicyForAction(flow *TransactionFlow, actionID string) *RetryPolicy {
	for _, a := range flow.Actions {
		if a.ID == actionID {
			if rp, ok := flow.RetryPolicies[a.RetryPolicy]; ok {
				return &rp
			}
			break
		}
	}
	return nil
}

func calculateDelay(p RetryPolicy, attempt int) time.Duration {
	switch p.BackoffAlgorithm {
	case "exponential":
		delay := float64(p.InitialIntervalMs) * pow(p.Multiplier, attempt)
		if max := float64(p.MaxIntervalMs); delay > max {
			delay = max
		}
		return time.Duration(delay) * time.Millisecond
	case "fixed":
		return time.Duration(p.IntervalMs) * time.Millisecond
	default:
		return time.Duration(p.InitialIntervalMs) * time.Millisecond
	}
}

func pow(base float64, exp int) float64 {
	if exp == 0 {
		return 1
	}
	res := 1.0
	for i := 0; i < exp; i++ {
		res *= base
	}
	return res
}

// SQLFlowLoader returns a FlowLoader that reads `defination` from task_group_flow table.
func SQLFlowLoader(db *sql.DB) FlowLoader {
	return func(ctx context.Context, flowID string) (*TransactionFlow, error) {
		var def string
		row := db.QueryRowContext(ctx, "SELECT defination FROM task_group_flow WHERE id = ?", flowID)
		if err := row.Scan(&def); err != nil {
			return nil, err
		}
		var flow TransactionFlow
		if err := json.Unmarshal([]byte(def), &flow); err == nil {
			return &flow, nil
		}
		if err := yaml.Unmarshal([]byte(def), &flow); err == nil {
			return &flow, nil
		}
		return nil, fmt.Errorf("cannot parse flow definition for id %s", flowID)
	}
}

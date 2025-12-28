// engine.go
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"time"

	"register"
)

// Minimal flow / policy / task models used by engine.
// If your project already defines these, keep those and remove duplicates.
type RetryPolicy struct {
	MaxAttempts       int     `yaml:"maxAttempts" json:"maxAttempts"`
	BackoffAlgorithm  string  `yaml:"backoff" json:"backoff"` // "exponential"|"fixed"
	InitialIntervalMs int     `yaml:"initialIntervalMs" json:"initialIntervalMs"`
	Multiplier        float64 `yaml:"multiplier" json:"multiplier"`
	MaxIntervalMs     int     `yaml:"maxIntervalMs" json:"maxIntervalMs"`
	IntervalMs        int     `yaml:"intervalMs" json:"intervalMs"`
}

type ActionConfig struct {
	ID          string                 `yaml:"id" json:"id"`
	Name        string                 `yaml:"name" json:"name"`
	Type        string                 `yaml:"type" json:"type"`
	Component   string                 `yaml:"component" json:"component"`
	Config      map[string]interface{} `yaml:"config" json:"config"`
	RetryPolicy struct {
		ID              int    `yaml:"id" json:"id"`
		MaxAttempts     int    `yaml:"max_attempts" json:"max_attempts"`
		BackoffStrategy string `yaml:"backoff_strategy" json:"backoff_strategy"`
	} `yaml:"retry_policy" json:"retry_policy"`
	RequestTemplate string `yaml:"requestTemplate" json:"requestTemplate"`
}

type TransactionFlow struct {
	Name          string                 `yaml:"name" json:"name"`
	Actions       []ActionConfig         `yaml:"actions" json:"actions"`
	RetryPolicies map[string]RetryPolicy `yaml:"retryPolicies" json:"retryPolicies"`
}

type TaskInstance struct {
	ID          string
	ActionID    string
	ActionName  string
	Component   register.TaskComponent
	Input       json.RawMessage
	Output      json.RawMessage
	TryTimes    int
	Status      string
	NextRetryAt time.Time
	ErrorMsg    string
}

// DistributedEngine is component-driven: it builds TaskComponents from action configs
// then runs Try for every action, and on success runs Confirm for all; on failure runs Cancel for already-tryed ones.
type FlowLoader func(ctx context.Context, flowID string) (*TransactionFlow, error)

type DistributedEngine struct {
	// FlowLoader loads a TransactionFlow by its ID from persistence.
	FlowLoader FlowLoader
	httpClient *HTTPClient
	taskStore  map[string]*TaskInstance
}

// NewDistributedEngine loads a yaml flow file and returns an engine.
// NewDistributedEngine creates an engine that uses the provided FlowLoader to
// load flows by ID at execution time.
func NewDistributedEngine(loader FlowLoader) (*DistributedEngine, error) {
	if loader == nil {
		return nil, errors.New("flow loader required")
	}
	e := &DistributedEngine{
		FlowLoader: loader,
		taskStore:  make(map[string]*TaskInstance),
		httpClient: NewHTTPClient(),
	}
	return e, nil
}

// ExecuteTransaction drives a single transaction identified by txID (generated if empty).
// params is serialized and passed as payload to Try.
func (e *DistributedEngine) ExecuteTransaction(ctx context.Context, flowID string, params map[string]interface{}) (string, error) {
	txID := common.GenerateTransactionID()
	// load flow by id
	flow, err := e.FlowLoader(ctx, flowID)
	if err != nil {
		return txID, fmt.Errorf("load flow %s failed: %w", flowID, err)
	}

	// build task instances (components + inputs)
	tasks, err := e.buildTasks(flow, params)
	if err != nil {
		return txID, err
	}

	log.Printf("Transaction %s started flow=%s tasks=%d", txID, flow.Name, len(tasks))

	// 1) Try phase: sequential here (can be parallel based on task config)
	tried := make([]*TaskInstance, 0, len(tasks))
	for _, t := range tasks {
		if err := e.tryWithRetry(ctx, txID, flow, t); err != nil {
			// Try failed after retries -> trigger compensation for those succeeded
			log.Printf("task Try failed: %s, triggering compensation", t.ActionName)
			e.compensate(ctx, txID, reverseTasks(tried))
			return txID, fmt.Errorf("transaction failed: %w", err)
		}
		tried = append(tried, t)
	}

	// 2) Confirm phase: if all Try succeeded, confirm all (in original order)
	for _, t := range tasks {
		if _, err := t.Component.Confirm(ctx, txID); err != nil {
			// Confirm failed -> this is serious: try to Cancel where possible and alert
			log.Printf("Confirm failed for %s: %v", t.ActionName, err)
			e.compensate(ctx, txID, reverseTasks(tasks))
			return txID, fmt.Errorf("confirm failed: %w", err)
		}
	}

	log.Printf("Transaction %s completed successfully", txID)
	return txID, nil
}

// buildTasks constructs TaskInstances by instantiating registered components.
func (e *DistributedEngine) buildTasks(flow *TransactionFlow, params map[string]interface{}) ([]*TaskInstance, error) {
	// serialize params once
	payload, _ := json.Marshal(params)
	var tasks []*TaskInstance
	for _, a := range flow.Actions {
		// marshal action config to json.RawMessage for factory
		cfgBytes, _ := json.Marshal(a.Config)
		comp, err := register.BuildComponent(a.Component, cfgBytes)
		if err != nil {
			return nil, fmt.Errorf("无法构建组件 %s: %w", a.Component, err)
		}
		// Prepare and Validate lifecycle
		if err := comp.Prepare(cfgBytes); err != nil {
			return nil, fmt.Errorf("组件 Prepare 失败: %w", err)
		}
		if err := comp.Validate(); err != nil {
			return nil, fmt.Errorf("组件 Validate 失败: %w", err)
		}

		t := &TaskInstance{
			ID:         GenerateTransactionID(), // task id; replace by better id if needed
			ActionID:   a.ID,
			ActionName: a.Name,
			Component:  comp,
			Input:      payload,
			Status:     "pending",
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

// tryWithRetry executes Try on a task with retry policy.
func (e *DistributedEngine) tryWithRetry(ctx context.Context, txID string, flow *TransactionFlow, t *TaskInstance) error {
	policy := e.resolveRetryPolicyForAction(flow, t.ActionID)
	max := 1
	if policy != nil && policy.MaxAttempts > 0 {
		max = policy.MaxAttempts
	}

	var lastErr error
	for attempt := 0; attempt < max; attempt++ {
		t.TryTimes = attempt + 1
		req := &register.TCCReq{
			TxID:    txID,
			Payload: t.Input,
		}
		resp, err := t.Component.Try(ctx, req)
		if err == nil && resp != nil && resp.Success {
			// success
			if resp.Result != nil {
				t.Output = resp.Result
			}
			t.Status = "success"
			return nil
		}
		// failure path
		if err == nil && resp != nil {
			lastErr = errors.New(resp.Error)
		} else {
			lastErr = err
		}
		t.ErrorMsg = lastErr.Error()
		t.Status = "retrying"
		// if last attempt, mark failed
		if attempt == max-1 {
			t.Status = "failed"
			return fmt.Errorf("task %s try failed after %d attempts: %w", t.ActionName, attempt+1, lastErr)
		}
		// wait according to backoff
		delay := time.Duration(1000) * time.Millisecond
		if policy != nil {
			delay = calculateDelay(*policy, attempt)
		}
		t.NextRetryAt = time.Now().Add(delay)
		log.Printf("task %s will retry after %v (attempt %d/%d)", t.ActionName, delay, attempt+1, max)
		select {
		case <-time.After(delay):
			// retry
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	return lastErr
}

// compensate runs Cancel on tasks (in reverse order). If a component doesn't support compensation it logs an alert.
func (e *DistributedEngine) compensate(ctx context.Context, txID string, tasks []*TaskInstance) {
	for _, t := range tasks {
		if !t.Component.SupportsCompensation() {
			log.Printf("组件 %s 不支持补偿, 任务 %s 需要人工介入", t.Component.ID(), t.ActionName)
			// here you could push an alert to monitoring/notification
			continue
		}
		if _, err := t.Component.Cancel(ctx, txID); err != nil {
			log.Printf("补偿失败 task=%s err=%v", t.ActionName, err)
			// optionally retry cancellation or escalate
		} else {
			log.Printf("补偿成功 task=%s", t.ActionName)
		}
	}
}

// helpers

func reverseTasks(in []*TaskInstance) []*TaskInstance {
	out := make([]*TaskInstance, 0, len(in))
	for i := len(in) - 1; i >= 0; i-- {
		out = append(out, in[i])
	}
	return out
}

func (e *DistributedEngine) resolveRetryPolicyForAction(flow *TransactionFlow, actionID string) *RetryPolicy {
	// find action config
	for _, a := range flow.Actions {
		if a.ID == actionID {
			// prefer inline retry_policy on action
			if a.RetryPolicy.MaxAttempts > 0 {
				rp := &RetryPolicy{}
				rp.MaxAttempts = a.RetryPolicy.MaxAttempts
				switch a.RetryPolicy.BackoffStrategy {
				case "exponential":
					rp.BackoffAlgorithm = "exponential"
					rp.InitialIntervalMs = 1000
					rp.Multiplier = 2
					rp.MaxIntervalMs = 30000
				case "fixed":
					rp.BackoffAlgorithm = "fixed"
					rp.IntervalMs = 1000
				default:
					rp.BackoffAlgorithm = "fixed"
					rp.IntervalMs = 1000
				}
				return rp
			}
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
		delay := float64(p.InitialIntervalMs) * math.Pow(p.Multiplier, float64(attempt))
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

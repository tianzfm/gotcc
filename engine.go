// engine.go

import (
	"context"
	"fmt"
	"log"
	"math"
	"time"
	"gopkg.in/yaml.v2"
)

type DistributedEngine struct {
	workflowConfig *TransactionFlow
	taskStore      map[string]*TaskInstance
	httpClient     *HTTPClient
}

func NewDistributedEngine(configPath string) (*DistributedEngine, error) {
	engine := &DistributedEngine{
		taskStore:  make(map[string]*TaskInstance),
		httpClient: NewHTTPClient(),
	}

	// 加载配置
	if err := engine.loadConfig(configPath); err != nil {
		return nil, err
	}

	return engine, nil
}

// 加载工作流配置
func (e *DistributedEngine) loadConfig(configPath string) error {
	data, err := ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("读取配置文件失败: %v", err)
	}

	var flow struct {
		TransactionFlow TransactionFlow `yaml:"transactionFlow"`
	}

	if err := yaml.Unmarshal(data, &flow); err != nil {
		return fmt.Errorf("解析YAML配置失败: %v", err)
	}

	e.workflowConfig = &flow.TransactionFlow
	return nil
}

// 执行事务流程
func (e *DistributedEngine) ExecuteTransaction(ctx context.Context, params map[string]interface{}) (string, error) {
	transactionID := GenerateTransactionID()

	log.Printf("开始执行事务流程: %s, ID: %s", e.workflowConfig.Name, transactionID)

	// 创建任务实例
	tasks := e.createTaskInstances(transactionID, params)

	// 顺序执行所有Action
	for i, task := range tasks {
		if err := e.executeActionWithRetry(ctx, task); err != nil {
			log.Printf("任务执行失败: %s, 开始回滚", task.ActionName)

			// 执行补偿操作
			e.executeCompensation(ctx, tasks[:i])

			return transactionID, fmt.Errorf("事务执行失败: %v", err)
		}
	}

	log.Printf("事务流程执行完成: %s", transactionID)
	return transactionID, nil
}

// 带重试的任务执行
func (e *DistributedEngine) executeActionWithRetry(ctx context.Context, task *TaskInstance) error {
	policyName := e.getActionRetryPolicy(task.ActionID)
	policy := e.workflowConfig.RetryPolicies[policyName]

	for attempt := 0; attempt < policy.MaxAttempts; attempt++ {
		task.TryTimes = attempt + 1
		task.Status = StatusExecuting
		task.ModifiedTime = time.Now()

		log.Printf("执行任务: %s, 尝试次数: %d/%d", task.ActionName, attempt+1, policy.MaxAttempts)

		// 执行具体任务
		result, err := e.executeSingleAction(ctx, task)

		if err == nil {
			// 执行成功
			task.Status = StatusSuccess
			task.OutputResult = result
			log.Printf("任务执行成功: %s", task.ActionName)
			return nil
		}

		// 执行失败
		task.ErrorMsg = err.Error()

		if attempt == policy.MaxAttempts-1 {
			// 最后一次尝试也失败
			task.Status = StatusFailed
			log.Printf("任务最终失败: %s, 错误: %v", task.ActionName, err)
			return err
		}

		// 计算下次重试时间
		delay := e.calculateRetryDelay(policy, attempt)
		task.NextRetryTime = time.Now().Add(delay)
		task.Status = StatusRetrying

		log.Printf("任务执行失败: %s, %v后重试", task.ActionName, delay)
		time.Sleep(delay)
	}

	return nil
}

// 计算重试延迟
func (e *DistributedEngine) calculateRetryDelay(policy RetryPolicy, attempt int) time.Duration {
	switch policy.BackoffAlgorithm {
	case "exponential":
		delay := float64(policy.InitialIntervalMs) * math.Pow(policy.Multiplier, float64(attempt))
		if max := float64(policy.MaxIntervalMs); delay > max {
			delay = max
		}
		return time.Duration(delay) * time.Millisecond

	case "fixed":
		return time.Duration(policy.IntervalMs) * time.Millisecond

	default:
		return time.Duration(policy.InitialIntervalMs) * time.Millisecond
	}
}

// 执行单个Action
func (e *DistributedEngine) executeSingleAction(ctx context.Context, task *TaskInstance) (string, error) {
	actionConfig := e.getActionConfig(task.ActionID)

	switch actionConfig.Type {
	case "HTTP_RPC":
		return e.executeHTTPAction(ctx, actionConfig, task.InputParams)
	default:
		return "", fmt.Errorf("不支持的Action类型: %s", actionConfig.Type)
	}
}

// 执行HTTP Action
func (e *DistributedEngine) executeHTTPAction(ctx context.Context, action *ActionConfig, params string) (string, error) {
	endpoint, _ := action.Config["endpoint"].(string)
	method, _ := action.Config["method"].(string)

	// 这里简化处理，实际应该渲染模板
	requestBody := action.RequestTemplate

	result, err := e.httpClient.DoRequest(method, endpoint, requestBody)
	if err != nil {
		return "", fmt.Errorf("HTTP请求失败: %v", err)
	}

	return result, nil
}
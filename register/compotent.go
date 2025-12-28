package register

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
)

// TCC 请求/响应通用结构
type TCCReq struct {
	TxID     string                 `json:"tx_id,omitempty"`
	Payload  json.RawMessage        `json:"payload,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

type TCCResp struct {
	Success bool            `json:"success"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   string          `json:"error,omitempty"`
}

// TaskComponent 是任务执行器组件接口，支持从配置构建、校验和三阶段（Try/Confirm/Cancel）操作。
// 设计目标：
// - 组件可由工厂根据 task 配置实例化（便于把 group 中的每个 task 组装为组件实例）
// - 提供能力声明（是否支持补偿/幂等/并发）
// - 提供生命周期钩子（Prepare/Validate）以解析 task 的 config
type TaskComponent interface {
	// 唯一标识（组件类型），如 "http", "db", "mq" 等
	ID() string
	// 从 task 的 json/yaml 配置构建或准备组件实例
	Prepare(cfg json.RawMessage) error
	// 校验配置（可在 Prepare 后调用）
	Validate() error

	// Try 阶段：预留资源/模拟/锁定。返回结果供 Confirm/Cancel 使用。
	Try(ctx context.Context, req *TCCReq) (*TCCResp, error)
	// Confirm 阶段：最终提交
	Confirm(ctx context.Context, txID string) (*TCCResp, error)
	// Cancel 阶段：回滚/释放资源
	Cancel(ctx context.Context, txID string) (*TCCResp, error)

	// 支持补偿（Cancel）吗？若不支持，引擎应记录并报警人工介入
	SupportsCompensation() bool
	// 是否是幂等操作（对重试友好）
	IsIdempotent() bool
}

// Factory 用于根据 task 的配置构建组件实例
type ComponentFactory func(taskCfg json.RawMessage) (TaskComponent, error)

// Registry 管理组件工厂与实例构建
type Registry struct {
	mu        sync.RWMutex
	factories map[string]ComponentFactory
}

// DefaultRegistry 全局注册器
var DefaultRegistry = NewRegistry()

func NewRegistry() *Registry {
	return &Registry{
		factories: make(map[string]ComponentFactory),
	}
}

// Register 注册一个组件工厂（按组件 ID）
// 若已存在会返回错误
func (r *Registry) Register(id string, factory ComponentFactory) error {
	if id == "" || factory == nil {
		return errors.New("invalid id or factory")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.factories[id]; ok {
		return errors.New("factory already registered: " + id)
	}
	r.factories[id] = factory
	return nil
}

// MustRegister 注册但在冲突时 panic（便于 init 注册使用）
func (r *Registry) MustRegister(id string, factory ComponentFactory) {
	if err := r.Register(id, factory); err != nil {
		panic(err)
	}
}

// GetFactory 获取工厂
func (r *Registry) GetFactory(id string) (ComponentFactory, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	f, ok := r.factories[id]
	return f, ok
}

// Build 使用已注册的工厂构建组件实例
func (r *Registry) Build(id string, cfg json.RawMessage) (TaskComponent, error) {
	f, ok := r.GetFactory(id)
	if !ok {
		return nil, errors.New("component factory not found: " + id)
	}
	return f(cfg)
}

// ListIDs 返回已注册组件的 id 列表
func (r *Registry) ListIDs() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ids := make([]string, 0, len(r.factories))
	for id := range r.factories {
		ids = append(ids, id)
	}
	return ids
}

// 便捷全局函数
func Register(id string, factory ComponentFactory) error {
	return DefaultRegistry.Register(id, factory)
}
func MustRegister(id string, factory ComponentFactory) { DefaultRegistry.MustRegister(id, factory) }
func BuildComponent(id string, cfg json.RawMessage) (TaskComponent, error) {
	return DefaultRegistry.Build(id, cfg)
}

package engine

import (
    "context"
    "sync"
    "time"
)

type Scheduler struct {
    mu          sync.Mutex
    tasks       map[string]*Task
    retryPolicy *RetryPolicy
}

type Task struct {
    ID        string
    Action    func() error
    Status    string
    CreatedAt time.Time
}

type RetryPolicy struct {
    MaxAttempts int
    Delay       time.Duration
}

func NewScheduler(retryPolicy *RetryPolicy) *Scheduler {
    return &Scheduler{
        tasks:       make(map[string]*Task),
        retryPolicy: retryPolicy,
    }
}

func (s *Scheduler) Schedule(taskID string, action func() error) {
    s.mu.Lock()
    defer s.mu.Unlock()

    task := &Task{
        ID:        taskID,
        Action:    action,
        Status:    "pending",
        CreatedAt: time.Now(),
    }
    s.tasks[taskID] = task
}

func (s *Scheduler) Execute(ctx context.Context, taskID string) error {
    s.mu.Lock()
    task, exists := s.tasks[taskID]
    s.mu.Unlock()

    if !exists {
        return nil // Task not found
    }

    for attempt := 1; attempt <= s.retryPolicy.MaxAttempts; attempt++ {
        err := task.Action()
        if err == nil {
            task.Status = "success"
            return nil
        }

        task.Status = "failed"
        time.Sleep(s.retryPolicy.Delay)
    }

    return nil
}

func (s *Scheduler) Cancel(taskID string) {
    s.mu.Lock()
    defer s.mu.Unlock()

    if task, exists := s.tasks[taskID]; exists {
        task.Status = "cancelled"
    }
}

func (s *Scheduler) GetTaskStatus(taskID string) string {
    s.mu.Lock()
    defer s.mu.Unlock()

    if task, exists := s.tasks[taskID]; exists {
        return task.Status
    }
    return "not found"
}
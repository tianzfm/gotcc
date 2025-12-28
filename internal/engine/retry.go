package engine

import (
    "time"
)

type RetryPolicy struct {
    MaxAttempts      int
    BaseInterval     time.Duration
    MaxInterval      time.Duration
    Multiplier       float64
    RandomizationFactor float64
}

type RetryExecutor struct {
    policy RetryPolicy
}

func NewRetryExecutor(policy RetryPolicy) *RetryExecutor {
    return &RetryExecutor{policy: policy}
}

func (re *RetryExecutor) Execute(task func() error) error {
    var err error
    for attempt := 1; attempt <= re.policy.MaxAttempts; attempt++ {
        err = task()
        if err == nil {
            return nil
        }

        // Calculate the wait time based on the retry policy
        waitTime := re.calculateWaitTime(attempt)
        time.Sleep(waitTime)
    }
    return err
}

func (re *RetryExecutor) calculateWaitTime(attempt int) time.Duration {
    waitTime := re.policy.BaseInterval * time.Duration(re.policy.Multiplier*float64(attempt-1))
    if waitTime > re.policy.MaxInterval {
        waitTime = re.policy.MaxInterval
    }
    return waitTime
}
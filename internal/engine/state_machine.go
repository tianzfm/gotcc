package engine

import (
    "errors"
    "sync"
)

type State string

const (
    StatePending      State = "pending"
    StateRunning      State = "running"
    StateSuccess      State = "success"
    StateFailed       State = "failed"
    StateCancelled    State = "cancelled"
)

type StateMachine struct {
    mu     sync.Mutex
    state  State
    events chan State
}

func NewStateMachine() *StateMachine {
    sm := &StateMachine{
        state:  StatePending,
        events: make(chan State),
    }
    go sm.run()
    return sm
}

func (sm *StateMachine) run() {
    for event := range sm.events {
        sm.mu.Lock()
        switch event {
        case StateRunning:
            if sm.state != StatePending {
                sm.mu.Unlock()
                continue
            }
            sm.state = StateRunning
        case StateSuccess:
            if sm.state != StateRunning {
                sm.mu.Unlock()
                continue
            }
            sm.state = StateSuccess
        case StateFailed:
            if sm.state != StateRunning {
                sm.mu.Unlock()
                continue
            }
            sm.state = StateFailed
        case StateCancelled:
            if sm.state != StateRunning {
                sm.mu.Unlock()
                continue
            }
            sm.state = StateCancelled
        }
        sm.mu.Unlock()
    }
}

func (sm *StateMachine) Trigger(event State) error {
    sm.mu.Lock()
    defer sm.mu.Unlock()

    if sm.state == StateSuccess || sm.state == StateFailed || sm.state == StateCancelled {
        return errors.New("cannot trigger event on completed state")
    }

    sm.events <- event
    return nil
}

func (sm *StateMachine) CurrentState() State {
    sm.mu.Lock()
    defer sm.mu.Unlock()
    return sm.state
}

func (sm *StateMachine) Close() {
    close(sm.events)
}
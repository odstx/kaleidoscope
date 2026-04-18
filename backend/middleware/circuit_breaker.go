package middleware

import (
	"sync"
	"time"

	"kaleidoscope/metrics"
)

type CircuitState int

const (
	CircuitClosed CircuitState = iota
	CircuitOpen
	CircuitHalfOpen
)

type CircuitBreaker struct {
	mu            sync.RWMutex
	state         CircuitState
	failures      int
	threshold     int
	timeout       time.Duration
	lastFailure   time.Time
	totalRequests int64
	totalFailures int64
}

func NewCircuitBreaker(threshold int, timeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		state:     CircuitClosed,
		threshold: threshold,
		timeout:   timeout,
	}
}

func (cb *CircuitBreaker) Allow() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {
	case CircuitOpen:
		if time.Since(cb.lastFailure) > cb.timeout {
			cb.state = CircuitHalfOpen
			cb.totalRequests++
			return true
		}
		return false
	case CircuitHalfOpen:
		cb.totalRequests++
		return true
	default:
		cb.totalRequests++
		return true
	}
}

func (cb *CircuitBreaker) RecordSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.failures = 0
	if cb.state == CircuitHalfOpen {
		cb.state = CircuitClosed
	}
}

func (cb *CircuitBreaker) RecordFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.failures++
	cb.lastFailure = time.Now()
	cb.totalFailures++

	if cb.state == CircuitHalfOpen {
		cb.state = CircuitOpen
	} else if cb.failures >= cb.threshold {
		cb.state = CircuitOpen
	}
}

func (cb *CircuitBreaker) GetState() CircuitState {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state
}

type AppCircuitBreakers struct {
	mu         sync.RWMutex
	breakers   map[string]*CircuitBreaker
	threshold  int
	timeout    time.Duration
}

func NewAppCircuitBreakers(threshold int, timeout time.Duration) *AppCircuitBreakers {
	return &AppCircuitBreakers{
		breakers: make(map[string]*CircuitBreaker),
		threshold: threshold,
		timeout:   timeout,
	}
}

func (acb *AppCircuitBreakers) GetBreaker(appName string) *CircuitBreaker {
	acb.mu.Lock()
	defer acb.mu.Unlock()

	if breaker, ok := acb.breakers[appName]; ok {
		return breaker
	}

	breaker := NewCircuitBreaker(acb.threshold, acb.timeout)
	acb.breakers[appName] = breaker
	return breaker
}

func (acb *AppCircuitBreakers) RecordMetrics() {
	acb.mu.RLock()
	defer acb.mu.RUnlock()

	for appName, breaker := range acb.breakers {
		state := breaker.GetState()
		stateValue := 0.0
		switch state {
		case CircuitClosed:
			stateValue = 0
		case CircuitOpen:
			stateValue = 1
		case CircuitHalfOpen:
			stateValue = 2
		}
		metrics.RecordCircuitBreakerState(appName, stateValue)
	}
}

var (
	appCircuitBreakers = NewAppCircuitBreakers(5, 30*time.Second)
)

func InitCircuitBreakerMetrics() {
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		for range ticker.C {
			appCircuitBreakers.RecordMetrics()
		}
	}()
}
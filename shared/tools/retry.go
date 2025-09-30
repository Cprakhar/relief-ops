package tools

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"
)

// RetryableError wraps an error to indicate if it should be retried
type RetryableError struct {
	Err       error
	Retryable bool
	Permanent bool
}

func (r *RetryableError) Error() string {
	return r.Err.Error()
}

// NewRetryableError creates a new retryable error
func NewRetryableError(err error, retryable bool) *RetryableError {
	return &RetryableError{
		Err:       err,
		Retryable: retryable,
		Permanent: false,
	}
}

// NewPermanentError creates a new permanent error (not retryable)
func NewPermanentError(err error) *RetryableError {
	return &RetryableError{
		Err:       err,
		Retryable: false,
		Permanent: true,
	}
}

// IsRetryable returns true if the error should be retried
func IsRetryable(err error) bool {
	if retryErr, ok := err.(*RetryableError); ok {
		return retryErr.Retryable && !retryErr.Permanent
	}
	return true // Default to retryable for unknown errors
}

// RetryConfig holds configuration for retry operations
type RetryConfig struct {
	MaxAttempts   int
	InitialDelay  time.Duration
	MaxDelay      time.Duration
	BackoffFactor float64
	Jitter        bool
}

// DefaultRetryConfig returns a sensible default configuration
func DefaultRetryConfig() *RetryConfig {
	return &RetryConfig{
		MaxAttempts:   3,
		InitialDelay:  100 * time.Millisecond,
		MaxDelay:      30 * time.Second,
		BackoffFactor: 2.0,
		Jitter:        true,
	}
}

// addJitter adds random jitter to prevent thundering herd
func addJitter(delay time.Duration) time.Duration {
	// Add ±25% jitter
	max := int64(delay / 2)
	var n int64
	if max > 0 {
		b := make([]byte, 8)
		_, err := rand.Read(b)
		if err == nil {
			n = int64(b[0]) % max
		}
	}
	jitter := time.Duration(n)
	return delay + jitter - delay/4
}

// RetryWithBackoff performs retries with exponential backoff and context support
func RetryWithBackoff(ctx context.Context, config *RetryConfig, fn func() error) error {
	if config == nil {
		config = DefaultRetryConfig()
	}

	var lastErr error
	delay := config.InitialDelay

	for attempt := 1; attempt <= config.MaxAttempts; attempt++ {
		select {
		case <-ctx.Done():
			return fmt.Errorf("retry cancelled: %w", ctx.Err())
		default:
		}

		lastErr = fn()
		if lastErr == nil {
			return nil
		}

		// Don't retry if error is not retryable
		if !IsRetryable(lastErr) {
			return fmt.Errorf("non-retryable error after %d attempts: %w", attempt, lastErr)
		}

		// Don't sleep after the last attempt
		if attempt == config.MaxAttempts {
			break
		}

		// Calculate next delay with jitter if enabled
		nextDelay := delay
		if config.Jitter {
			nextDelay = addJitter(delay)
		}

		// Cap the delay at maxDelay
		if nextDelay > config.MaxDelay {
			nextDelay = config.MaxDelay
		}

		select {
		case <-ctx.Done():
			return fmt.Errorf("retry cancelled during backoff: %w", ctx.Err())
		case <-time.After(nextDelay):
		}

		// Calculate next delay for next iteration
		delay = time.Duration(float64(delay) * config.BackoffFactor)
	}

	return fmt.Errorf("max retry attempts (%d) exceeded: %w", config.MaxAttempts, lastErr)
}

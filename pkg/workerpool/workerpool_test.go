package workerpool

import (
	"errors"
	"testing"
	"time"
)

func TestWorkerPool(t *testing.T) {
	pool := NewWorkerPool(5)

	// Test task execution
	taskExecuted := false
	pool.Tasks <- Task{
		Execute: func() error {
			taskExecuted = true
			return nil
		},
	}

	time.Sleep(100 * time.Millisecond) // Give the worker some time
	if !taskExecuted {
		t.Error("Task was not executed")
	}

	// Test error handling
	expectedError := errors.New("an error occurred")
	pool.Tasks <- Task{
		Execute: func() error {
			return expectedError
		},
	}

	select {
	case task := <-pool.Errors:
		if task.Error != expectedError {
			t.Errorf("Expected error %v, but got %v", expectedError, task.Error)
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("Expected an error, but did not get one")
	}

	// Test shutdown
	pool.Stop()
	select {
	case <-pool.Shutdown:
		// Expected
	case <-time.After(100 * time.Millisecond):
		t.Error("Workers did not shut down in time")
	}
}

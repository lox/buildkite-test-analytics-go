package example

import (
	"testing"
	"time"
)

func TestSleep10(t *testing.T) {
	time.Sleep(time.Millisecond * 10)
}

func TestSleep100(t *testing.T) {
	time.Sleep(time.Millisecond * 100)
}

func TestLogging(t *testing.T) {
	t.Log("Testing logging")
}

func TestFailure(t *testing.T) {
	t.Error("An error occurred")
}

func TestSubtest(t *testing.T) {
	t.Run("sub1", func(tt *testing.T) {
		time.Sleep(time.Millisecond * 10)
	})
	t.Run("sublog", func(tt *testing.T) {
		tt.Log("Test sublog")
	})
}

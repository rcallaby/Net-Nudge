package tests

import (
    "testing"
    "net-nudge/internal/jobs"
    "net-nudge/internal/registry"
    "net-nudge/mocks"
)

func TestEndToEnd(t *testing.T) {
    reg := registry.NewRegistry()
    reg.Register(&mocks.MockTool{Name: "mock", Success: true})

    job := jobs.Job{Target: "endtoend.com", Tool: &mocks.MockTool{Name: "mock", Success: true}}
    result, err := job.Run()
    if err != nil {
        t.Fatalf("E2E test failed: %v", err)
    }
    if result == "" {
        t.Errorf("expected non-empty result")
    }
}
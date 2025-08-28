package jobs

import (
    "testing"
    "net-nudge/mocks"
)

func TestJobRun(t *testing.T) {
    tool := &mocks.MockTool{Name: "mock", Success: true}
    job := Job{Target: "test.com", Tool: tool}

    result, err := job.Run()
    if err != nil {
        t.Fatalf("job failed: %v", err)
    }
    if result == "" {
        t.Errorf("expected result, got empty string")
    }
}
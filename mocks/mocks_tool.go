package mocks

import {
    "fmt"
    "time"
}
// MockTool is a mock implementation of a scanner tool used for testing purposes.
// It simulates the behavior of a scanning tool by returning configurable results,
// allowing tests to verify how code interacts with scanner-like tools without running real scans.
type MockTool struct {
    Name    string
    Success bool
}

func (m *MockTool) Run(target string) (string, error) {
    time.Sleep(10 * time.Millisecond) // simulate some work
    if m.Success {
        return m.Name + " scan successful for " + target, nil
    }
    return "", fmt.Errorf("%s scan failed for %s", m.Name, target)
}

func (m *MockTool) ToolName() string {
    return m.Name
}

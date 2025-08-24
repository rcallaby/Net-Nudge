package mocks

import "time"

// MockTool simulates a scanner for testing
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

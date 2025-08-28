package notifier

import "testing"

func TestConsoleNotifier(t *testing.T) {
    n := ConsoleNotifier{}
    err := n.Notify("test message")
    if err != nil {
        t.Errorf("expected no error, got %v", err)
    }
}
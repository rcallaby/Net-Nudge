package config

import (
    "os"
    "testing"
)

func TestLoadConfig(t *testing.T) {
    yamlData := `
scheduler:
  interval: "1h"
targets:
  - domain: "example.com"
scans:
  - tool: "nmap"
    args: ["-sV"]
`
    tmp := "test_config.yaml"
    os.WriteFile(tmp, []byte(yamlData), 0644)
    defer os.Remove(tmp)

    cfg, err := LoadConfig(tmp)
    if err != nil {
        t.Fatalf("failed to load config: %v", err)
    }

    if cfg.Scheduler.Interval != "1h" {
        t.Errorf("expected interval 1h, got %s", cfg.Scheduler.Interval)
    }
    if cfg.Targets[0].Domain != "example.com" {
        t.Errorf("expected domain example.com, got %s", cfg.Targets[0].Domain)
    }
}
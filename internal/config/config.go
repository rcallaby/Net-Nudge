package config

import (
    "os"
    "time"

    yaml "gopkg.in/yaml.v3"
)

type Config struct {
    General struct {
        DataDir   string   `yaml:"data_dir"`
        TempDir   string   `yaml:"temp_dir"`
        Workers   int      `yaml:"workers"`
        EnvPath   []string `yaml:"env_path"`
        Timeout   string   `yaml:"global_timeout"`
    } `yaml:"general"`

    Notifiers struct {
        Email  *EmailConfig  `yaml:"email"`
        Slack  *SlackConfig  `yaml:"slack"`
        Webhook *WebhookConfig `yaml:"webhook"`
    } `yaml:"notifiers"`

    Tools ToolsConfig `yaml:"tools"`

    Pipelines []Pipeline `yaml:"pipelines"`
}

type ToolsConfig struct {
    Subfinder *BinaryTool `yaml:"subfinder"`
    Amass     *BinaryTool `yaml:"amass"`
    HTTPX     *BinaryTool `yaml:"httpx"`
    Naabu     *BinaryTool `yaml:"naabu"`
    Nmap      *BinaryTool `yaml:"nmap"`
    Masscan   *BinaryTool `yaml:"masscan"`
    Nuclei    *BinaryTool `yaml:"nuclei"`
}

type BinaryTool struct {
    Enabled bool     `yaml:"enabled"`
    Path    string   `yaml:"path"`
    Args    []string `yaml:"args"`
    Rate    int      `yaml:"rate"`
}

type Pipeline struct {
    Name     string   `yaml:"name"`
    Cron     string   `yaml:"cron"` // robfig/cron spec
    Scope    Scope    `yaml:"scope"`
    Stages   []string `yaml:"stages"` // e.g. ["subdomains","probe","ports","vuln"]
    Reports  ReportCfg `yaml:"reports"`
}

type Scope struct {
    Domains       []string `yaml:"domains"`
    SubdomainWordlist string `yaml:"subdomain_wordlist"`
    ExcludeDomains []string `yaml:"exclude_domains"`
    CIDRs         []string `yaml:"cidrs"`
    Ports         []string `yaml:"ports"` // e.g. ["80,443","top-1000"]
}

type ReportCfg struct {
    Dir    string `yaml:"dir"`
    SaveJSON bool `yaml:"json"`
    SaveCSV  bool `yaml:"csv"`
}

// Notifier configs

type EmailConfig struct {
    Enabled  bool   `yaml:"enabled"`
    From     string `yaml:"from"`
    To       []string `yaml:"to"`
    SMTPHost string `yaml:"smtp_host"`
    SMTPPort int    `yaml:"smtp_port"`
    Username string `yaml:"username"`
    Password string `yaml:"password"`
    StartTLS bool   `yaml:"starttls"`
}

type SlackConfig struct {
    Enabled bool   `yaml:"enabled"`
    Webhook string `yaml:"webhook"`
    Channel string `yaml:"channel"`
}

type WebhookConfig struct {
    Enabled bool   `yaml:"enabled"`
    URL     string `yaml:"url"`
    Secret  string `yaml:"secret"`
}

func Load(path string) (*Config, error) {
    b, err := os.ReadFile(path)
    if err != nil { return nil, err }
    var cfg Config
    if err := yaml.Unmarshal(b, &cfg); err != nil { return nil, err }
    if cfg.General.Workers <= 0 { cfg.General.Workers = 6 }
    if cfg.General.DataDir == "" { cfg.General.DataDir = "data" }
    if cfg.General.TempDir == "" { cfg.General.TempDir = "tmp" }
    if cfg.General.Timeout == "" { cfg.General.Timeout = "2h" }
    return &cfg, nil
}

func (c *Config) GlobalTimeout() time.Duration {
    d, _ := time.ParseDuration(c.General.Timeout)
    if d == 0 { d = 2 * time.Hour }
    return d
}
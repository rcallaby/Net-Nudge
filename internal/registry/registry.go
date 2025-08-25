package registry

import "sync"

// Tool interface allows adding new tools easily
type Tool interface {
	Name() string
	Binary() string
	BaseArgs() []string
}

var (
	toolRegistry = make(map[string]Tool, 8) // Preallocate space for 8 tools; adjust as needed
	mu           sync.RWMutex
)

// RegisterTool adds a tool to the registry
func RegisterTool(t Tool) {
	mu.Lock()
	defer mu.Unlock()
	toolRegistry[t.Name()] = t
}

// GetTool fetches a tool by name
func GetTool(name string) (Tool, bool) {
	mu.RLock()
	defer mu.RUnlock()
	t, ok := toolRegistry[name]
	return t, ok
}

// ListTools returns all registered tools
func ListTools() []Tool {
	mu.RLock()
	defer mu.RUnlock()
	tools := make([]Tool, 0, len(toolRegistry))
	for _, t := range toolRegistry {
		tools = append(tools, t)
	}
	return tools
}

// Example tool implementations
type Nmap struct{}

func (Nmap) Name() string     { return "nmap" }
func (Nmap) Binary() string   { return "nmap" }
func (Nmap) BaseArgs() []string { return []string{"-sV"} }

type Gobuster struct{}

func (Gobuster) Name() string     { return "gobuster" }
func (Gobuster) Binary() string   { return "gobuster" }
func (Gobuster) BaseArgs() []string { return []string{"dir", "-u"} }

// InitDefaultTools registers the default set of tools in the registry.
func InitDefaultTools() {
	RegisterTool(Nmap{})
	RegisterTool(Gobuster{})
}


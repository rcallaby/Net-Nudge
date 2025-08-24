package tools


import (
"context"
"encoding/json"
"os/exec"
)


type nuclei struct { bin string; args []string }

func NewNuclei(path string, args []string) VulnScanner { return &nuclei{bin: path, args: args} }

func (n *nuclei) Scan(ctx context.Context, targets []string) ([]map[string]any, error) {
args := append([]string{"-json","-l","-"}, n.args...)
cmd := exec.CommandContext(ctx, n.bin, args...)
cmd.Stdin = sliceReader(targets)
lines, err := runAndCollect(ctx, cmd)
if err != nil { return nil, err }
var f []map[string]any
for _, ln := range lines {
var m map[string]any
if json.Unmarshal([]byte(ln), &m) == nil { f = append(f, m) }
}
return f, nil
}
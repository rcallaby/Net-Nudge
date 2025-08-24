package tools

import (
	"context"
	"encoding/json"
	"os/exec"
	"strings"
	"bufio"
	"bytes"
)



// PortScanner defines the interface for scanning ports on given hosts.
type PortScanner interface {
	Scan(ctx context.Context, hosts []string, ports []string) (map[string][]string, error)
}

// naabu implements the PortScanner interface using the naabu binary.
type naabu struct {
	bin  string
	args []string
}
func NewNaabu(path string, args []string) PortScanner { return &naabu{bin: path, args: args} }
func (n *naabu) Scan(ctx context.Context, hosts []string, ports []string) (map[string][]string, error) {
	args := append([]string{"-json", "-host", "-"}, n.args...)
	if len(ports) > 0 {
		args = append(args, "-p", strings.Join(ports, ","))
	}
	cmd := exec.CommandContext(ctx, n.bin, args...)
	cmd.Stdin = sliceReader(hosts)
	lines, err := runAndCollect(ctx, cmd)
	if err != nil {
		return nil, err
	}
	res := make(map[string][]string)
	for _, ln := range lines {
		var m map[string]any
		if json.Unmarshal([]byte(ln), &m) == nil {
			host, ok := m["host"].(string)
			port := fmtAny(m["port"])
			if ok && host != "" && port != "" {
				res[host] = append(res[host], port)
			}
		}
	}
	return res, nil
}

// runAndCollect runs the given exec.Cmd and returns its stdout output as a slice of lines.
func runAndCollect(ctx context.Context, cmd *exec.Cmd) ([]string, error) {
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	var lines []string
	scanner := bufio.NewScanner(&out)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
// fmtAny returns the string representation of the provided value v.
// It handles common types such as string, float64, and integer types.
// For float64 values that are whole numbers, it omits the decimal part.
// If the type is not recognized, it returns an empty string.
func fmtAny(v interface{}) string {
	switch t := v.(type) {
	case string:
		return t
	case float64:
		// Remove decimal if it's an integer value
		if t == float64(int64(t)) {
			return fmt.Sprintf("%d", int64(t))
		}
		return fmt.Sprintf("%v", t)
	case int, int32, int64:
		return fmt.Sprintf("%v", t)
	default:
		return ""
	}
}}
}
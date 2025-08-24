package tools

import (
	"context"
	"os/exec"
	"strings"
	"fmt"
	"regexp"
)

type masscan struct {
	bin  string
	args []string
	bin  string
	args []string
}

func NewMasscan(path string, args []string) PortScanner {
	return &masscan{bin: path, args: args}
}

func (m *masscan) Scan(ctx context.Context, hosts []string, ports []string) (map[string][]string, error) {
	args := []string{"-oL", "-"}
	if len(ports) > 0 {
		portArg := strings.Join(ports, ",")
		args = append(args, "-p", portArg)
	} else {
		args = append(args, "-p", "1-65535")
	}
	if len(hosts) > 0 {
		args = append(args, "--range", strings.Join(hosts, ","))
	}
	args = append(args, m.args...)
	cmd := exec.CommandContext(ctx, m.bin, args...)
	cmd.Stdin = sliceReader(hosts)
	lines, err := runAndCollect(ctx, cmd)
	if err != nil {
		return nil, err
	}
	out := map[string][]string{}
	for _, ln := range lines {
		host, port := parseMasscanLine(ln)
		if host != "" && port != "" {
			out[host] = append(out[host], port)
		}
	}
	return out, nil
}

// parseMasscanLine parses a line from masscan -oL output and returns host and port if found.
func parseMasscanLine(line string) (string, string) {
	// Example masscan -oL line: "open tcp 80 192.168.1.1 12345"
	re := regexp.MustCompile(`^open\s+\w+\s+(\d+)\s+([\d\.]+)`)
	matches := re.FindStringSubmatch(line)
	if len(matches) == 3 {
		return matches[2], matches[1]
	}
	return "", ""
}
return out, nil

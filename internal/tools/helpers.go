package tools


import (
"bytes"
"strings"
)


type sliceReader []string


func (s sliceReader) Read(p []byte) (int, error) {
b := []byte(strings.Join([]string(s), "\n"))
r := bytes.NewReader(b)
return r.Read(p)
}


func parseNmapGrepable(line string) (string, []string) {
// very light parser
// Format: Host: 1.2.3.4 () Ports: 80/open/tcp//http///, 443/open/tcp//ssl///
if !strings.Contains(line, "Host:") || !strings.Contains(line, "Ports:") { return "", nil }
parts := strings.Split(line, "Ports:")
hostPart := strings.TrimSpace(strings.TrimPrefix(strings.Fields(parts[0])[1], ""))
var ports []string
for _, seg := range strings.Split(parts[1], ",") {
fields := strings.Split(strings.TrimSpace(seg), "/")
if len(fields) > 0 { ports = append(ports, fields[0]) }
}
return hostPart, ports
}


func parseMasscanLine(line string) (string, string) {
// example: Host: 1.2.3.4 () Ports: 80/open
if !strings.Contains(line, "Ports:") { return "", "" }
fields := strings.Fields(line)
if len(fields) < 6 { return "", "" }
host := fields[1]
port := strings.Split(fields[5], "/")[0]
return host, port
}


func fmtAny(v any) string { return strings.TrimSpace(strings.Trim(strings.ReplaceAll(strings.TrimSpace(fmt.Sprintf("%v", v)), "\"", ""), "[]")) }
package tools


import (
"context"
"os/exec"
)


type nmap struct { bin string; args []string }
func NewNmap(path string, args []string) PortScanner { return &nmap{bin: path, args: args} }
func (n *nmap) Scan(ctx context.Context, hosts []string, ports []string) (map[string][]string, error) {
// Simplified: nmap -Pn -p <ports> -iL - -oG - | parse open ports
args := append([]string{"-Pn","-oG","-","-iL","-"}, n.args...)
if len(ports) > 0 { args = append(args, "-p", ports[0]) }
cmd := exec.CommandContext(ctx, n.bin, args...)
cmd.Stdin = sliceReader(hosts)
lines, err := runAndCollect(ctx, cmd)
if err != nil { return nil, err }
out := map[string][]string{}
for _, ln := range lines {
// grepable: Host: 1.2.3.4 () Ports: 80/open/tcp//http///
host, ports := parseNmapGrepable(ln)
if host != "" && len(ports) > 0 { out[host] = append(out[host], ports...) }
}
return out, nil
}
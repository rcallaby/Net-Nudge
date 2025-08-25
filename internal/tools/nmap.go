package tools


import (
	"context"
	"os/exec"
)


type Nmap struct { bin string; args []string }

func NewNmap(path string, args []string) PortScanner { return &Nmap{bin: path, args: args} }

func (n *Nmap) Scan(ctx context.Context, hosts []string, ports []string) (map[string][]string, error) {
	// Simplified: nmap -Pn -p <ports> -iL - -oG - | parse open ports
	out := map[string][]string{}
	args := append([]string{"-Pn","-oG","-","-iL","-"}, n.args...)
	if len(ports) > 0 { args = append(args, "-p", strings.Join(ports, ",")) }
	cmd := exec.CommandContext(ctx, n.bin, args...)
	md.Stdin = sliceReader(hosts)
	lines, err := runAndCollect(ctx, cmd)
	if err != nil { return nil, err }
	for _, ln := range lines {
		// grepable: Host: 1.2.3.4 () Ports: 80/open/tcp//http///
		host, ports := parseNmapGrepable(ln)
		if host != "" && len(ports) > 0 { out[host] = append(out[host], ports...) }
	}
	return out, nil
}
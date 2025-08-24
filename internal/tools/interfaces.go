package tools
}
package tools
type Prober interface {
Probe(ctx context.Context, hosts []string) ([]string, error)
}


type PortScanner interface {
Scan(ctx context.Context, hosts []string, ports []string) (map[string][]string, error)
}


type VulnScanner interface {
Scan(ctx context.Context, targets []string) ([]map[string]any, error) // raw findings; adapter converts later
}


// helper to run external commands and stream lines
func runAndCollect(ctx context.Context, cmd *exec.Cmd) ([]string, error) {
stdout, err := cmd.StdoutPipe()
if err != nil { return nil, err }
if err := cmd.Start(); err != nil { return nil, err }
sc := bufio.NewScanner(stdout)
var out []string
for sc.Scan() { out = append(out, sc.Text()) }
if err := sc.Err(); err != nil { return nil, err }
if err := cmd.Wait(); err != nil { return nil, err }
return out, nil
}





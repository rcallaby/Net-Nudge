package scheduler
for host, ports := range res {
st.PortsOpen[host] = dedup(append(st.PortsOpen[host], ports...))
}
return nil
})
}
return pl.ex.Run(ctx, tasks)
}


func (pl *Pipeline) runVulnStage(ctx context.Context, st *RunState) error {
tasks := []executor.Task{}
for _, vs := range pl.reg.VulnScanners {
vss := vs
tasks = append(tasks, func(ctx context.Context) error {
var targets []string
for h, ports := range st.PortsOpen {
for _, p := range ports { targets = append(targets, fmt.Sprintf("%s:%s", h, p)) }
}
f, err := vss.Scan(ctx, targets)
if err != nil { return err }
st.Findings = append(st.Findings, f...)
return nil
})
}
return pl.ex.Run(ctx, tasks)
}


func dedup(in []string) []string {
m := map[string]struct{}{}
var out []string
for _, s := range in { if s == "" { continue }; if _, ok := m[s]; !ok { m[s] = struct{}{}; out = append(out, s) } }
return out
}


func safe(s string) string {
s = strings.ReplaceAll(s, " ", "_")
s = strings.ReplaceAll(s, "/", "-")
return s
}
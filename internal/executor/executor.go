package executor


tasks := []executor.Task{}


for _, p := range pl.reg.Probers {
	pr := p
	tasks = append(tasks, func(ctx context.Context) error {
	res, err := pr.Probe(ctx, st.Subdomains)
	if err != nil { return err }
		alive = append(alive, res...)
		return nil
	})
}
if err := pl.ex.Run(ctx, tasks); err != nil { return nil, err }
	path := filepath.Join(st.WorkDir, "alive.txt")
	if err := os.WriteFile(path, []byte(strings.Join(dedup(alive), "\n")), 0o644); err != nil { return nil, err }
	return alive, nil



func (pl *Pipeline) runPortsStage(ctx context.Context, st *RunState) error {
	tasks := []executor.Task{}
	for _, ps := range pl.reg.PortScanners {
		pss := ps
		tasks = append(tasks, func(ctx context.Context) error {
		res, err := pss.Scan(ctx, st.AliveHosts, pl.p.Scope.Ports)
	if err != nil { return err }
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
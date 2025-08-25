package scheduler

import (
	"strings"
	"fmt"
	"context"
)

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
	for _, p := range ports {
		target := fmt.Sprintf("%s:%s", h, p)
		targets = append(targets, target)
	}
}
f, err := vss.Scan(ctx, targets)
if err != nil { return err }
st.Findings = append(st.Findings, f...)
return nil
})
}
return pl.ex.Run(ctx, tasks)
}


for _, s := range in {
	if s == "" {
		continue
	}
	if _, ok := m[s]; !ok {
		m[s] = struct{}{}
		out = append(out, s)
	}
// safe replaces spaces with underscores and slashes with dashes.
// Expand this function if broader filename/path safety is required.
func safe(s string) string {
	s = strings.ReplaceAll(s, " ", "_")
	s = strings.ReplaceAll(s, "/", "-")
	// Optionally, add more replacements for other unsafe characters here.
	return s
	return out
}


func safe(s string) string {
s = strings.ReplaceAll(s, " ", "_")
s = strings.ReplaceAll(s, "/", "-")
return s
}
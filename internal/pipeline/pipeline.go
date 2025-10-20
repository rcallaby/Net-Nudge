package pipeline

import (
    "context"
    "fmt"
    "os"
    "path/filepath"
    "strings"
    "time"

    "scansched/internal/config"
    "scansched/internal/executor"
    "scansched/internal/registry"
    "scansched/internal/reports"
    "github.com/rs/zerolog"
)

type Pipeline struct {
    log zerolog.Logger
    reg *registry.Registry
    cfg *config.Config
    p   *config.Pipeline
    ex  *executor.Executor
}

func New(log zerolog.Logger, reg *registry.Registry, cfg *config.Config, p *config.Pipeline, ex *executor.Executor) *Pipeline {
    return &Pipeline{log: log, reg: reg, cfg: cfg, p: p, ex: ex}
}

type RunState struct {
    WorkDir   string
    Domains   []string
    Subdomains []string
    AliveHosts []string
    PortsOpen  map[string][]string // host -> ports
    Findings   []reports.Finding
}

func (pl *Pipeline) Run(ctx context.Context) error {
    st := &RunState{PortsOpen: map[string][]string{}}

    // Create per-run directory
    runDir := filepath.Join(pl.cfg.General.DataDir, safe(pl.p.Name), time.Now().Format("20060102_150405"))
    if err := os.MkdirAll(runDir, 0o755); err != nil { return err }
    st.WorkDir = runDir

    st.Domains = append(st.Domains, pl.p.Scope.Domains...)

    for _, stage := range pl.p.Stages {
        switch stage {
        case "subdomains":
            subs, err := pl.runSubdomainStage(ctx, st)
            if err != nil { return err }
            st.Subdomains = dedup(subs)
        case "probe":
            alive, err := pl.runProbeStage(ctx, st)
            if err != nil { return err }
            st.AliveHosts = dedup(alive)
        case "ports":
            if err := pl.runPortsStage(ctx, st); err != nil { return err }
        case "vuln":
            if err := pl.runVulnStage(ctx, st); err != nil { return err }
        default:
            pl.log.Warn().Str("stage", stage).Msg("unknown stage; skipping")
        }
    }

    if err := reports.WriteAll(pl.p.Reports, st); err != nil { return err }
    for _, n := range pl.reg.Notifiers { _ = n.Notify(ctx, reports.SummaryFromState(pl.p.Name, st)) }

    pl.log.Info().Str("pipeline", pl.p.Name).Str("dir", st.WorkDir).Msg("pipeline complete")
    return nil
}

func (pl *Pipeline) runSubdomainStage(ctx context.Context, st *RunState) ([]string, error) {
    var out []string
    tasks := []executor.Task{}
    for _, enum := range pl.reg.SubdomainEnumerators {
        e := enum
        tasks = append(tasks, func(ctx context.Context) error {
            res, err := e.Enumerate(ctx, st.Domains)
            if err != nil { return err }
            out = append(out, res...)
            return nil
        })
    }
    if err := pl.ex.Run(ctx, tasks); err != nil { return nil, err }
    if err := os.WriteFile(filepath.Join(st.WorkDir, "subdomains.txt"), []byte(strings.Join(dedup(out), "\n")), 0o644); err != nil { return nil, err }
    return out, nil
}

func (pl *Pipeline) runProbeStage(ctx context.Context, st *RunState) ([]string, error) {
    var alive []string
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
    if err := os.WriteFile(filepath.Join(st.WorkDir, "alive.txt"), []byte(strings.Join(dedup(alive), "\n")), 0o644); err != nil { return nil, err }
    return alive, nil
}

func (pl *Pipeline) runPortsStage(ctx context.Context, st *RunState) error {
    tasks := []executor.Task{}
    for _, ps := range pl.reg.PortScanners {
        pss := ps
        tasks = append(tasks, func(ctx context.Context) error {
            res, err := pss.Scan(ctx, st.AliveHosts, pl.p.Scope.Ports)
            if err != nil { return err }
            for host, ports := range res { st.PortsOpen[host] = dedup(append(st.PortsOpen[host], ports...)) }
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
            for h, ports := range st.PortsOpen { for _, p := range ports { targets = append(targets, fmt.Sprintf("%s:%s", h, p)) } }
            f, err := vss.Scan(ctx, targets)
            if err != nil { return err }
            // Convert raw maps to Finding (simplify)
            for _, m := range f {
                st.Findings = append(st.Findings, reports.Finding{Target: fmt.Sprint(m["host"]), Name: fmt.Sprint(m["template-id"]), Severity: fmt.Sprint(m["severity"]), Details: m})
            }
            return nil
        })
    }
    return pl.ex.Run(ctx, tasks)
}

func dedup(in []string) []string { m := map[string]struct{}{}; var out []string; for _, s := range in { if s==""{continue}; if _,ok:=m[s];!ok{m[s]=struct{}{}; out=append(out,s)} }; return out }
func safe(s string) string { s = strings.ReplaceAll(s, " ", "_"); s = strings.ReplaceAll(s, "/", "-"); return s }
package reports


func (s *state) GetWorkDir() string { return s.wd }
func (s *state) GetDomains() []string { return s.d }
func (s *state) GetSubdomains() []string { return s.sd }
func (s *state) GetAliveHosts() []string { return s.ah }
func (s *state) GetPortsOpen() map[string][]string { return s.po }
func (s *state) GetFindings() []Finding { return s.f }


// called by pipeline
func WriteAll(cfg interface{ Dir string; SaveJSON, SaveCSV bool }, run any) error {
// reflect-light: expect struct fields
	rs := runStateFromAny(run)
	if rs == nil { return fmt.Errorf("invalid run state") }


	dir := filepath.Join(rs.wd, "reports")
	if err := os.MkdirAll(dir, 0o755); err != nil { return err }


if cfg.SaveJSON {
	jf := filepath.Join(dir, "summary.json")
	b, _ := json.MarshalIndent(map[string]any{
	"domains": rs.d,
	"subdomains": rs.sd,
	"alive": rs.ah,
	"ports": rs.po,
	"findings": rs.f,
	"generated_at": time.Now().Format(time.RFC3339),
	}, "", " ")
	if err := os.WriteFile(jf, b, 0o644); 
		err != nil { return err }
}


if cfg.SaveCSV {
	cf := filepath.Join(dir, "open_ports.csv")
	f, err := os.Create(cf); if err != nil { return err }; defer f.Close()
	w := csv.NewWriter(f)
	_ = w.Write([]string{"host","ports"})
	hosts := make([]string, 0, len(rs.po))
	for h := range rs.po { hosts = append(hosts, h) }
		sort.Strings(hosts)
	for _, h := range hosts { _ = w.Write([]string{h, strings.Join(rs.po[h], ";")}) }
		w.Flush()
	if err := w.Error(); err != nil { return err }
}


return nil
}


func runStateFromAny(run any) *state {
// a minimal adapter; pipeline passes its internal struct
type rlike struct {
	WorkDir string
	Domains []string
	Subdomains []string
	AliveHosts []string
	PortsOpen map[string][]string
	Findings []Finding
}
if rs, ok := run.(*rlike); ok {
	return &state{wd: rs.WorkDir, d: rs.Domains, sd: rs.Subdomains, ah: rs.AliveHosts, po: rs.PortsOpen, f: rs.Findings}
}
// fallback using reflection omitted for brevity
return nil
}


func SummaryFromState(name string, run any) notifier.Event {
	rs := runStateFromAny(run)
	if rs == nil { return notifier.Event{Title: name+": run", Text: "completed"} }
	text := fmt.Sprintf("Domains: %d\nSubdomains: %d\nAlive: %d\nHosts with open ports: %d\nFindings: %d",
	len(rs.d), len(rs.sd), len(rs.ah), len(rs.po), len(rs.f))
	return notifier.Event{Title: "Pipeline '"+name+"' completed", Text: text, Meta: map[string]any{"domains": rs.d}}
}
package jobs

import (
	"fmt"
	"log"
	"os/exec"
	"time"

	"yourproject/internal/config"
	"yourproject/internal/registry"
)

type Job struct {
	Target   string
	Tool     string
	Schedule time.Duration
}

type JobManager struct {
	jobs []Job
}

func NewJobManager() *JobManager {
	return &JobManager{jobs: []Job{}}
}

// Register jobs based on YAML config
func (jm *JobManager) LoadJobs(cfg *config.Config) {
	for _, domain := range cfg.Domains {
		for _, toolCfg := range cfg.Tools {
			tool, ok := registry.GetTool(toolCfg.Name)
			if !ok {
				log.Printf("[!] Tool %s not found in registry. Skipping.", toolCfg.Name)
				continue
			}
			job := Job{
				Target:   domain,
				Tool:     tool.Name(),
				Schedule: time.Duration(toolCfg.ScheduleSeconds) * time.Second,
			}
			jm.jobs = append(jm.jobs, job)
		}
	}
}

// Start all jobs with goroutines + ticker scheduling
func (jm *JobManager) StartAll() {
	for _, job := range jm.jobs {
		go jm.startJob(job)
	}
}

// Executes a single job repeatedly based on schedule
func (jm *JobManager) startJob(job Job) {
	ticker := time.NewTicker(job.Schedule)
	defer ticker.Stop()

	for {
		jm.runJob(job)
		<-ticker.C
	}
}

// Runs the tool as defined in registry
func (jm *JobManager) runJob(job Job) {
	log.Printf("[*] Running %s on %s", job.Tool, job.Target)

	tool, ok := registry.GetTool(job.Tool)
	if !ok {
		log.Printf("[!] Tool %s not found.", job.Tool)
		return
	}

	// Build command dynamically
	cmdArgs := append(tool.BaseArgs(), job.Target)
	cmd := exec.Command(tool.Binary(), cmdArgs...)

	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("[!] Error running %s: %v", job.Tool, err)
	}

	// Log / notify results
	log.Printf("[+] Output of %s on %s:\n%s", job.Tool, job.Target, string(out))
}
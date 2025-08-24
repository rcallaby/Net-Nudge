package tools


import (
	"context"
	"os/exec"
)


type subfinder struct { bin string; args []string }

func NewSubfinder(path string, args []string) SubdomainEnumerator { return &subfinder{bin: path, args: args} }

func (s *subfinder) Enumerate(ctx context.Context, domains []string) ([]string, error) {
	args := append(append([]string{"-silent"}, s.args...), "-dL", "-")
	cmd := exec.CommandContext(ctx, s.bin, args...)
	cmd.Stdin = sliceReader(domains)
	return runAndCollect(ctx, cmd)
}

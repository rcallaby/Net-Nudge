package tools


import (
	"context"
	"os/exec"
)


type httpx struct { bin string; args []string }

func NewHTTPX(path string, args []string) Prober { return &httpx{bin: path, args: args} }

func (h *httpx) Probe(ctx context.Context, hosts []string) ([]string, error) {
	args := append(append([]string{"-silent"}, h.args...), "-l", "-")
	cmd := exec.CommandContext(ctx, h.bin, args...)
	cmd.Stdin = sliceReader(hosts)
	return runAndCollect(ctx, cmd)
}
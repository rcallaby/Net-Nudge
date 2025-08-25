package tools


import (
	"context"
	"os/exec"
)


type amass struct { bin string; args []string }

func NewAmass(path string, args []string) SubdomainEnumerator { return &amass{bin: path, args: args} }

func (a *amass) Enumerate(ctx context.Context, domains []string) ([]string, error) {
// amass enum -silent -d domain
	var all []string
	for _, d := range domains {
		args := append([]string{"enum","-silent","-d", d}, a.args...)
		out, err := runAndCollect(ctx, exec.CommandContext(ctx, a.bin, args...))
		if err != nil { return nil, err }
		all = append(all, out...)
	}
	return all, nil
}
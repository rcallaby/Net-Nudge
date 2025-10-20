package executor

import (
    "context"
    "sync"

    "github.com/rs/zerolog"
)

type Task func(ctx context.Context) error

type Executor struct {
    log zerolog.Logger
    workers int
}

func New(log zerolog.Logger, workers int) *Executor {
    if workers <= 0 { workers = 4 }
    return &Executor{log: log, workers: workers}
}

func (e *Executor) Run(ctx context.Context, tasks []Task) error {
    sem := make(chan struct{}, e.workers)
    var wg sync.WaitGroup
    errs := make(chan error, len(tasks))

    for _, t := range tasks {
        wg.Add(1)
        sem <- struct{}{}
        go func(tf Task) {
            defer wg.Done(); defer func() { <-sem }()
            if err := tf(ctx); err != nil { errs <- err }
        }(t)
    }

    wg.Wait(); close(errs)
    var first error
    for err := range errs { if first == nil { first = err } }
    return first
}


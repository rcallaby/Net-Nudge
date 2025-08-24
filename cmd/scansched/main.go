package main

import (
    "context"
    "flag"
    "fmt"
    "os"
    "os/signal"
    "syscall"

    "scansched/internal/config"
    "scansched/internal/logging"
    "scansched/internal/registry"
    "scansched/internal/scheduler"
)

func main() {
    cfgPath := flag.String("config", "configs/sample.yaml", "Path to YAML configuration")
    once := flag.Bool("once", false, "Run scheduled jobs once (no cron loop)")
    validate := flag.Bool("validate", false, "Validate config and exit")
    flag.Parse()

    log := logging.New()

    cfg, err := config.Load(*cfgPath)
    if err != nil { log.Fatal().Err(err).Msg("failed to load config") }

    if *validate { log.Info().Msg("config OK"); return }

    ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
    defer cancel()

    // Registry wires tool adapters and notifiers based on config
    reg := registry.New(log)
    if err := reg.WireFromConfig(ctx, cfg); err != nil {
        log.Fatal().Err(err).Msg("failed to wire registry")
    }

    sched := scheduler.New(log, reg, cfg)

    if *once {
        if err := sched.RunOnce(ctx); err != nil { log.Fatal().Err(err).Msg("run-once failed") }
        fmt.Println("Done")
        return
    }

    if err := sched.Run(ctx); err != nil {
        log.Fatal().Err(err).Msg("scheduler failed")
    }
}


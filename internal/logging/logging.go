package logging

import (
    "os"

    "github.com/rs/zerolog"
    "github.com/rs/zerolog/log"
)

type Logger = zerolog.Logger

func New() Logger {
    l := log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
    return l.Level(zerolog.InfoLevel)
}

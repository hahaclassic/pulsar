package main

import (
	"log/slog"

	"github.com/hahaclassic/pulsar/internal/pulse"
)

func main() {
	err := pulse.Start()
	if err != nil {
		slog.Error("pulse error", "err", err)
	}
}

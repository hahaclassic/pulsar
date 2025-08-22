package main

import (
	"log"
	"log/slog"

	"github.com/hahaclassic/pulsar/internal/pulse"
	cpustat "github.com/hahaclassic/pulsar/internal/stat"
)

func main() {
	cpuStatParser, err := cpustat.NewParser()
	if err != nil {
		log.Fatalf("failed to init cpu stat parser: %s", err)
	}

	pulsar := pulse.NewPulsar(cpuStatParser)

	if err := pulsar.Start(); err != nil {
		slog.Error("pulse error", "err", err)
	}
}

package cpustat

import (
	"os"
	"strconv"
	"strings"
)

func ReadCPU() (idle, total uint64, err error) {
	data, err := os.ReadFile("/proc/stat")
	if err != nil {
		return
	}
	lines := strings.Split(string(data), "\n")
	fields := strings.Fields(lines[0]) // "cpu ..."
	var values []uint64
	for _, f := range fields[1:] {
		v, _ := strconv.ParseUint(f, 10, 64)
		values = append(values, v)
	}
	for _, v := range values {
		total += v
	}
	idle = values[3] // idle field
	return
}

// Usage returns cpu usage like: 56.2 (%)
func Usage(prevIdle, prevTotal, currIdle, currTotal uint64) float64 {
	idleTicks := float64(currIdle - prevIdle)
	totalTicks := float64(currTotal - prevTotal)

	return 100 * (1.0 - idleTicks/totalTicks)
}

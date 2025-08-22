package cpustat

import (
	"bufio"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
)

const (
	procStatPath = "/proc/stat"
)

var (
	ErrOpenFile          = errors.New("cpu parser: failed to open file")
	ErrNoTokensAvailable = errors.New("cpu parser: no tokens available")
	ErrScanRow           = errors.New("cpu parser: scan error")
)

type CPUStatParser struct {
	idle  [2]uint64
	total [2]uint64
}

func NewParser() (*CPUStatParser, error) {
	c := &CPUStatParser{}

	if err := c.ReadCPU(); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *CPUStatParser) ReadCPU() error {
	file, err := os.Open(procStatPath)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrOpenFile, err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			slog.Error("failed to close file", "path", procStatPath, "err", err)
		}
	}()

	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return fmt.Errorf("%w: %w", ErrScanRow, err)
		}

		return ErrNoTokensAvailable
	}

	row := scanner.Text() // "cpu ..."
	fields := strings.Fields(row)

	var total, idle uint64
	for i, f := range fields[1:] {
		v, _ := strconv.ParseUint(f, 10, 64)
		total += v
		if i == 3 { // idle field
			idle = v
		}
	}

	c.idle[0], c.total[0] = c.idle[1], c.total[1]
	c.total[1], c.idle[1] = total, idle

	return nil
}

func (c *CPUStatParser) Usage() float64 {
	if c.total[1] < c.total[0] {
		return 0
	}

	idleTicks := float64(c.idle[1] - c.idle[0])
	totalTicks := float64(c.total[1] - c.total[0])

	return 100 * (1.0 - idleTicks/totalTicks)
}

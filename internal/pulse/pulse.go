package pulse

import (
	"fmt"
	"math/rand/v2"
	"os"
	"os/signal"
	"syscall"
	"time"

	cpustat "github.com/hahaclassic/pulsar/internal/stat"
)

//
//         |                   |
//        /|                  /|
//     	 / |   |             / |   |
// ___ 	/  |  /| /\_______  /  |  /| /\____
//    \/   | / |/         \/   | / |/
//         |/                  |/
//         |	               |
//

// TODO:
// 2. graceful shoutdown (c затуханием пульса и возвращением курсора+)
// 3. вынести cpustat из пакета (сделать через интерфейс)
// 4. Resize

const (
	maxBPM   = 180
	minBPM   = 55
	bpmRatio = (float64(maxBPM) - float64(minBPM)) / 100

	width     = 60
	height    = 7
	beatWidth = 13

	symbPerSec     = int(float64(maxBPM) / 60.0 * (beatWidth + 1))
	sleepPerSymbMs = 1000 / symbPerSec
)

var beat = [height][beatWidth]rune{
	{' ', ' ', ' ', ' ', ' ', '|', ' ', ' ', ' ', ' ', ' ', ' ', ' '},
	{' ', ' ', ' ', ' ', '/', '|', ' ', ' ', ' ', ' ', ' ', ' ', ' '},
	{' ', ' ', ' ', '/', ' ', '|', ' ', ' ', ' ', '|', ' ', ' ', ' '},
	{' ', ' ', '/', ' ', ' ', '|', ' ', ' ', '/', '|', ' ', '/', '\\'},
	{'\\', '/', ' ', ' ', ' ', '|', ' ', '/', ' ', '|', '/', ' ', ' '},
	{' ', ' ', ' ', ' ', ' ', '|', '/', ' ', ' ', ' ', ' ', ' ', ' '},
	{' ', ' ', ' ', ' ', ' ', '|', ' ', ' ', ' ', ' ', ' ', ' ', ' '},
}

func getBPM(cpu float64) int {
	jitter := -2 + rand.Int32N(6) // A person doesn't have to be perfect.

	return min(int(bpmRatio*cpu+minBPM)+int(jitter), maxBPM)
}

func Start() error {
	fmt.Print("\033[?25l")       // скрыть курсор
	defer fmt.Print("\033[?25h") // показать курсор при выходе

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sig
		fmt.Print("\033[?25h")
		os.Exit(0)
	}()

	buffer := make([][]rune, height)
	for i := range buffer {
		buffer[i] = make([]rune, width)
		for j := range buffer[i] {
			buffer[i][j] = ' '
		}
	}

	fmt.Print("\033[2J") // очистить весь экран
	fmt.Print("\033[H")  // переместить курсор в левый верхний угол

	currIdx := -1

	prevIdle, prevTotal, err := cpustat.ReadCPU()
	if err != nil {
		return err
	}
	time.Sleep(time.Duration(sleepPerSymbMs) * time.Millisecond)

	for {
		currIdle, currTotal, err := cpustat.ReadCPU()
		if err != nil {
			return err
		}

		cpu := cpustat.Usage(prevIdle, prevTotal, currIdle, currTotal)
		bpm := getBPM(cpu)
		bps := float64(bpm) / 60.0 // beats per second
		interval := int(float64(symbPerSec)/bps) - beatWidth

		for step := range beatWidth + interval {
			// move
			if currIdx >= width-1 {
				for y := range height {
					copy(buffer[y], buffer[y][1:])
					buffer[y][width-1] = ' '
				}
			} else {
				currIdx++
			}

			// set new column
			switch {
			case step < beatWidth:
				for y := range height {
					buffer[y][currIdx] = beat[y][step]
				}
			case (step-beatWidth)%2 == 0:
				buffer[height/2][currIdx] = '_'
			default:
				buffer[height/2+1][currIdx] = '‾'
			}

			// print updated buffer
			fmt.Printf("\033[%dA", height) // поднять курсор вверх
			for y := range height {
				last := len(buffer[y]) - beatWidth
				fmt.Print(string(buffer[y][:last]))
				fmt.Printf("\033[3%dm%s\033[0m\n", getColor(bpm), string(buffer[y][last:]))
			}

			time.Sleep(time.Duration(sleepPerSymbMs) * time.Millisecond)
			prevIdle, prevTotal = currIdle, currTotal

			printStats(cpu, bpm)
			fmt.Printf("\033[%dA", height+1)
		}
	}
}

func getColor(bpm int) int {
	var color int

	switch {
	case bpm > 150:
		color = 1
	case bpm > 120:
		color = 3
	default:
		color = 2
	}

	return color
}

func printStats(cpu float64, bpm int) {
	fmt.Printf("\033[Kcpu: %0.1f%%   ", cpu)

	fmt.Printf("bpm: \033[3%dm%d\033[0m\n", getColor(bpm), bpm)
}

package pulse

import (
	"fmt"
	"math/rand/v2"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/term"
)

const (
	maxBPM   = 180
	minBPM   = 55
	bpmRatio = (float64(maxBPM) - float64(minBPM)) / 100

	highBPM     = 120
	criticalBPM = 150

	defaultWidth  = 60
	defaultHeight = 7
	beatWidth     = 13

	symbPerSec     = int(float64(maxBPM) / 60.0 * (beatWidth + 1))
	sleepPerSymbMs = 1000 / symbPerSec
)

type escapeCode string

const (
	// Cursor control
	esc            escapeCode = "\033"
	escHideCursor             = esc + "[?25l"
	escShowCursor             = esc + "[?25h"
	escClearScreen            = esc + "[2J"
	escMoveHome               = esc + "[H"

	// Colors
	escColorReset  = esc + "[0m"
	escColorRed    = esc + "[31m"
	escColorGreen  = esc + "[32m"
	escColorYellow = esc + "[33m"

	// Line control
	escClearLine = esc + "[K"
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

var beat = [defaultHeight][beatWidth]rune{
	{' ', ' ', ' ', ' ', ' ', '|', ' ', ' ', ' ', ' ', ' ', ' ', ' '},
	{' ', ' ', ' ', ' ', '/', '|', ' ', ' ', ' ', ' ', ' ', ' ', ' '},
	{' ', ' ', ' ', '/', ' ', '|', ' ', ' ', ' ', '|', ' ', ' ', ' '},
	{' ', ' ', '/', ' ', ' ', '|', ' ', ' ', '/', '|', ' ', '/', '\\'},
	{'\\', '/', ' ', ' ', ' ', '|', ' ', '/', ' ', '|', '/', ' ', ' '},
	{' ', ' ', ' ', ' ', ' ', '|', '/', ' ', ' ', ' ', ' ', ' ', ' '},
	{' ', ' ', ' ', ' ', ' ', '|', ' ', ' ', ' ', ' ', ' ', ' ', ' '},
}

type CPUStatParser interface {
	ReadCPU() error
	Usage() float64
}

type Pulsar struct {
	cpu CPUStatParser
	buf *buffer
}

func NewPulsar(parser CPUStatParser) *Pulsar {
	return &Pulsar{
		cpu: parser,
		buf: newBuffer(defaultWidth, defaultHeight),
	}
}

// getBPM calculates heart rate (BPM) based on CPU usage and adds small jitter.
func (p *Pulsar) getBPM(cpu float64) int {
	jitter := -2 + rand.Int32N(4) // A person doesn't have to be perfect.

	return min(int(bpmRatio*cpu+minBPM)+int(jitter), maxBPM)
}

func (p *Pulsar) getColor(bpm int) escapeCode {
	var color escapeCode

	switch {
	case bpm > criticalBPM:
		color = escColorRed
	case bpm > highBPM:
		color = escColorYellow
	default:
		color = escColorGreen
	}

	return color
}

// Start runs the main loop: updates CPU, calculates BPM, and draws cardiogram.
func (p *Pulsar) Start() error {
	fmt.Print(escHideCursor)
	defer fmt.Print(escShowCursor)

	sigExit := make(chan os.Signal, 1)
	signal.Notify(sigExit, os.Interrupt, syscall.SIGTERM)

	sigResize := make(chan os.Signal, 1)
	signal.Notify(sigResize, syscall.SIGWINCH) // changing the size of the terminal

	p.updateSize()
	go func() {
		for range sigResize {
			p.updateSize()
		}
	}()

	if err := p.cpu.ReadCPU(); err != nil {
		return err
	}
	time.Sleep(time.Duration(sleepPerSymbMs) * time.Millisecond)

	for {
		select {
		case <-sigExit:
			p.printZeroBPM()
			return nil

		default:
			if err := p.cpu.ReadCPU(); err != nil {
				return err
			}

			cpu := p.cpu.Usage()
			bpm := p.getBPM(cpu)
			bps := float64(bpm) / 60.0 // beats per second
			interval := int(float64(symbPerSec)/bps) - beatWidth

			p.printCardiacCycle(cpu, bpm, interval)
		}
	}
}

func (p *Pulsar) printZeroBPM() {
	fmt.Print(escClearScreen)
	fmt.Print(escMoveHome)

	for range symbPerSec {
		p.buf.shift()
		p.buf.setRune(defaultHeight/2, '_')
		p.buf.printAllColored(escColorReset)

		time.Sleep(time.Duration(sleepPerSymbMs) * time.Millisecond)
	}
}

func (p *Pulsar) printCardiacCycle(cpu float64, bpm int, interval int) {
	for step := range beatWidth + interval {
		p.buf.shift()

		// set new column
		switch {
		case step < beatWidth:
			for i := range defaultHeight {
				p.buf.setRune(i, beat[i][step])
			}
		case (step-beatWidth)%2 == 0:
			p.buf.setRune(defaultHeight/2, '_')
		default:
			p.buf.setRune(defaultHeight/2+1, '‾')
		}

		p.printBufferWithStat(cpu, bpm)

		time.Sleep(time.Duration(sleepPerSymbMs) * time.Millisecond)
	}
}

// printBufferWithStat prints the buffer and CPU/BPM stats below it.
func (p *Pulsar) printBufferWithStat(cpu float64, bpm int) {
	color := p.getColor(bpm)

	// print cardiogram
	p.buf.printWithColoredBeat(color)

	// print stat
	fmt.Printf("%scpu: %0.1f%%   ", escClearLine, cpu)
	fmt.Printf("bpm: %s%d%s\n", color, bpm, escColorReset)
	fmt.Print(escMoveHome)
}

// updateSize fetches terminal size and resizes the buffer accordingly.
func (p *Pulsar) updateSize() {
	fmt.Print(escClearScreen)
	fmt.Print(escMoveHome)

	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return
	}

	p.buf.resize(width)
}

package pulse

import (
	"fmt"
	"sync"
)

type buffer struct {
	data    [][]rune
	w       int
	h       int
	currIdx int // current column index

	mu *sync.RWMutex
}

// newBuffer creates a new buffer of given width and height, filled with spaces.
func newBuffer(width, height int) *buffer {
	data := make([][]rune, height)
	for i := range data {
		data[i] = make([]rune, width)
		for j := range data[i] {
			data[i][j] = ' '
		}
	}

	return &buffer{
		data: data,
		w:    width,
		h:    height,
		mu:   &sync.RWMutex{},
	}
}

func (b *buffer) resize(width int) { // Height doesn't matter.
	b.mu.Lock()
	defer b.mu.Unlock()

	if width <= 0 || width == b.w {
		return
	}

	if width > b.w {
		for i := range b.data {
			for range width - b.w {
				b.data[i] = append(b.data[i], ' ')
			}
		}
	} else {
		diff := b.w - width
		b.currIdx = b.currIdx - diff

		for i := range b.data {
			copy(b.data[i], b.data[i][diff:])
			b.data[i] = b.data[i][:width]
		}
	}

	b.w = width
}

// shift moves the buffer one column to the left.
// If currIdx is not yet at the rightmost column, just increment it.
func (b *buffer) shift() {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.currIdx < b.w-1 {
		b.currIdx++
		return
	}

	for i := range b.h {
		copy(b.data[i], b.data[i][1:])
		b.data[i][b.w-1] = ' '
	}
}

// printWithColoredBeat prints the buffer to the terminal,
// coloring the last part (current beat) in the given color.
func (b *buffer) printWithColoredBeat(color escapeCode) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	fmt.Print(escMoveHome) // move the cursor up
	for i := range b.h {
		if b.currIdx < beatWidth {
			fmt.Printf("%s%s%s\n", color, string(b.data[i]), escColorReset)
		} else {
			idx := b.currIdx - beatWidth
			fmt.Print(string(b.data[i][:idx]))
			fmt.Printf("%s%s%s\n", color, string(b.data[i][idx:]), escColorReset)
		}
	}
}

// printAllColored prints the entire buffer in a single color.
func (b *buffer) printAllColored(color escapeCode) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	fmt.Print(escMoveHome) // move the cursor up
	fmt.Print(color)
	for i := range b.h {
		if b.currIdx < beatWidth {
			fmt.Printf("%s\n", string(b.data[i]))
		} else {
			idx := b.currIdx - beatWidth
			fmt.Print(string(b.data[i][:idx]))
			fmt.Printf("%s\n", string(b.data[i][idx:]))
		}
	}
	fmt.Print(escColorReset)
}

// setRune sets a rune at the current column (currIdx) for a given row.
func (b *buffer) setRune(rowIdx int, r rune) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if rowIdx >= 0 && rowIdx < b.h {
		b.data[rowIdx][b.currIdx] = r
	}
}

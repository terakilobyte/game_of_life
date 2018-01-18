// Package life is the manages the "game" state
// Shamelessly taken from https://golang.org/doc/play/life.go
package life

import "math/rand"

// Life stores the state of a round of Conway's Game of Life.
type Life struct {
	A, b *Grid
	w, h int
}

// NewLife returns a new Life game state with a random initial state.
func NewLife(w, h, size int) *Life {
	a := NewGrid(w, h, size)
	for i := 0; i < (w * h / 15); i++ {
		a.Set(rand.Intn(w), rand.Intn(h), true)
	}
	return &Life{
		A: a, b: NewGrid(w, h, size),
		w: w, h: h,
	}
}

// Step advances the game by one instant, recomputing and updating all cells.
func (l *Life) Step() {
	// Update the state of the next field (b) from the current field (a).
	for y := 0; y < l.h; y++ {
		for x := 0; x < l.w; x++ {
			l.b.Set(x, y, l.A.Next(x, y))
		}
	}
	// Swap fields a and b.
	l.A, l.b = l.b, l.A
}

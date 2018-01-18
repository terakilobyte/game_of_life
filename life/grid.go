package life

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"golang.org/x/image/colornames"
)

// Shamelessly taken/inspired by https://golang.org/doc/play/life.go
// Grid is the structure in which the cellular automota live
type Grid struct {
	w        int
	h        int
	cellSize int
	Cells    [][]bool
}

// NewGrid constructs a new Grid
func NewGrid(cols, rows, size int) *Grid {
	cells := make([][]bool, cols)
	for i := 0; i < rows; i++ {
		cells[i] = make([]bool, rows)
	}
	return &Grid{w: cols, h: rows, cellSize: size, Cells: cells}
}

// Alive returns whether the specified position is alive
func (g *Grid) Alive(x, y int) bool {
	x += g.w
	x %= g.w
	y += g.h
	y %= g.h
	return g.Cells[y][x]
}

// Set sets the state of a specific location
func (g *Grid) Set(x, y int, state bool) {
	g.Cells[y][x] = state
}

// Draw draws the grid
func (g *Grid) Draw(imd *imdraw.IMDraw) {
	for i := 0; i < g.w; i++ {
		for j := 0; j < g.h; j++ {
			imd.Push(
				pixel.V(float64(i*g.cellSize), float64(j*g.cellSize)),
				pixel.V(float64(i*g.cellSize+g.cellSize), float64(j*g.cellSize+g.cellSize)),
			)
			if g.Alive(i, j) {
				imd.Color = colornames.Black
			} else {
				imd.Color = colornames.White
			}
			imd.Rectangle(0)
		}
	}
}

// Next returns the next state
func (g *Grid) Next(x, y int) bool {
	// Count the adjacent cells that are alive.
	alive := 0
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if (j != 0 || i != 0) && g.Alive(x+i, y+j) {
				alive++
			}
		}
	}
	// Return next state according to the game rules:
	//   exactly 3 neighbors: on,
	//   exactly 2 neighbors: maintain current state,
	//   otherwise: off.
	return alive == 3 || alive == 2 && g.Alive(x, y)
}

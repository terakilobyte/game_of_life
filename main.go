package main

import (
	"math/rand"
	"time"

	"golang.org/x/image/colornames"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
)

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Pixel Rocks!",
		Bounds: pixel.R(0, 0, 800, 800),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	rand.Seed(time.Now().UnixNano())
	win.Clear(colornames.White)

	size := 5
	gridHeight := win.Bounds().H()
	gridWidth := win.Bounds().W()
	rows := int(gridHeight) / size
	cols := int(gridWidth) / size

	gridDraw := imdraw.New(nil)
	game := NewLife(cols, rows, size)
	last := time.Now()
	for !win.Closed() {
		// game loop
		if time.Since(last).Nanoseconds() > int64(33*time.Millisecond) {
			gridDraw.Clear()
			game.a.Draw(gridDraw)
			gridDraw.Draw(win)
			game.Step()
			last = time.Now()
		}
		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}

// Life stores the state of a round of Conway's Game of Life.
type Life struct {
	a, b *Grid
	w, h int
}

// NewLife returns a new Life game state with a random initial state.
func NewLife(w, h, size int) *Life {
	a := NewGrid(w, h, size)
	for i := 0; i < (w * h / 20); i++ {
		a.Set(rand.Intn(w), rand.Intn(h), true)
	}
	return &Life{
		a: a, b: NewGrid(w, h, size),
		w: w, h: h,
	}
}

// Step advances the game by one instant, recomputing and updating all cells.
func (l *Life) Step() {
	// Update the state of the next field (b) from the current field (a).
	for y := 0; y < l.h; y++ {
		for x := 0; x < l.w; x++ {
			l.b.Set(x, y, l.a.Next(x, y))
		}
	}
	// Swap fields a and b.
	l.a, l.b = l.b, l.a
}

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

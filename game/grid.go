package game

import "errors"

var (
	ErrMalformedGrid = errors.New("malformed grid")
)

type Grid [][]int

func newGrid(width, height int) Grid {
	grid := make(Grid, height)
	for y := 0; y < height; y++ {
		grid[y] = make([]int, width)
	}
	return grid
}

func (g Grid) Get(x, y int) int {
	return g[y][x]
}

func (g Grid) Set(x, y, value int) {
	g[y][x] = value
}

func (g Grid) SetAll(value int) {
	for y := 0; y < g.GetHeight(); y++ {
		for x := 0; x < g.GetWidth(); x++ {
			g.Set(x, y, value)
		}
	}
}

func (g Grid) GetWidth() int {
	return len(g[0])
}

func (g Grid) GetHeight() int {
	return len(g)
}

func (g Grid) Validate() error {
	for _, row := range g {
		if len(row) != len(g[0]) {
			return ErrMalformedGrid
		}
	}
	return nil
}

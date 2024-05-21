package game

import (
	"errors"
	"math/rand"
	"time"
)

var (
	ErrInvalidFieldSize = errors.New("invalid field size")
	ErrOutOfBounds      = errors.New("out of bounds")
	ErrDuplicateMine    = errors.New("duplicate mine")
)

const (
	CellUnrevealed = -1
	CellMine       = -2
	CellFlag       = -3
)

type Coordinate struct {
	X int
	Y int
}

type State int

const (
	StatePlaying State = iota
	StateWon
	StateLost
)

type Game struct {
	gridWidth     int
	gridHeight    int
	mines         []Coordinate
	flags         []Coordinate
	revealedCells []Coordinate
	gameOver      bool
}

type Option func(*Game)

func New(width, height int) (*Game, error) {
	if width <= 0 || height <= 0 {
		return nil, ErrInvalidFieldSize
	}

	g := &Game{
		gridWidth:  width,
		gridHeight: height,
	}

	return g, nil
}

func NewFromGrid(grid Grid) (*Game, error) {
	if err := grid.Validate(); err != nil {
		return nil, err
	}

	g, err := New(grid.GetWidth(), grid.GetHeight())
	if err != nil {
		return nil, err
	}

	err = g.PlaceMines(grid)
	if err != nil {
		return nil, err
	}

	return g, nil
}

func (g *Game) GetGridWidth() int {
	return g.gridWidth
}

func (g *Game) GetGridHeight() int {
	return g.gridHeight
}

func (g *Game) PlaceMine(x int, y int) error {
	if !g.coordinatesInBounds(x, y) {
		return ErrOutOfBounds
	}

	if g.cellHasMine(x, y) {
		return ErrDuplicateMine
	}

	g.mines = append(g.mines, Coordinate{x, y})
	return nil
}

func (g *Game) GetMineCount() int {
	return len(g.mines)
}

func (g *Game) RevealCell(x int, y int) int {
	if g.gameOver {
		return -1
	}
	if g.cellHasMine(x, y) {
		g.gameOver = true
		return -1
	}

	g.revealedCells = append(g.revealedCells, Coordinate{x, y})
	g.RemoveFlag(x, y)

	n := g.getNumberOfAdjacentMines(x, y)
	if n > 0 {
		return n
	}

	g.revealAdjacentCells(x, y)

	return 0
}

func (g *Game) revealAdjacentCells(x int, y int) {
	for x2 := x - 1; x2 <= x+1; x2++ {
		for y2 := y - 1; y2 <= y+1; y2++ {
			if !g.coordinatesInBounds(x2, y2) {
				continue
			}

			if g.cellIsRevealed(x2, y2) {
				continue
			}

			if !g.cellHasMine(x2, y2) && !g.cellHasFlag(x2, y2) {
				g.RevealCell(x2, y2)
			}
		}
	}
}

func (g *Game) RemoveFlag(x int, y int) {
	for i, f := range g.flags {
		if f.X == x && f.Y == y {
			g.flags = append(g.flags[:i], g.flags[i+1:]...)
			return
		}
	}
}

func (g *Game) State() State {
	if g.gameOver {
		return StateLost
	}

	if g.checkWinCondition() {
		return StateWon
	}

	return StatePlaying
}

func (g *Game) cellHasMine(x int, y int) bool {
	for _, m := range g.mines {
		if m.X == x && m.Y == y {
			return true
		}
	}
	return false
}

func (g *Game) cellHasFlag(x int, y int) bool {
	for _, f := range g.flags {
		if f.X == x && f.Y == y {
			return true
		}
	}
	return false
}

func (g *Game) coordinatesInBounds(x int, y int) bool {
	return x >= 0 && x < g.gridWidth && y >= 0 && y < g.gridHeight
}

func (g *Game) PlaceFlag(x int, y int) error {
	if !g.coordinatesInBounds(x, y) {
		return ErrOutOfBounds
	}

	g.flags = append(g.flags, Coordinate{x, y})

	return nil
}

func (g *Game) GetFlagCount() int {
	return len(g.flags)
}

func (g *Game) checkWinCondition() bool {
	if len(g.flags) != len(g.mines) {
		return false
	}

	for _, m := range g.mines {
		if !g.cellHasFlag(m.X, m.Y) {
			return false
		}
	}

	return true
}

func (g *Game) getNumberOfAdjacentMines(revealX int, revealY int) int {
	num := 0
	for x := revealX - 1; x <= revealX+1; x++ {
		for y := revealY - 1; y <= revealY+1; y++ {
			if !g.coordinatesInBounds(x, y) {
				continue
			}

			if g.cellHasMine(x, y) {
				num++
			}
		}
	}

	return num
}

func (g *Game) PlaceMines(grid Grid) error {
	if grid.GetHeight() != g.gridHeight || grid.GetWidth() != g.gridWidth {
		return ErrInvalidFieldSize
	}
	if err := grid.Validate(); err != nil {
		return err
	}

	for y, row := range grid {
		for x, cell := range row {
			if cell == 1 {
				err := g.PlaceMine(x, y)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (g *Game) GetGrid() Grid {
	grid := newGrid(g.gridWidth, g.gridHeight)
	grid.SetAll(CellUnrevealed)

	for _, c := range g.revealedCells {
		grid.Set(c.X, c.Y, g.getNumberOfAdjacentMines(c.X, c.Y))
	}

	for _, f := range g.flags {
		grid.Set(f.X, f.Y, CellFlag)
	}

	if g.State() == StateLost {
		for _, m := range g.mines {
			grid.Set(m.X, m.Y, CellMine)
		}
	}

	return grid
}

func (g *Game) cellIsRevealed(x int, y int) bool {
	for _, c := range g.revealedCells {
		if c.X == x && c.Y == y {
			return true
		}
	}
	return false
}

func (g *Game) ToggleFlag(x int, y int) {
	if g.cellHasFlag(x, y) {
		g.RemoveFlag(x, y)
	} else {
		g.PlaceFlag(x, y)
	}
}

func (g *Game) PlaceRandomMines(count int) error {
	if count < 0 || count > g.gridWidth*g.gridHeight {
		return ErrInvalidFieldSize
	}

	rand.Seed(time.Now().UnixNano())

	availablePositions := make([]Coordinate, 0, g.gridWidth*g.gridHeight)
	for x := 0; x < g.gridWidth; x++ {
		for y := 0; y < g.gridHeight; y++ {
			availablePositions = append(availablePositions, Coordinate{x, y})
		}
	}

	// Shuffle the list of positions
	rand.Shuffle(len(availablePositions), func(i, j int) {
		availablePositions[i], availablePositions[j] = availablePositions[j], availablePositions[i]
	})

	// Place mines in the first 'count' positions
	for i := 0; i < count; i++ {
		pos := availablePositions[i]
		err := g.PlaceMine(pos.X, pos.Y)
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *Game) Reset() {
	g.mines = nil
	g.flags = nil
	g.revealedCells = nil
	g.gameOver = false
	g.PlaceRandomMines(10)
}

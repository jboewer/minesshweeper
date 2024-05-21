package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/jboewer/mine-ssh-weeper/game"
	"log"
	"strconv"
	"strings"
)

func NewGameModel(g *game.Game) GameModel {
	return GameModel{
		Game: g,
		Cursor: Cursor{
			game: g,
		},
	}
}

type Cursor struct {
	x    int
	y    int
	game *game.Game
}

func (c *Cursor) Up() {
	c.y--
	if c.y < 0 {
		c.y = 0
	}
}

func (c *Cursor) Down() {
	c.y++
	if c.y >= c.game.GetGridHeight() {
		c.y = c.game.GetGridHeight() - 1
	}
}

func (c *Cursor) Left() {
	c.x--
	if c.x < 0 {
		c.x = 0
	}
}

func (c *Cursor) Right() {
	c.x++
	if c.x >= c.game.GetGridWidth() {
		c.x = c.game.GetGridWidth() - 1
	}
}

type GameModel struct {
	Game   *game.Game
	Cursor Cursor
}

func (gv GameModel) View() string {
	rendered := &strings.Builder{}

	gv.renderGameGrid(rendered)

	// Write game state
	switch gv.Game.State() {
	case game.StatePlaying:
		rendered.WriteString("Playing")
	case game.StateWon:
		rendered.WriteString("Won")
	case game.StateLost:
		rendered.WriteString("Lost")
	}

	rendered.WriteString("\n")

	gv.renderInstructions(rendered)

	// Send the UI for rendering
	return rendered.String()
}

func (gv GameModel) renderGameGrid(rendered *strings.Builder) {
	grid := gv.Game.GetGrid()

	rows := make([][]string, grid.GetHeight())
	for y := 0; y < grid.GetHeight(); y++ {
		rows[y] = make([]string, grid.GetWidth())

		for x := 0; x < grid.GetWidth(); x++ {
			c := grid.Get(x, y)

			switch c {
			case game.CellUnrevealed:
				rows[y][x] = " "
			case game.CellFlag:
				rows[y][x] = "F"
			case game.CellMine:
				rows[y][x] = "M"
			default:
				rows[y][x] = strconv.Itoa(c)
			}

		}
	}

	tbl := table.New().
		Border(lipgloss.NormalBorder()).
		BorderRow(true).
		BorderColumn(true).
		Rows(rows...).
		StyleFunc(func(row, col int) lipgloss.Style {
			fg, bg := getCellColors(rows[row-1][col])

			if row == gv.Cursor.y+1 && col == gv.Cursor.x {
				fg, bg = getCursorColors(fg, bg)
			}

			return lipgloss.NewStyle().
				Foreground(fg).
				Background(bg).
				Padding(0, 1)
		})

	rendered.WriteString(tbl.Render())
	rendered.WriteString("\n")
}

func (gv GameModel) Init() tea.Cmd {
	return nil
}

func (gv GameModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return gv, tea.Quit
		case "w":
			gv.Cursor.Up()
		case "a":
			gv.Cursor.Left()
		case "s":
			gv.Cursor.Down()
		case "d":
			gv.Cursor.Right()
		case "f":
			gv.Game.ToggleFlag(gv.Cursor.x, gv.Cursor.y)
		case " ":
			gv.Game.RevealCell(gv.Cursor.x, gv.Cursor.y)
		case "r":
			gv.Reset()
		}
	}

	return gv, nil
}

func (gv GameModel) renderInstructions(rendered *strings.Builder) {
	rendered.WriteString("WASD: Move Around, F: Toggle Flag, Space: Reveal, R: Reset, Q: Quit")
}

func (gv GameModel) Reset() {
	log.Println("Resetting game")
	gv.Game.Reset()
}

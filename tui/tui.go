package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/jboewer/minesweeper/game"
	"strconv"
	"strings"
)

func NewGameView(g *game.Game) GameView {
	return GameView{
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

type GameView struct {
	Game   *game.Game
	Cursor Cursor
}

func (t GameView) View() string {
	rendered := &strings.Builder{}

	t.renderGameGrid(rendered)

	// Write game state
	switch t.Game.State() {
	case game.StatePlaying:
		rendered.WriteString("Playing")
	case game.StateWon:
		rendered.WriteString("Won")
	case game.StateLost:
		rendered.WriteString("Lost")
	}

	// Send the UI for rendering
	return rendered.String()
}

func (t GameView) renderGameGrid(rendered *strings.Builder) {
	grid := t.Game.GetGrid()

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
			var fg lipgloss.TerminalColor = lipgloss.NoColor{}
			var bg lipgloss.TerminalColor = lipgloss.NoColor{}

			switch rows[row-1][col] {
			case "0":
				fg = lipgloss.Color("#292929")
			case "1":
				fg = lipgloss.Color("#74adf2")
			case "2":
				fg = lipgloss.Color("#00FF00")
			case "3":
				fg = lipgloss.Color("#FF0000")
			case "4":
				fg = lipgloss.Color("#28706d")
			case "5":
				fg = lipgloss.Color("#b06446")
			case "6":
				fg = lipgloss.Color("#FF0000")
			case "7":
				fg = lipgloss.Color("#8a7101")
			case "8":
				fg = lipgloss.Color("#111")
				bg = lipgloss.Color("#bfbfbf")
			case "*":
				fg = lipgloss.Color("#FF33FF")
			case "F":
				bg = lipgloss.Color("#ffee00")
				fg = lipgloss.Color("#111")
			case "B":
				bg = lipgloss.Color("#FF0000")
				fg = lipgloss.Color("#111")
			case " ":
				bg = lipgloss.Color("#bfbfbf")
			}

			if row == t.Cursor.y+1 && col == t.Cursor.x {
				fg = lipgloss.Color("#111")
				bg = lipgloss.Color("#FF33FF")
			}

			return lipgloss.NewStyle().
				Foreground(fg).
				Background(bg).
				Padding(0, 1)
		})

	rendered.WriteString(tbl.Render())
	rendered.WriteString("\n")
}
func (t GameView) Init() tea.Cmd {
	return nil
}

func (t GameView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return t, tea.Quit
		case "w":
			t.Cursor.Up()
		case "a":
			t.Cursor.Left()
		case "s":
			t.Cursor.Down()
		case "d":
			t.Cursor.Right()
		case "f":
			t.Game.ToggleFlag(t.Cursor.x, t.Cursor.y)
		case " ":
			t.Game.RevealCell(t.Cursor.x, t.Cursor.y)
		}

	}

	return t, nil
}

package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jboewer/minesweeper/game"
	"github.com/jboewer/minesweeper/tui"
	"os"
)

func main() {
	g, err := game.New(10, 10)
	if err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}

	g.PlaceRandomMines(10)

	p := tea.NewProgram(
		tui.NewGameView(g),
		tea.WithAltScreen(),
	)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

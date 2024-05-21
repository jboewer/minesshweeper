package tui

import "github.com/charmbracelet/lipgloss"

func getCursorColors(
	fg lipgloss.TerminalColor,
	bg lipgloss.TerminalColor,
) (lipgloss.TerminalColor, lipgloss.TerminalColor) {
	fg = lipgloss.Color("#111")
	bg = lipgloss.Color("#FF33FF")
	return fg, bg
}

func getCellColors(cell string) (lipgloss.TerminalColor, lipgloss.TerminalColor) {
	var fg lipgloss.TerminalColor = lipgloss.NoColor{}
	var bg lipgloss.TerminalColor = lipgloss.NoColor{}

	switch cell {
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
	return fg, bg
}

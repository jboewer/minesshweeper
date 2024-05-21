package tui

import "github.com/charmbracelet/lipgloss"

var colorMap = map[string]struct {
	fg lipgloss.TerminalColor
	bg lipgloss.TerminalColor
}{
	"0": {fg: lipgloss.Color("#292929"), bg: lipgloss.NoColor{}},
	"1": {fg: lipgloss.Color("#74adf2"), bg: lipgloss.NoColor{}},
	"2": {fg: lipgloss.Color("#00FF00"), bg: lipgloss.NoColor{}},
	"3": {fg: lipgloss.Color("#FF0000"), bg: lipgloss.NoColor{}},
	"4": {fg: lipgloss.Color("#28706d"), bg: lipgloss.NoColor{}},
	"5": {fg: lipgloss.Color("#b06446"), bg: lipgloss.NoColor{}},
	"6": {fg: lipgloss.Color("#FF0000"), bg: lipgloss.NoColor{}},
	"7": {fg: lipgloss.Color("#8a7101"), bg: lipgloss.NoColor{}},
	"8": {fg: lipgloss.Color("#111"), bg: lipgloss.Color("#bfbfbf")},
	"*": {fg: lipgloss.Color("#FF33FF"), bg: lipgloss.NoColor{}},
	"F": {fg: lipgloss.Color("#111"), bg: lipgloss.Color("#ffee00")},
	"B": {fg: lipgloss.Color("#111"), bg: lipgloss.Color("#FF0000")},
	" ": {fg: lipgloss.NoColor{}, bg: lipgloss.Color("#bfbfbf")},
}

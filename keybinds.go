package main

import (
	"github.com/charmbracelet/bubbles/key"
)

type keyMap struct {
	Up    key.Binding
	Down  key.Binding
	Quit  key.Binding
	Back  key.Binding
	Top   key.Binding
	Enter key.Binding
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc", "left"),
		key.WithHelp("←/esc", "go back"),
	),
	Top: key.NewBinding(
		key.WithKeys("t"),
		key.WithHelp("t", "go to top"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
	),
}

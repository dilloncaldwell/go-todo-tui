package main

import (
	"github.com/charmbracelet/bubbles/key"
)

// Key bindings
type keyMap struct {
	up      key.Binding
	down    key.Binding
	add     key.Binding
	delete  key.Binding
	edit    key.Binding
	toggle  key.Binding
	sort    key.Binding
	filter  key.Binding
	confirm key.Binding
	cancel  key.Binding
	quit    key.Binding
	help    key.Binding
}

var keys = keyMap{
	up: key.NewBinding(
		key.WithKeys("k", "up"),
		key.WithHelp("↑/k", "up"),
	),
	down: key.NewBinding(
		key.WithKeys("j", "down"),
		key.WithHelp("↓/j", "down"),
	),
	add: key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "add"),
	),
	delete: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "delete"),
	),
	edit: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "edit"),
	),
	toggle: key.NewBinding(
		key.WithKeys("t", "space"),
		key.WithHelp("t/space", "toggle"),
	),
	sort: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "Sort"),
	),
	filter: key.NewBinding(
		key.WithKeys("f"),
		key.WithHelp("f", "filter"),
	),
	confirm: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "confirm"),
	),
	cancel: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "cancel"),
	),
	quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "help"),
	),
}

// ShortHelp returns keybindings to be shown in the mini help view. It's part
// of the key.Map interface.
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.add, k.delete, k.edit, k.toggle, k.sort, k.filter, k.help, k.quit}
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.up, k.down, k.add, k.delete, k.edit, k.toggle, k.sort, k.filter, k.help, k.quit}, // first column
	}
}

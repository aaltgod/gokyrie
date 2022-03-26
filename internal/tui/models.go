package cliclient

import (
	"os"

	"github.com/aaltgod/gokyrie/internal/tui/bubbles/graph"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
)

type Model struct {
	Logger   *os.File
	List     list.Model
	Graph    graph.Model
	ViewPort viewport.Model
}

type GraphMsg struct {
	Body string
}

type DelegateKeyMap struct {
	Enter     key.Binding
	BackSpace key.Binding
}

package app

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
)

type Team struct {
	TeamName, TeamIP string
}

type Model struct {
	Team Team
	List list.Model
}

type DelegateKeyMap struct {
	Enter     key.Binding
	BackSpace key.Binding
}

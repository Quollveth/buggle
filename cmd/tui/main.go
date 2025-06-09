package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	quitting bool
}

func initModel() model        { return model{} }
func (m model) Init() tea.Cmd { return nil }
func (m model) View() string  { return "" }

func (m model) Update(tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func main() {
	p := tea.NewProgram(initModel())

	if _, err := p.Run(); err != nil {
		fmt.Printf("peepee went poopoo: %v", err)
		os.Exit(1)
	}
}

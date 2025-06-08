package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

// no enum 😭
type screen int

const (
	SCREEN_MAIN screen = iota
)

// returned by the submodel if it wishes to change the main screen
type ChangeScreenMsg struct {
	To screen
}

// we have inheritance at home
type SubModel interface {
	Init() tea.Cmd
	View() string
	// a submodel can return a message passed to the main model
	Update(tea.Msg) (SubModel, tea.Msg, tea.Cmd)
}

type mainModel struct {
	currentScreen screen
	screens       []SubModel
	quitting      bool
}

func MainModel() mainModel {
	//TODO: this thing lmao
	mainScreen := MainScreenModel()

	return mainModel{
		currentScreen: SCREEN_MAIN,
		quitting:      false,
		screens: []SubModel{mainScreen},
	}
}

func (m mainModel) Init() tea.Cmd {return m.screens[m.currentScreen].Init()}
func (m mainModel) View() string {return m.screens[m.currentScreen].View()}
func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	//-- global input
	if msg, ok := msg.(tea.KeyMsg); ok {
		if k := msg.String(); k == "ctrl+c" {
			return m, tea.Quit
		}
	}
	//-- screen specific input
	updatedModel, received, cmd := m.screens[m.currentScreen].Update(msg)
	m.screens[m.currentScreen] = updatedModel

	if(received == nil){
		return m, cmd
	}
	
	if received, ok := received.(ChangeScreenMsg); ok {
		m.currentScreen = received.To
	}

	return m, cmd
}

func main() {
	p := tea.NewProgram(MainModel())

	if _, err := p.Run(); err != nil {
		fmt.Printf("peepee went poopoo: %v", err)
		os.Exit(1)
	}
}
